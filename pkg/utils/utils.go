package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// 执行外部命令
func ExecCmd(cmdStr string) {
	cmd := exec.Command(cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	//log.Println(string(out))
}

func ExecuteCmd(cmdStr string) string {
	log.Println("执行的命令:", cmdStr)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(output)
}

func IsReadable(f string) bool {
	err := syscall.Access(f, syscall.O_RDONLY)
	if err != nil {
		return false
	} else {
		return true
	}
}

// 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false
	}
	//文件存在，判断是否可读
	return IsReadable(path)
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但不是不存在的错误，所以把这个错误原封不动的返回
}
