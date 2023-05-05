#!/bin/bash

WORKDIR=$(cd `dirname $0`; pwd)
HOMEDIR=$(cd ~; pwd)/codepass

# 字体颜色
Green="\033[32m"
Red="\033[31m"
GreenBG="\033[42;37m"
RedBG="\033[41;37m"
Font="\033[0m"

# 输出 OK
OK() {
	echo -e "${Green}[OK]${Font} ${GreenBG}$1${Font}"
}

# 输出 ERR
ERR() {
	echo -e "${Red}[ERR]${Font} ${RedBG}$1${Font}"
}

# 判断输出
JUDGE() {
    if [[ 0 -eq $? ]]; then
        OK "$1 完成"
        sleep 1
    else
        ERR "$1 失败"
        exit 1
    fi
}

# 检测系统架构
PLATFORM=$(uname -s)
FRAMEWORK=$(uname -m)
if [[ ${PLATFORM} != "Darwin" ]] && [[ ${PLATFORM} != "Linux" ]]; then
    ERR "不支持的系统：${PLATFORM}，仅支持 Linux 和 Mac"
    exit 1
fi
if [[ ${FRAMEWORK} == "aarch64" ]]; then
    FRAMEWORK="arm64"
fi
if [[ ${FRAMEWORK} == "x86_64" ]]; then
    FRAMEWORK="amd64"
fi
if [[ ${FRAMEWORK} != "amd64" ]] && [[ ${FRAMEWORK} != "arm64" ]]; then
    ERR "不支持的架构：${FRAMEWORK}，仅支持 amd64 和 arm64"
    exit 1
fi

# 获取github最新发布releases版本
VERSION=$(curl -s https://api.github.com/repos/kuaifan/codepass/releases/latest | grep 'tag_name' | cut -d\" -f4)
if [[ -z ${VERSION} ]]; then
    ERR "获取最新版本失败"
    exit 1
fi

# 输入必须信息
while [ -z "$DOMAIN" ]; do
    read -rp "请输入你的主域名(例如: abc.com):" DOMAIN
done
while [ -z "$CLIENT_ID" ]; do
    read -rp "请输入 GitHub Client ID:" CLIENT_ID
done
while [ -z "$CLIENT_SECRET" ]; do
    read -rp "请输入 GitHub Client secrets:" CLIENT_SECRET
done

# 检测域名
DOIP=$(ping "${DOMAIN}" -c 1 | sed '1{s/[^(]*(//;s/).*//;q}')
OK "正在获取 公网ip 信息，请耐心等待"
LOIP=$(curl https://api-ipv4.ip.sb/ip)
echo -e "域名dns解析IP：${DOIP}"
echo -e "本机IP: ${LOIP}"
sleep 2
if [[ $(echo "${LOIP}" | tr '.' '+' | bc) -eq $(echo "${DOIP}" | tr '.' '+' | bc) ]]; then
    OK "域名DNS解析IP 与 本机IP 匹配"
    sleep 2
else
    ERR "域名DNS解析IP 与 本机IP 不匹配，是否继续安装？（Y/n）"
    read -r CONTINUEIN
    [[ -z ${CONTINUEIN} ]] && CONTINUEIN="Y"
    case $CONTINUEIN in
    [yY][eE][sS] | [yY])
        echo -e "${GreenBG}继续安装${Font}"
        sleep 2
        ;;
    *)
        echo -e "${RedBG}安装终止${Font}"
        exit 2
        ;;
    esac
fi

# 安装证书
curl https://get.acme.sh | sh
JUDGE "安装 SSL 证书生成脚本"
alias acme.sh=~/.acme.sh/acme.sh

# 生成证书目录
SSLPATH=$HOMEDIR/ssl/$DOMAIN
mkdir -p $SSLPATH

# 申请证书
acme.sh --register-account -m admin@admin.com
acme.sh --issue -d $DOMAIN -d *.$DOMAIN --standalone
acme.sh --installcert -d $DOMAIN -d *.$DOMAIN --key-file $SSLPATH/site.key --fullchain-file $SSLPATH/site.crt
if [ ! -f $SSLPATH/site.key ]; then
    ERR "证书申请失败"
    exit 1
fi

# 安装 codepass
PKGPATH=$HOMEDIR/pkg
mkdir -p $PKGPATH
curl -s -L "https://github.com/kuaifan/codepass/releases/download/${VERSION}/codepass_${PLATFORM}_${FRAMEWORK}.tar.gz" -o $PKGPATH/package.tar.gz
tar -zxf $PKGPATH/package.tar.gz -C $PKGPATH
if [ ! -f $PKGPATH/codepass/codepass ]; then
    ERR "codepass 安装失败"
    exit 1
fi
mv $PKGPATH/codepass/* $HOMEDIR
rm -r $PKGPATH
chmod +x $HOMEDIR/codepass

# 配置 codepass
$HOMEDIR/codepass config \
--host $DOMAIN \
--port 443 \
--ssl-cert $SSLPATH/site.crt \
--ssl-key $SSLPATH/site.key \
--client-id $CLIENT_ID \
--client-secret $CLIENT_SECRET \
--path $HOMEDIR/config.yaml

echo ""
echo ""
OK "[codepass 程序安装完成]"
OK "执行程序: $HOMEDIR/codepass"
OK "配置文件: $HOMEDIR/config.yaml"
echo ""
OK "1.安装环境: $HOMEDIR/codepass install"
OK "2.启动服务: $HOMEDIR/codepass service --conf $HOMEDIR/config.yaml"
