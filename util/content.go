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
	#
	docker pull openresty/openresty:alpine
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
PROXY_DOMAIN="{{.PROXY_DOMAIN}}"
PROXY_URI="{{.PROXY_URI}}"

CPUS="{{.CPUS}}"
DISK="{{.DISK}}"
MEMORY="{{.MEMORY}}"

# 保存状态
CREATE() {
	echo "$1"
	echo "$1" > {{.RUN_PATH}}/.codepass/workspaces/$NAME/create
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
EOF
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/$NAME/config
mkdir -p {{.RUN_PATH}}/.codepass/workspaces/$NAME/workspace

# 启动虚拟机
CREATE "Launching"
start="multipass launch focal --name $NAME"
start="$start --cloud-init {{.RUN_PATH}}/.codepass/workspaces/config.yaml"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/$NAME/config:/config"
start="$start --mount {{.RUN_PATH}}/.codepass/workspaces/$NAME/workspace:/workspace"
[ -n "$CPUS" ] && start="$start --cpus $CPUS"
[ -n "$DISK" ] && start="$start --disk $DISK"
[ -n "$MEMORY" ] && start="$start --memory $MEMORY"
$start

# 安装 code-server
CREATE "Installing"
multipass exec $NAME -- sh -c 'curl -fsSL https://code-server.dev/install.sh | sh'

# 初始化 code-server 配置
CREATE "Configuring"
multipass exec $NAME -- sh <<-EOE
mkdir -p ~/.config/code-server
mkdir -p ~/wwwroot
cat > ~/.config/code-server/config.yaml <<-EOF
bind-addr: 0.0.0.0:55123
auth: password
password: $PASS
cert: false
EOF
EOE

# 优化 code-server 页面资源
multipass exec $NAME -- sudo sh <<-EOE
echo ".card-box > .header {display:none}" >> /usr/lib/code-server/src/browser/pages/login.css
rm -f /usr/lib/code-server/src/browser/media/pwa-icon-192.png
rm -f /usr/lib/code-server/src/browser/media/pwa-icon-512.png
rm -f /usr/lib/code-server/src/browser/media/favicon-dark-support.svg
rm -f /usr/lib/code-server/src/browser/media/favicon.ico
EOE

# 启动 code-server
CREATE "Starting"
multipass exec $NAME -- sudo sh <<-EOE
systemctl set-environment PROXY_DOMAIN=$PROXY_DOMAIN
systemctl set-environment VSCODE_PROXY_URI=$PROXY_URI
systemctl set-environment DEFAULT_WORKSPACE=/workspace
systemctl enable --now code-server@ubuntu
if [ 0 -eq \$? ]; then
	echo "success" > /tmp/.code-server
else
	echo "error" > /tmp/.code-server
fi
EOE
server=$(multipass exec $NAME -- sh -c 'cat /tmp/.code-server')
if [ "$server" != "success" ]; then
	CREATE "Failed"
	rm -f $CmdPath
	exit 1
fi

# 保存密码
echo "$PASS" > {{.RUN_PATH}}/.codepass/workspaces/$NAME/pass

# 输出成功
CREATE "Success"

# 删除脚本
rm -f $CmdPath
`)

// NginxDefaultConf nginx 默认配置
const NginxDefaultConf = string(`
map $http_upgrade $connection_upgrade {
	default upgrade;
	'' close;
}
map $http_host $this_host {
	"" $host;
	default $http_host;
}
map $http_x_forwarded_proto $the_scheme {
	 default $http_x_forwarded_proto;
	 "" $scheme;
}
map $http_x_forwarded_host $the_host {
	default $http_x_forwarded_host;
	"" $this_host;
}
`)

// NginxUpsteamLua nginx lua配置
const NginxUpsteamLua = string(`
local ip = ngx.arg[1];
local host = ngx.var.host;
local pattern = "(%d+)-"
local port = string.match(host, pattern)

if port then
    return string.format("%s:%s", ip, port)
else
    return string.format("%s:55123", ip)
end
`)

// NginxDomainConf nginx 域名配置
const NginxDomainConf = string(`
# {{.IP}}
# {{.DOMAIN}}

server {
	listen 80;
	listen 443 ssl http2;
	server_name ~^(\d+)?{{.DOMAIN}}$;

	if ($server_port !~ 443){
		rewrite ^(/.*)$ https://$host$1 permanent;
	}

	ssl_certificate /etc/nginx/cert/crt;
	ssl_certificate_key /etc/nginx/cert/key;
	
	ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
	ssl_ciphers EECDH+CHACHA20:EECDH+CHACHA20-draft:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
	ssl_prefer_server_ciphers on;
	ssl_session_cache shared:SSL:10m;
	ssl_session_timeout 10m;

	location / {
		set_by_lua_file $ups /etc/nginx/lua/upsteam.lua {{.IP}};
		proxy_http_version 1.1;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Real-PORT $remote_port;
		proxy_set_header X-Forwarded-Host $the_host;
		proxy_set_header X-Forwarded-Proto $the_scheme;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header Host $http_host;
		proxy_set_header Scheme $scheme;
		proxy_set_header Server-Protocol $server_protocol;
		proxy_set_header Server-Name $server_name;
		proxy_set_header Server-Addr $server_addr;
		proxy_set_header Server-Port $server_port;
		proxy_set_header Upgrade $http_upgrade;
		proxy_set_header Connection $connection_upgrade;
		proxy_pass http://$ups;
	}
}
`)

// DockerComposeContent docker-compose 配置
const DockerComposeContent = string(`
# {{.HTTP_PORT}}
# {{.HTTPS_PORT}}

version: '3'

services:
  nginx:
    image: "openresty/openresty:alpine"
    ports:
      - "{{.HTTP_PORT}}:80"
      - "{{.HTTPS_PORT}}:443"
    volumes:
      - {{.RUN_PATH}}/.codepass/nginx/cert:/etc/nginx/cert
      - {{.RUN_PATH}}/.codepass/nginx/lua:/etc/nginx/lua
      - {{.RUN_PATH}}/.codepass/nginx/conf.d:/etc/nginx/conf.d
    restart: unless-stopped
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
