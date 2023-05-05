# CodePass

## 介绍

CodePass 是一个基于 Go 语言开发的，用于快速搭建在线代码开发平台。

## 安装启动

#### 自动安装

```shell
bash <(curl -sSL https://raw.githubusercontent.com/kuaifan/codepass/main/install.sh)
```

#### 手动安装

```shell
# 请手动将 v0.0.1 替换为最新版本
curl -L "https://github.com/kuaifan/codepass/releases/download/v0.0.1/codepass_$(uname -s)_$(uname -m).tar.gz" -o ./codepass.tar.gz
tar -zxvf codepass.tar.gz && rm -f codepass.tar.gz && cd codepass
chmod +x ./codepass
./codepass install 
./codepass service --conf ./config.yaml
```