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