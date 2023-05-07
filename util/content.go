package utils

import (
	"bytes"
	"fmt"
	"strings"
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
	echo -e "${Green}[OK]${Font} ${GreenBG}$1${Font}"
}

# 输出 ERR
ERR() {
	echo -e "${Red}[ERR]${Font} ${RedBG}$1${Font}"
	rm -f $CmdPath
	exit 1
}

# 检测 multipass
check_multipass() {
	multipass version > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "开始安装 multipass ..."
		if [ "$(uname -s)" = "Darwin" ]; then
			brew cask install multipass
		elif [ "$(uname -s)" = "Linux" ]; then
			sudo snap install multipass
		else
			ERR "不支持的操作系统"
		fi
		multipass version > /dev/null 2>&1
		if [ $? -ne 0 ]; then
			ERR "multipass 安装失败"
		else
			OK "multipass 安装成功"
		fi
		echo "正在启动 multipassd ..."
		waitnum=0
		while [ -z "${multipassd}" ]; do
			waitnum=$((waitnum+1))
			multipassd=$(multipass version | grep "multipassd" | awk '{print $2}')
			if [ ${waitnum} -gt 60 ]; then
				ERR "multipassd 启动超时"
			fi
			sleep 1
		done
		OK "multipassd 启动成功"
	fi
	if [ ! -f "{{.RUN_PATH}}/.codepass/install/default" ]; then
		echo "下载默认镜像 ubuntu:20.04 ..."
		multipass launch 20.04 --name codepass-default
		local list=$(multipass list | grep "codepass-default")
		if [ -z "$list" ]; then
			ERR "下载镜像失败 ubuntu:20.04"
		fi
		multipass delete --purge codepass-default
		echo "20.04" > {{.RUN_PATH}}/.codepass/install/default
		OK "下载镜像完成 ubuntu:20.04"
	fi
}

# 运行脚本
check_multipass
OK "环境安装完成"

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
# {{.PASSWORD}}
# {{.PROXY_DOMAIN}}
# {{.PROXY_URI}}

# {{.OWNER_NAME}}
# {{.REPOS_OWNER}}
# {{.REPOS_NAME}}
# {{.REPOS_URL}}
# {{.CLONE_CMD}}

# {{.CPUS}}
# {{.DISK}}
# {{.MEMORY}}
# {{.IMAGE}}

# {{.CREATED_AT}}

# 保存状态
STATUS() {
	echo ""
	echo "[$1]"
	echo "$1" > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/status
}

# 判断状态
JUDGEA() {
	if [ 0 -eq $? ]; then
		echo "$1 完成"
	else
		echo "$1 失败"
		STATUS "Failed"
		rm -f $CmdPath
		exit 1
	fi
}
JUDGEB() {
	local desc="$1"
	local state=$(multipass exec {{.NAME}} -- sudo sh -c 'cat /tmp/.code-judge && rm -f /tmp/.code-judge')
	if [ "$state" = "success" ]; then
		echo "$desc 完成"
	else
		echo "$desc 失败"
		STATUS "Failed"
		rm -f $CmdPath
		exit 1
	fi
}

# 准备工作
STATUS "Preparing"
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/code-server
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/workspace
cat > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/info.yaml <<-EOF
owner_name: {{.OWNER_NAME}}
repos_owner: {{.REPOS_OWNER}}
repos_name: {{.REPOS_NAME}}
repos_url: {{.REPOS_URL}}
created_at: {{.CREATED_AT}}
image: {{.IMAGE}}
EOF
cat > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/init.yaml <<-EOF
runcmd:
  - curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
  - sudo apt-get install -y nodejs
  - sudo curl -sSL https://get.docker.com/ | sh
  - sudo systemctl start docker
EOF
cat > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/code-server/config.yaml <<-EOF
app-name: CodePass
disable-workspace-trust: true
bind-addr: 0.0.0.0:55123
auth: password
password: {{.PASSWORD}}
cert: false
EOF

# 启动虚拟机
STATUS "Launching"
start="multipass launch {{.IMAGE}} --name {{.NAME}}"
start="$start --cloud-init {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/init.yaml"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config:~/.config"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/workspace:~/workspace"
[ -n "{{.CPUS}}" ] && start="$start --cpus {{.CPUS}}"
[ -n "{{.DISK}}" ] && start="$start --disk {{.DISK}}"
[ -n "{{.MEMORY}}" ] && start="$start --memory {{.MEMORY}}"
$start
multipass info {{.NAME}} > /dev/null 2>&1
JUDGEA "Launch"

# 安装 code-server
STATUS "Installing"
multipass exec {{.NAME}} -- sh <<-EOE
curl -fsSL https://code-server.dev/install.sh | sh
sudo sh -c 'echo ".card-box > .header {display:none}" >> /usr/lib/code-server/src/browser/pages/login.css'
sudo rm -f /usr/lib/code-server/src/browser/media/pwa-icon-192.png
sudo rm -f /usr/lib/code-server/src/browser/media/pwa-icon-512.png
sudo rm -f /usr/lib/code-server/src/browser/media/favicon-dark-support.svg
sudo rm -f /usr/lib/code-server/src/browser/media/favicon.ico
sudo ln -s \${HOME}/workspace /workspace
code-server --version > /dev/null 2>&1
if [ 0 -eq \$? ]; then
	echo "success" > /tmp/.code-judge
else
	echo "error" > /tmp/.code-judge
fi
EOE
JUDGEB "Install"

# Cloning
STATUS "Cloning"
multipass exec {{.NAME}} -- sh <<-EOE
{{.CLONE_CMD}} /workspace/{{.REPOS_NAME}}
if [ -d "/workspace/{{.REPOS_NAME}}/.git/" ]; then
	echo "success" > /tmp/.code-judge
else
	echo "error" > /tmp/.code-judge
fi
EOE
JUDGEB "Clone"

# 启动 code-server
STATUS "Starting"
multipass exec {{.NAME}} -- sudo sh <<-EOE
systemctl set-environment CODE_PASS_IMAGE={{.IMAGE}}
systemctl set-environment PROXY_DOMAIN={{.PROXY_DOMAIN}}
systemctl set-environment VSCODE_PROXY_URI={{.PROXY_URI}}
systemctl set-environment DEFAULT_WORKSPACE=/workspace/{{.REPOS_NAME}}
systemctl enable --now code-server@ubuntu
if [ 0 -eq \$? ]; then
	echo "success" > /tmp/.code-judge
else
	echo "error" > /tmp/.code-judge
fi
EOE
JUDGEB "Start"

# 输出成功
STATUS "Success"

# 删除脚本
rm -f $CmdPath
`)

// OperationContent 操作内容
const OperationContent = string(`
#!/bin/bash

# 脚本路径
CmdPath=$0

# 全局变量
# {{.NAME}}
# {{.PROXY_DOMAIN}}
# {{.PROXY_URI}}

# {{.REPOS_NAME}}

# {{.IMAGE}}
# {{.OPERATION}}

# 保存状态
STATUS() {
	echo "$1" > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/status
}

SERVER() {
	multipass exec {{.NAME}} -- sudo sh <<-EOE
systemctl set-environment CODE_PASS_IMAGE={{.IMAGE}}
systemctl set-environment PROXY_DOMAIN={{.PROXY_DOMAIN}}
systemctl set-environment VSCODE_PROXY_URI={{.PROXY_URI}}
systemctl set-environment DEFAULT_WORKSPACE=/workspace/{{.REPOS_NAME}}
systemctl restart --now code-server@ubuntu
EOE
}

start() {
	STATUS "Starting"
	echo "Starting..."
	multipass start {{.NAME}}
	if [ 0 -eq $? ]; then
		SERVER
		STATUS "Success"
		echo "Started"
	else
		STATUS "Error"
		echo "Start failed"
	fi
}

stop() {
	STATUS "Stopping"
	echo "Stopping..."
	multipass stop {{.NAME}}
	if [ 0 -eq $? ]; then
		STATUS "Success"
		echo "Stopped"
	else
		STATUS "Error"
		echo "Stop failed"
	fi
}

restart() {
	STATUS "Restarting"
	echo "Restarting..."
	multipass restart {{.NAME}}
	if [ 0 -eq $? ]; then
		SERVER
		STATUS "Success"
		echo "Restarted"
	else
		STATUS "Error"
		echo "Restart failed"
	fi
}

delete() {
	STATUS "Deleting"
	echo "Deleting..."
	multipass delete --purge {{.NAME}}
	if [ 0 -eq $? ]; then
		STATUS "Success"
		echo "Deleted"
	else
		STATUS "Error"
		echo "Delete failed"
	fi
	if [ -d "{{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}" ]; then
		rm -rf {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}
	fi
}

# 执行命令
{{.OPERATION}}

# 删除脚本
rm -f $CmdPath
`)

// TemplateContent 从模板中获取内容
func TemplateContent(templateContent string, envMap map[string]interface{}) string {
	templateContent = strings.ReplaceAll(templateContent, "\t", "    ")
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
