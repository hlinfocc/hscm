#!/bin/bash

# Check if user is root
if [ $(id -u) != "0" ]; then
    echo "Error: You must be root to run this script, please use root to install"
    exit 1
fi

installCronJob(){
    # 检查 crontab 中是否存在定时任务
    existing_task=$(/usr/bin/crontab -l | grep "/usr/local/bin/hscm")

    if [ -z "$existing_task" ]; then
        (/usr/bin/crontab -l 2>/dev/null; echo "0 3 * * * /usr/local/bin/hscm") | /usr/bin/crontab -
        echo "新增[hscm]定时任务  ......成功"
    fi
}

if [ -f "./hscm" ];then
    \cp -f ./hscm /usr/local/bin/hscm
    installCronJob
    touch /etc/hscm.yml
    cat > /etc/hscm.yml <<EOF
webhook: https://oapi.dingtalk.com/robot/send?access_token=xxx
hosts:
  - xx.example.com

EOF
    echo "安装完成，请到/etc/hscm.yml配置webhook地址和需要监控的域名(将webhook和hosts换成实际的)"
fi
