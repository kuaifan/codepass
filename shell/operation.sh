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
    echo "$1" > {{.WORK_PATH}}/workspaces/{{.NAME}}/status
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
    if [ -d "{{.WORK_PATH}}/workspaces/{{.NAME}}" ]; then
        rm -rf {{.WORK_PATH}}/workspaces/{{.NAME}}
    fi
}

# 执行命令
{{.OPERATION}}

# 删除脚本
rm -f $CmdPath