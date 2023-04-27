package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// Create 创建实例
func (model *MultipassModel) Create(c *gin.Context) {
	// 参数校验
	name := c.Query("name")
	pass := c.Query("pass")
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例名称不能为空",
		})
		return
	}
	if !utils.Test(name, "^[a-zA-Z][a-zA-Z0-9_]*$") {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例名称只允许字母开头，数字、字母、下划线组成",
		})
		return
	}
	if pass == "" {
		pass = utils.GenerateString(16)
	}
	if !utils.Test(pass, "^[a-zA-Z0-9_]*$") {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例密码只允许数字、字母、下划线组成",
		})
		return
	}
	// 检测实例是否已存在
	dirPath := fmt.Sprintf("/tmp/.codepass/instances/%s", name)
	if utils.IsDir(dirPath) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例已存在",
		})
		return
	}
	// 生成创建脚本
	cmdFile := fmt.Sprintf("/tmp/.codepass/instances/%s/launch.sh", name)
	logFile := fmt.Sprintf("/tmp/.codepass/instances/%s/launch.log", name)
	err := utils.WriteFile(cmdFile, utils.FromTemplateContent(utils.CreateExecContent, gin.H{
		"NAME": name,
		"PASS": pass,
		"CPUS": "",
		"MEM":  "",
		"DISK": "",
	}))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "创建实例失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	// 执行创建脚本
	go func() {
		_, _ = utils.Cmd("-c", fmt.Sprintf("chmod +x %s", cmdFile))
		_, _ = utils.Cmd("-c", fmt.Sprintf("/bin/sh %s > %s", cmdFile, logFile))
	}()
	//
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "创建实例成功",
		"data": gin.H{
			"name": name,
			"pass": pass,
		},
	})
}

// CreateLog 查看创建日志
func (model *MultipassModel) CreateLog(c *gin.Context) {
	name := c.Query("name")
	tail, _ := strconv.Atoi(c.Query("tail"))
	if tail <= 0 {
		tail = 200
	}
	if tail > 10000 {
		tail = 10000
	}
	logFile := fmt.Sprintf("/tmp/.codepass/instances/%s/launch.log", name)
	statusFile := fmt.Sprintf("/tmp/.codepass/instances/%s/status", name)
	if !utils.IsFile(logFile) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "日志文件不存在",
		})
		return
	}
	logContent, _ := utils.Cmd("-c", fmt.Sprintf("tail -%d %s", tail, logFile))
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "读取成功",
		"data": gin.H{
			"status": strings.TrimSpace(utils.ReadFile(statusFile)),
			"log":    strings.TrimSpace(logContent),
		},
	})
}

// Info 查看实例信息
func (model *MultipassModel) Info(c *gin.Context) {
	name := c.Query("name")
	dirPath := fmt.Sprintf("/tmp/.codepass/instances/%s", name)
	if !utils.IsDir(dirPath) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例不存在",
		})
		return
	}
	result, err := utils.Cmd("-c", fmt.Sprintf("multipass info %s --format json", name))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "获取失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	var data infoModel
	if err = json.Unmarshal([]byte(result), &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "解析失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	statusFile := fmt.Sprintf("/tmp/.codepass/instances/%s/status", name)
	passFile := fmt.Sprintf("/tmp/.codepass/instances/%s/pass", name)
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "获取成功",
		"data": gin.H{
			"status": strings.TrimSpace(utils.ReadFile(statusFile)),
			"pass":   strings.TrimSpace(utils.ReadFile(passFile)),
			"info":   data.Info[name],
		},
	})
}

// Delete 删除实例
func (model *MultipassModel) Delete(c *gin.Context) {
	name := c.Query("name")
	dirPath := fmt.Sprintf("/tmp/.codepass/instances/%s", name)
	if !utils.IsDir(dirPath) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例不存在",
		})
		return
	}
	_, _ = utils.Cmd("-c", fmt.Sprintf("multipass umount %s", name))            // 取消目录挂载
	_, _ = utils.Cmd("-c", fmt.Sprintf("rm -rf %s", dirPath))                   // 删除实例目录
	_, err := utils.Cmd("-c", fmt.Sprintf("multipass delete --purge %s", name)) // 删除实例
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "实例删除失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "实例删除成功",
	})
}