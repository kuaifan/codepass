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

DOMAIN="{{.DOMAIN}}"
KEY="{{.KEY}}"
CRT="{{.CRT}}"

# 保存状态
CREATE() {
	echo "$1"
	echo "$1" > /tmp/.codepass/instances/$NAME/create
}

# 域名/证书
if [ -n "$DOMAIN" ]; then
	mkdir -p /tmp/.codepass/instances/$NAME/certs
	cat > /tmp/.codepass/instances/$NAME/certs/$DOMAIN.key <<-EOF
$KEY
EOF
	cat > /tmp/.codepass/instances/$NAME/certs/$DOMAIN.crt <<-EOF
$CRT
EOF
	echo "$DOMAIN" > /tmp/.codepass/instances/$NAME/certs/domain
fi

# 启动虚拟机
CREATE "Launching"
start="multipass launch focal --name $NAME"
[ -n "$CPUS" ] && start="$start --cpus $CPUS"
[ -n "$MEM" ] && start="$start --mem $MEM"
[ -n "$DISK" ] && start="$start --disk $DISK"
$start

# 挂载目录
CREATE "Mounting"
mkdir -p /tmp/.codepass/instances/$NAME/share
multipass mount /tmp/.codepass/instances/$NAME/share $NAME:/share

# 安装 code-server
CREATE "Installing"
multipass exec $NAME -- sh -c 'curl -fsSL https://code-server.dev/install.sh | sh'

# 初始化配置
CREATE "Configuring"
multipass exec $NAME -- sh <<-EOE
mkdir -p ~/.config/code-server
cat > ~/.config/code-server/config.yaml <<-EOF
bind-addr: 0.0.0.0:51234
auth: password
password: $PASS
cert: false
EOF
EOE
multipass exec $NAME -- sudo sh -c 'echo ".card-box > .header {display:none}" >> /usr/lib/code-server/src/browser/pages/login.css'

# 启动 code-server
CREATE "Starting"
multipass exec $NAME -- sudo sh -c 'systemctl enable --now code-server@ubuntu'

# 保存密码
echo "$PASS" > /tmp/.codepass/instances/$NAME/pass

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

// NginxDomainConf nginx 域名配置
const NginxDomainConf = string(`
# {{.NAME}}
# {{.IP}}
# {{.DOMAIN}}

upstream server_{{.NAME}} {
	server {{.IP}}:51234 weight=5 max_fails=3 fail_timeout=30s;
	keepalive 16;
}

server {
	listen 80;
	listen 443 ssl http2;
	server_name {{.DOMAIN}};

	if ($server_port !~ 443){
		rewrite ^(/.*)$ https://$host$1 permanent;
	}

	ssl_certificate /instances/{{.NAME}}/certs/{{.DOMAIN}}.crt;
	ssl_certificate_key /instances/{{.NAME}}/certs/{{.DOMAIN}}.key;
	
	ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
	ssl_ciphers EECDH+CHACHA20:EECDH+CHACHA20-draft:EECDH+AES128:RSA+AES128:EECDH+AES256:RSA+AES256:EECDH+3DES:RSA+3DES:!MD5;
	ssl_prefer_server_ciphers on;
	ssl_session_cache shared:SSL:10m;
	ssl_session_timeout 10m;

	location / {
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
		proxy_pass http://server_{{.NAME}};
	}
}
`)

// DockerComposeContent docker-compose 配置
const DockerComposeContent = string(`
version: '3'

services:
  nginx:
    image: "nginx:alpine"
    ports:
      - "1180:80"
      - "11443:443"
    volumes:
      - /tmp/.codepass/nginx/default.conf:/etc/nginx/conf.d/default.conf
      - /tmp/.codepass/instances:/instances
    restart: unless-stopped
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
