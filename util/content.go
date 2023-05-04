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

# 保存状态
CREATE() {
	echo "\n[$1]"
	echo "$1" > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/create
}

# 判断状态
JUDGEA() {
	if [ 0 -eq $? ]; then
		echo "$1 完成"
	else
		echo "$1 失败"
		CREATE "Failed"
		rm -f $CmdPath
		exit 1
	fi
}
JUDGEB() {
	local desc="$1"
	local state=$(multipass exec {{.NAME}} -- sh -c 'cat /tmp/.code-judge')
	if [ "$state" = "success" ]; then
		echo "$desc 完成"
	else
		echo "$desc 失败"
		CREATE "Failed"
		rm -f $CmdPath
		exit 1
	fi
}

# 准备工作
CREATE "Preparing"
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/code-server
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/workspace
cat > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/info.yaml <<-EOF
owner_name: {{.OWNER_NAME}}
repos_owner: {{.REPOS_OWNER}}
repos_name: {{.REPOS_NAME}}
repos_url: {{.REPOS_URL}}
EOF
cat > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/init.yaml <<-EOF
runcmd:
  - curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
  - sudo apt-get install -y nodejs
  - sudo curl -sSL https://get.daocloud.io/docker | sh
  - sudo systemctl start docker
EOF
cat > {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/code-server/config.yaml <<-EOF
bind-addr: 0.0.0.0:55123
auth: password
password: {{.PASSWORD}}
cert: false
EOF

# 启动虚拟机
CREATE "Launching"
start="multipass launch focal --name {{.NAME}}"
start="$start --cloud-init {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config/init.yaml"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/config:~/.config"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/{{.NAME}}/workspace:~/workspace"
[ -n "{{.CPUS}}" ] && start="$start --cpus {{.CPUS}}"
[ -n "{{.DISK}}" ] && start="$start --disk {{.DISK}}"
[ -n "{{.MEMORY}}" ] && start="$start --memory {{.MEMORY}}"
$start
multipass info {{.NAME}} &> /dev/null
JUDGEA "Launch"

# 安装 code-server
CREATE "Installing"
multipass exec {{.NAME}} -- sh <<-EOE
curl -fsSL https://code-server.dev/install.sh | sh
sudo echo ".card-box > .header {display:none}" >> /usr/lib/code-server/src/browser/pages/login.css
sudo rm -f /usr/lib/code-server/src/browser/media/pwa-icon-192.png
sudo rm -f /usr/lib/code-server/src/browser/media/pwa-icon-512.png
sudo rm -f /usr/lib/code-server/src/browser/media/favicon-dark-support.svg
sudo rm -f /usr/lib/code-server/src/browser/media/favicon.ico
sudo ln -s \${HOME}/workspace /workspace
code-server --version &> /dev/null
if [ 0 -eq \$? ]; then
	echo "success" > /tmp/.code-judge
else
	echo "error" > /tmp/.code-judge
fi
EOE
JUDGEB "Install"

# Cloning
CREATE "Cloning"
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
CREATE "Starting"
multipass exec {{.NAME}} -- sudo sh <<-EOE
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
