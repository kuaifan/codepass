package utils

import (
	"bytes"
	"fmt"
	"text/template"
)

// InstallExecContent 安装脚本
const InstallExecContent = string(`
#!/bin/bash

# 字体颜色
Green="\033[32m"
Red="\033[31m"
GreenBG="\033[42;37m"
RedBG="\033[41;37m"
Font="\033[0m"

# 脚本路径
CmdPath=$0

# 输出 OK
OK() {
    echo "${Green}[OK]${Font} ${GreenBG}$1${Font}"
}

# 输出 ERR
ERR() {
    echo "${Red}[ERR]${Font} ${RedBG}$1${Font}"
    rm -f $CmdPath
    exit 1
}

# 检测 docker
check_docker() {
    docker --version &> /dev/null
    if [ $? -ne  0 ]; then
        echo "安装docker环境..."
        curl -sSL https://get.daocloud.io/docker | sh
        OK "Docker环境安装完成"
    fi
	if [ "$(uname)" == "Linux" ]; then
    	systemctl start docker
        if [[ 0 -ne $? ]]; then
            ERR "Docker 启动 失败"
        fi
	fi
    #
    docker-compose --version &> /dev/null
    if [ $? -ne  0 ]; then
        echo "安装docker-compose..."
        curl -s -L "https://get.daocloud.io/docker/compose/releases/download/v2.17.3/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
        ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
        OK "Docker-compose安装完成"
		if [ "$(uname)" == "Linux" ]; then
        	service docker restart
		fi
    fi
}

# 检测 multipass
check_multipass() {
	multipass version &> /dev/null
	if [ $? -ne  0 ]; then
		echo "开始安装 multipass..."
		if [ "$(uname)" == "Darwin" ]; then
			brew install --cask multipass
		else
			sudo snap install multipass
		fi
		multipass version &> /dev/null
		if [ $? -ne  0 ]; then
			ERR "安装失败 multipass"
			exit 1
		fi
		OK "安装完成 multipass"
	fi
	
    if [ ! -f "/tmp/.codepass/install/focal" ]; then
        echo "拉取镜像 ubuntu:focal..."
        multipass launch focal --name codepass-testing
        local list=$(multipass list | grep "codepass-testing")
        if [ -z "$list" ]; then
            ERR "拉取失败 ubuntu:focal"
            exit 1
        fi
        multipass delete --purge codepass-testing
        echo "1" > /tmp/.codepass/install/focal
        OK "拉取完成 ubuntu:focal"
    fi
}

# 运行脚本
check_docker
check_multipass
OK "安装完成"

# 删除脚本
rm -f $CmdPath
`)

// CreateExecContent 创建脚本
const CreateExecContent = string(`
#!/bin/bash

# 脚本路径
CmdPath=$0

# 全局变量
NAME="{{.NAME}}"
PASS="{{.PASS}}"
CPUS="{{.CPUS}}"
MEM="{{.MEM}}"
DISK="{{.DISK}}"

# 保存状态
STATE() {
	echo "$1"
    echo "$1" > /tmp/.codepass/instances/$NAME/state
}

# 启动虚拟机
STATE "Launching"
start="multipass launch focal --name $NAME"
[ -n "$CPUS" ] && start="$start --cpus $CPUS"
[ -n "$MEM" ] && start="$start --mem $MEM"
[ -n "$DISK" ] && start="$start --disk $DISK"
$start

# 挂载目录
STATE "Mounting"
mkdir -p /tmp/.codepass/instances/$NAME/work
multipass mount /tmp/.codepass/instances/$NAME/work $NAME:/work

# 安装 code-server
STATE "Installing"
multipass exec $NAME -- sh -c 'curl -fsSL https://code-server.dev/install.sh | sh'

# 初始化配置
STATE "Configuring"
multipass exec $NAME -- sh <<-EOE
mkdir -p ~/.config/code-server
cat > ~/.config/code-server/config.yaml <<-EOF
bind-addr: 0.0.0.0:51234
auth: password
password: '$PASS'
cert: false
EOF
EOE
multipass exec $NAME -- sudo sh -c 'echo ".card-box > .header {display:none}" >> /usr/lib/code-server/src/browser/pages/login.css'

# 启动 code-server
STATE "Starting"
multipass exec $NAME -- sudo sh -c 'systemctl enable --now code-server@ubuntu'

# 保存密码
echo "$PASS" > /tmp/.codepass/instances/$NAME/pass

# 输出成功
STATE "Success"

# 删除脚本
rm -f $CmdPath
`)

// FromTemplateContent 从模板中获取内容
func FromTemplateContent(templateContent string, envMap map[string]interface{}) string {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			PrintError(fmt.Sprintf("模板分析失败: %s", err))
		}
	}()
	if err != nil {
		panic(1)
	}
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}
