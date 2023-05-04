# CodePass

## 介绍

CodePass 是一个基于 Go 语言开发的，用于快速搭建在线代码开发平台。

## 安装启动

```bash
curl -s -L "https://github.com/kuaifan/codepass/releases/download/v0.0.1/codepass_$(uname -s)_$(uname -m).tar.gz" -o ./codepass.tar.gz
tar -zxvf codepass.tar.gz && rm -f codepass.tar.gz && cd codepass
chmod +x ./codepass
./codepass install 
./codepass service --conf ./config.yaml
```

