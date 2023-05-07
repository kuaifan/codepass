package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// WorkspacesCreate 创建工作区（post）
func (model *ServiceModel) WorkspacesCreate(c *gin.Context) {
	// 参数校验
	var (
		repos    = utils.GinInput(c, "repos")
		password = utils.GinInput(c, "password")
		cpus     = utils.GinInput(c, "cpus")
		disk     = utils.GinInput(c, "disk")
		memory   = utils.GinInput(c, "memory")
		image    = utils.GinInput(c, "image")
	)
	if repos == "" {
		utils.GinResult(c, http.StatusBadRequest, "储存库地址不能为空")
		return
	}
	if !utils.Test(repos, "^https*://") {
		utils.GinResult(c, http.StatusBadRequest, "储存库地址格式错误")
		return
	}
	re := "^https*://github.com/([^\\/\\.]+)/([^\\/\\.]+)(\\.git)?\\/?$"
	name := ""
	reposOwner := ""
	reposName := ""
	if utils.Test(repos, re) {
		reg := regexp.MustCompile(re)
		match := reg.FindStringSubmatch(repos)
		reposOwner = strings.ToLower(match[1])
		reposName = strings.ToLower(match[2])
		name = strings.ToLower(fmt.Sprintf("%s-%s-%s", reposOwner, reposName, utils.GenerateString(8)))
	} else {
		utils.GinResult(c, http.StatusBadRequest, "暂不支持此储存库地址")
	}
	if password == "" {
		password = utils.GenerateString(32)
	}
	if !utils.Test(password, "^[a-zA-Z0-9_]*$") {
		utils.GinResult(c, http.StatusBadRequest, "工作区密码只允许数字、字母、下划线组成")
		return
	}
	if cpus != "" && !utils.Test(cpus, "^\\d+$") {
		utils.GinResult(c, http.StatusBadRequest, "CPU只能是存数字")
		return
	}
	if disk != "" && utils.Test(disk, "^\\d+$") {
		disk = fmt.Sprintf("%sGB", disk)
	}
	if memory != "" && utils.Test(memory, "^\\d+$") {
		memory = fmt.Sprintf("%sGB", memory)
	}
	if image == "" {
		image = "20.04"
	} else if !utils.InArray(image, []string{"18.04", "20.04", "22.04", "22.10"}) {
		utils.GinResult(c, http.StatusBadRequest, "请选择有效的系统版本")
		return
	}
	// 检测工作区是否已存在
	dirPath := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s", name))
	if utils.IsDir(dirPath) {
		utils.GinResult(c, http.StatusBadRequest, "工作区已存在")
		return
	}
	// 端口代理地址
	_, url := instanceDomain(name)
	proxyRegexp := regexp.MustCompile(`^(https*://)`)
	proxyDomain := proxyRegexp.ReplaceAllString(url, "")
	proxyUri := proxyRegexp.ReplaceAllString(url, "$1{{port}}-")
	// 生成创建脚本
	cmdFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create", name))
	logFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/logs", name))
	err := utils.WriteFile(cmdFile, utils.TemplateContent(utils.CreateExecContent, map[string]any{
		"NAME":         name,
		"PASSWORD":     password,
		"PROXY_DOMAIN": proxyDomain,
		"PROXY_URI":    proxyUri,

		"OWNER_NAME":  ServiceConf.GithubUserInfo.Login,
		"REPOS_OWNER": reposOwner,
		"REPOS_NAME":  reposName,
		"REPOS_URL":   repos,
		"CLONE_CMD":   fmt.Sprintf("git clone https://oauth2:%s@github.com/%s/%s.git", ServiceConf.GithubUserInfo.AccessToken, reposOwner, reposName),

		"CPUS":   cpus,
		"DISK":   disk,
		"MEMORY": memory,
		"IMAGE":  image,

		"CREATED_AT": utils.FormatYmdHis(time.Now()),
	}))
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, "创建工作区失败", gin.H{
			"err": err.Error(),
		})
		return
	}
	// 执行创建脚本
	go func() {
		_, _ = utils.Cmd("-c", fmt.Sprintf("chmod +x %s", cmdFile))
		_, _ = utils.Cmd("-c", fmt.Sprintf("/bin/bash %s > %s 2>&1", cmdFile, logFile))
		UpdateProxy()
	}()
	//
	utils.GinResult(c, http.StatusOK, "创建工作区成功", gin.H{
		"name":     name,
		"password": password,
	})
}

// WorkspacesLog 查看创建日志
func (model *ServiceModel) WorkspacesLog(c *gin.Context) {
	name := c.Query("name")
	tail, _ := strconv.Atoi(c.Query("tail"))
	if tail <= 0 {
		tail = 200
	}
	if tail > 10000 {
		tail = 10000
	}
	logFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/logs", name))
	statusFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/status", name))
	if !utils.IsFile(logFile) {
		utils.GinResult(c, http.StatusBadRequest, "日志文件不存在")
		return
	}
	logContent, _ := utils.Cmd("-c", fmt.Sprintf("tail -%d %s", tail, logFile))
	utils.GinResult(c, http.StatusOK, "读取成功", gin.H{
		"status": strings.TrimSpace(utils.ReadFile(statusFile)),
		"log":    strings.TrimSpace(logContent),
	})
}

// WorkspacesList 获取工作区列表
func (model *ServiceModel) WorkspacesList(c *gin.Context) {
	var list []*instanceModel
	for _, entry := range workspacesList() {
		if entry.OwnerName == ServiceConf.GithubUserInfo.Login {
			list = append(list, entry)
		}
	}
	utils.GinResult(c, http.StatusOK, "获取成功", gin.H{
		"list": list,
	})
}

// WorkspacesInfo 查看工作区信息
func (model *ServiceModel) WorkspacesInfo(c *gin.Context) {
	name := c.Query("name")
	format := c.Query("format")
	//
	_, err := instanceInfo(name, true)
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, err.Error())
		return
	}
	var result string
	if format == "hard" {
		result, err = utils.Cmd("-c", fmt.Sprintf("multipass get local.%s.cpus && multipass get local.%s.disk && multipass get local.%s.memory", name, name, name))
	} else if format == "json" {
		result, err = utils.Cmd("-c", fmt.Sprintf("multipass info %s --format json", name))
	} else {
		result, err = utils.Cmd("-c", fmt.Sprintf("multipass info %s", name))
	}
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, "获取失败", gin.H{
			"err": err.Error(),
		})
		return
	}
	var info any
	if format == "hard" {
		// 按换行分割 result
		res := strings.Split(result, "\n")
		utils.GinResult(c, http.StatusOK, "获取成功", gin.H{
			"cpus":   res[0],
			"disk":   res[1],
			"memory": res[2],
		})
		return
	} else if format == "json" {
		var data infoModel
		if err = json.Unmarshal([]byte(result), &data); err != nil {
			utils.GinResult(c, http.StatusBadRequest, "解析失败", gin.H{
				"err": err.Error(),
			})
			return
		}
		info = data.Info[name]
	} else {
		info = result
	}
	statusFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/status", name))
	viper.SetConfigFile(utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/config/code-server/config.yaml", name)))
	_ = viper.ReadInConfig()
	password := viper.GetString("password")
	viper.SetConfigFile(utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/config/info.yaml", name)))
	_ = viper.ReadInConfig()
	var (
		ownerName  = viper.GetString("owner_name")
		reposOwner = viper.GetString("repos_owner")
		reposName  = viper.GetString("repos_name")
		reposUrl   = viper.GetString("repos_url")
	)
	utils.GinResult(c, http.StatusOK, "获取成功", gin.H{
		"status":      strings.TrimSpace(utils.ReadFile(statusFile)),
		"password":    password,
		"owner_name":  ownerName,
		"repos_owner": reposOwner,
		"repos_name":  reposName,
		"repos_url":   reposUrl,
		"info":        info,
	})
}

// WorkspacesModify 修改工作区
func (model *ServiceModel) WorkspacesModify(c *gin.Context) {
	// 参数校验
	var (
		name   = utils.GinInput(c, "name")
		cpus   = utils.GinInput(c, "cpus")
		disk   = utils.GinInput(c, "disk")
		memory = utils.GinInput(c, "memory")
	)
	info, err := instanceInfo(name, true)
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, err.Error())
		return
	}
	if info.State != "Stopped" {
		utils.GinResult(c, http.StatusBadRequest, "请先停止工作区")
		return
	}
	if cpus != "" && !utils.Test(cpus, "^\\d+$") {
		utils.GinResult(c, http.StatusBadRequest, "CPU只能是存数字")
		return
	}
	if disk != "" && utils.Test(disk, "^\\d+$") {
		disk = fmt.Sprintf("%sGB", disk)
	}
	if memory != "" && utils.Test(memory, "^\\d+$") {
		memory = fmt.Sprintf("%sGB", memory)
	}
	result, err := utils.Cmd("-c", fmt.Sprintf("multipass set local.%s.cpus=%s && multipass set local.%s.disk=%s && multipass set local.%s.memory=%s", name, cpus, name, disk, name, memory))
	if err != nil {
		if result == "" {
			result = "修改失败"
		}
		utils.GinResult(c, http.StatusBadRequest, result, gin.H{
			"err": err.Error(),
		})
		return
	}
	utils.GinResult(c, http.StatusOK, "修改成功", gin.H{
		"result": result,
	})
}

// WorkspacesOperation 操作工作区（启动、停止、重启、删除）
func (model *ServiceModel) WorkspacesOperation(c *gin.Context) {
	name := c.Query("name")
	operation := c.Query("operation")
	status := map[string]string{
		"start":   "Starting",
		"stop":    "Stoping",
		"restart": "Restarting",
		"delete":  "Deleting",
	}[operation]
	if status == "" {
		utils.GinResult(c, http.StatusBadRequest, "操作类型错误")
		return
	}
	info, err := instanceInfo(name, true)
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, err.Error())
		return
	}
	// 端口代理地址
	_, url := instanceDomain(name)
	proxyRegexp := regexp.MustCompile(`^(https*://)`)
	proxyDomain := proxyRegexp.ReplaceAllString(url, "")
	proxyUri := proxyRegexp.ReplaceAllString(url, "$1{{port}}-")
	cmdFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/operation", name))
	logFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/logs", name))
	err = utils.WriteFile(cmdFile, utils.TemplateContent(utils.OperationContent, map[string]any{
		"NAME":         name,
		"PROXY_DOMAIN": proxyDomain,
		"PROXY_URI":    proxyUri,

		"REPOS_NAME": info.ReposName,

		"IMAGE":     info.Image,
		"OPERATION": operation,
	}))
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, "创建操作失败", gin.H{
			"err": err.Error(),
		})
		return
	}
	// 执行操作脚本
	go func() {
		_, _ = utils.Cmd("-c", fmt.Sprintf("chmod +x %s", cmdFile))
		_, _ = utils.Cmd("-c", fmt.Sprintf("/bin/bash %s >> %s 2>&1", cmdFile, logFile))
		if operation == "delete" {
			UpdateProxy()
		}
	}()
	utils.GinResult(c, http.StatusOK, "操作成功", gin.H{
		"status": status,
	})
}
