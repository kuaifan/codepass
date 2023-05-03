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
)

// WorkspacesCreate 创建工作区
func (model *ServiceModel) WorkspacesCreate(c *gin.Context) {
	// 参数校验
	var (
		name   = c.Query("name")
		pass   = c.Query("pass")
		cpus   = c.Query("cpus")
		disk   = c.Query("disk")
		memory = c.Query("memory")
	)
	if name == "" {
		utils.GinResult(c, http.StatusBadRequest, "工作区名称不能为空")
		return
	}
	if !utils.Test(name, "^[a-zA-Z][a-zA-Z0-9_]*$") {
		utils.GinResult(c, http.StatusBadRequest, "工作区名称只允许字母开头，数字、字母、下划线组成")
		return
	}
	if pass == "" {
		pass = utils.GenerateString(16)
	}
	if !utils.Test(pass, "^[a-zA-Z0-9_]*$") {
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
	cmdFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create.sh", name))
	logFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create.log", name))
	err := utils.WriteFile(cmdFile, utils.TemplateContent(utils.CreateExecContent, map[string]any{
		"NAME":         name,
		"PASS":         pass,
		"PROXY_DOMAIN": proxyDomain,
		"PROXY_URI":    proxyUri,

		"CPUS":   cpus,
		"DISK":   disk,
		"MEMORY": memory,
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
		_, _ = utils.Cmd("-c", fmt.Sprintf("/bin/sh %s > %s 2>&1", cmdFile, logFile))
		UpdateProxy()
	}()
	//
	utils.GinResult(c, http.StatusOK, "创建工作区成功", gin.H{
		"name": name,
		"pass": pass,
	})
}

// WorkspacesLog 查看创建日志
func (model *ServiceModel) WorkspacesLog(c *gin.Context) {
	type_ := c.Query("type")
	name := c.Query("name")
	tail, _ := strconv.Atoi(c.Query("tail"))
	if tail <= 0 {
		tail = 200
	}
	if tail > 10000 {
		tail = 10000
	}
	if type_ == "create" {
		// 创建日志
		logFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create.log", name))
		createFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create", name))
		if !utils.IsFile(logFile) {
			utils.GinResult(c, http.StatusBadRequest, "日志文件不存在")
			return
		}
		logContent, _ := utils.Cmd("-c", fmt.Sprintf("tail -%d %s", tail, logFile))
		utils.GinResult(c, http.StatusOK, "读取成功", gin.H{
			"create": strings.TrimSpace(utils.ReadFile(createFile)),
			"log":    strings.TrimSpace(logContent),
		})
	} else {
		// 其他日志
		utils.GinResult(c, http.StatusBadRequest, "暂不支持")
	}
}

// WorkspacesList 获取工作区列表
func (model *ServiceModel) WorkspacesList(c *gin.Context) {
	list := workspacesList()
	if list == nil {
		utils.GinResult(c, http.StatusBadRequest, "暂无数据")
		return
	}
	utils.GinResult(c, http.StatusOK, "获取成功", gin.H{
		"list": list,
	})
}

// WorkspacesInfo 查看工作区信息
func (model *ServiceModel) WorkspacesInfo(c *gin.Context) {
	name := c.Query("name")
	format := c.Query("format")
	dirPath := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s", name))
	if !utils.IsDir(dirPath) {
		utils.GinResult(c, http.StatusBadRequest, "工作区不存在")
		return
	}
	var result string
	var err error
	if format == "json" {
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
	if format == "json" {
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
	createFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create", name))
	viper.SetConfigFile(utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/config/code-server/config.yaml", name)))
	_ = viper.ReadInConfig()
	utils.GinResult(c, http.StatusOK, "获取成功", gin.H{
		"create": strings.TrimSpace(utils.ReadFile(createFile)),
		"pass":   viper.GetString("password"),
		"info":   info,
	})
}

// WorkspacesDelete 删除工作区
func (model *ServiceModel) WorkspacesDelete(c *gin.Context) {
	name := c.Query("name")
	//
	_, err := utils.Cmd("-c", fmt.Sprintf("multipass info %s", name))
	if err == nil {
		_, err = utils.Cmd("-c", fmt.Sprintf("multipass delete --purge %s", name)) // 删除工作区
	}
	dirPath := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s", name))
	if utils.IsDir(dirPath) {
		_, _ = utils.Cmd("-c", fmt.Sprintf("rm -rf %s", dirPath)) // 删除工作区目录
	}
	UpdateProxy()
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, "工作区删除失败", gin.H{
			"err": err.Error(),
		})
		return
	}
	utils.GinResult(c, http.StatusOK, "工作区删除成功")
}
