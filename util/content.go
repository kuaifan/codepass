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
	
	if [ ! -f "{{.RUN_PATH}}/.codepass/install/focal" ]; then
		echo "拉取镜像 ubuntu:focal..."
		multipass launch focal --name codepass-testing
		local list=$(multipass list | grep "codepass-testing")
		if [ -z "$list" ]; then
			ERR "拉取失败 ubuntu:focal"
			exit 1
		fi
		multipass delete --purge codepass-testing
		echo "1" > {{.RUN_PATH}}/.codepass/install/focal
		OK "拉取完成 ubuntu:focal"
	fi
}

# 运行脚本
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
# {{.NAME}}
# {{.PASS}}
# {{.PROXY_DOMAIN}}
# {{.PROXY_URI}}

# {{.OWNER_NAME}}
# {{.REPOS_OWNER}}
# {{.REPOS_NAME}}
# {{.CLONE_CMD}}

# {{.CPUS}}
# {{.DISK}}
# {{.MEMORY}}

# 保存状态
CREATE() {
	echo "$1"
	echo "$1" > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/create
}

# 准备工作
CREATE "Preparing"
cat > {{.RUN_PATH}}/.codepass/workspaces/config.yaml <<-EOF
#cloud-config
runcmd:
  - curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
  - sudo apt-get install -y nodejs
  - sudo curl -sSL https://get.daocloud.io/docker | sh
  - sudo systemctl start docker
  - {{.CLONE_CMD}}
EOF
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/workspace

# 启动虚拟机
CREATE "Launching"
start="multipass launch focal --name {{.NAME}}"
start="$start --cloud-init {{.RUN_PATH}}/.codepass/workspaces/config.yaml"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config:~/.config"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/workspace:~/workspace"
[ -n "{{.CPUS}}" ] && start="$start --cpus {{.CPUS}}"
[ -n "{{.DISK}}" ] && start="$start --disk {{.DISK}}"
[ -n "{{.MEMORY}}" ] && start="$start --memory {{.MEMORY}}"
$start

# 安装 code-server
CREATE "Installing"
multipass exec {{.NAME}} -- sh -c 'curl -fsSL https://code-server.dev/install.sh | sh'

# 初始化 code-server 配置
CREATE "Configuring"
multipass exec {{.NAME}} -- sh <<-EOE
mkdir -p ~/.config/code-server
cat > ~/.config/code-server/config.yaml <<-EOF
bind-addr: 0.0.0.0:55123
auth: password
password: {{.PASS}}
cert: false

owner-name: {{.OWNER_NAME}}
repos-owner: {{.REPOS_OWNER}}
repos-name: {{.REPOS_NAME}}
EOF
sudo ln -s \${HOME}/workspace /workspace
EOE

# 优化 code-server 页面资源
multipass exec {{.NAME}} -- sudo sh <<-EOE
echo ".card-box > .header {display:none}" >> /usr/lib/code-server/src/browser/pages/login.css
rm -f /usr/lib/code-server/src/browser/media/pwa-icon-192.png
rm -f /usr/lib/code-server/src/browser/media/pwa-icon-512.png
rm -f /usr/lib/code-server/src/browser/media/favicon-dark-support.svg
rm -f /usr/lib/code-server/src/browser/media/favicon.ico
EOE

# 启动 code-server
CREATE "Starting"
multipass exec {{.NAME}} -- sudo sh <<-EOE
systemctl set-environment PROXY_DOMAIN={{.PROXY_DOMAIN}}
systemctl set-environment VSCODE_PROXY_URI={{.PROXY_URI}}
systemctl set-environment DEFAULT_WORKSPACE=/workspace/{{.REPOS_NAME}}
systemctl enable --now code-server@ubuntu
if [ 0 -eq \$? ]; then
	echo "success" > /tmp/.code-server
else
	echo "error" > /tmp/.code-server
fi
EOE
server=$(multipass exec {{.NAME}} -- sh -c 'cat /tmp/.code-server')
if [ "$server" != "success" ]; then
	CREATE "Failed"
	rm -f $CmdPath
	exit 1
fi

# 输出成功
CREATE "Success"

# 删除脚本
rm -f $CmdPath
`)

// TemplateContent 从模板中获取内容
func TemplateContent(templateContent string, envMap map[string]interface{}) string {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			PrintError(fmt.Sprintf("模板分析失败: %s", err))
		}
	}()
	if err != nil {
		panic(1)
	}
	envMap["RUN_PATH"] = RunDir("")
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}
