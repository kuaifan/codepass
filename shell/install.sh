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
    if [ ! -f "{{.WORK_PATH}}/install/default" ]; then
        echo "下载默认镜像 ubuntu:20.04 ..."
        multipass launch 20.04 --name codepass-default
        local list=$(multipass list | grep "codepass-default")
        if [ -z "$list" ]; then
            ERR "下载镜像失败 ubuntu:20.04"
        fi
        multipass delete --purge codepass-default
        echo "20.04" > {{.WORK_PATH}}/install/default
        OK "下载镜像完成 ubuntu:20.04"
    fi
}

# 运行脚本
check_multipass
OK "环境安装完成"

# 删除脚本
rm -f $CmdPath