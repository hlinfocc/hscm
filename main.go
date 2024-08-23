package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"hscm/pkg/utils"
	"hscm/pkg/version"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WebHook string   `yaml:"webhook"`
	Hosts   []string `yaml:"hosts"`
}

type Payload struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

/**
* 命令行参数结构体
 */
type Args struct {
	Version bool
	Monitor bool
}

func main() {
	// 打开日志文件，如果文件不存在则创建
	file, err := os.OpenFile("/var/log/hscm.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// 设置日志输出到文件
	log.SetOutput(file)

	args := initParams()
	cfg := loadConfig()
	log.Println(cfg.Hosts)
	if args.Version {
		fmt.Println(version.Full())
	} else {
		for i := 0; i < len(cfg.Hosts); i++ {
			log.Println(cfg.Hosts[i])
			expiresDays := getCertificateExpiresDays(cfg.Hosts[i])
			log.Println("有效期剩余: ", expiresDays)
			if expiresDays < 10 {
				webhookNotice(cfg.WebHook, cfg.Hosts[i], expiresDays)
			}
		}
	}

}

/**
* 初始化命令行参数信息
 */
func initParams() Args {
	args := Args{}
	flag.BoolVar(&args.Version, "v", args.Version, "显示版本信息")
	flag.BoolVar(&args.Monitor, "m", args.Monitor, "执行监控")
	flag.Parse()
	return args
}

func loadConfig() Config {
	var cfg Config
	isexist := utils.FileExists("/etc/hscm.yml")
	if isexist {
		// 读取 YAML 文件
		yamlFile, err := os.ReadFile("/etc/hscm.yml")
		if err != nil {
			log.Fatalf("Error reading YAML file: %v", err)
		}
		// 解析 YAML 文件
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			log.Fatalf("Error unmarshalling YAML data: %v", err)
		}
		// log.Println(cfg)
	}
	return cfg
}
func webhookNotice(webhookUrl string, domain string, days int) {
	payload := Payload{}
	payload.Msgtype = "text"
	payload.Text.Content = fmt.Sprintf("以下域名https证书即将到期：%s，有效期仅剩余%d天", domain, days)
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}
	// 创建 HTTP 客户端
	client := &http.Client{}
	// 创建 POST 请求
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	// 设置 Content-Type 为 application/json
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// 输出响应状态码和响应体
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", body)
}
func getCertificateExpiresDays(urlStr string) int {
	// 获取主机名和端口号
	host, port, err := netHostPort(urlStr)
	if err != nil {
		fmt.Println("Error getting host and port:", err)
		return -1
	}

	address := net.JoinHostPort(host, port)
	log.Println("address:", address)
	// 建立到目标主机的 TLS 连接
	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true, // 为了演示，忽略证书验证
	})
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return -1
	}
	defer conn.Close()

	// 获取证书链
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		fmt.Println("No certificates received")
		return -1
	}

	// 取第一个证书
	cert := certs[0]

	// 计算证书过期时间
	expiry := cert.NotAfter
	now := time.Now()
	duration := expiry.Sub(now)

	// 一天等于 24 小时
	dayDuration := 24 * time.Hour

	// 总天数
	totalDays := int(duration / dayDuration)

	fmt.Printf("Certificate expires on: %s\n", expiry)
	fmt.Printf("Time until expiration: %s\n", duration)

	return totalDays
}

// netHostPort 从 URL 中提取主机名和端口号
func netHostPort(urlStr string) (host string, port string, err error) {
	if !strings.HasPrefix(urlStr, "https://") || !strings.HasPrefix(urlStr, "http://") {
		urlStr = fmt.Sprintf("https://%s", urlStr)
	}
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
	}
	host, port, err = net.SplitHostPort(parsedURL.Host)
	if err != nil {
		host = parsedURL.Host
		// 如果没有指定端口，默认为 ""
		if parsedURL.Scheme == "http" {
			port = "80"
		} else if parsedURL.Scheme == "https" {
			port = "443"
		} else if parsedURL.Scheme == "" {
			host = parsedURL.Host
			port = "443"
		} else {
			return "", "", err
		}
	}
	return host, port, nil

}
