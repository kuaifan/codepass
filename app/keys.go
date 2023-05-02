package app

import (
	utils "codepass/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"net/http"
)

// KeysInfo KEY详情
func (model *ServiceModel) KeysInfo(c *gin.Context) {
	filePath := utils.RunDir("/.codepass/keys/default")
	if !utils.IsFile(filePath) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "未设置",
		})
		return
	}
	result := utils.ReadFile(filePath)
	var data keyModel
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "解析失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret":  1,
		"msg":  "获取成功",
		"data": data,
	})
}

// KeysSave 保存KEY
func (model *ServiceModel) KeysSave(c *gin.Context) {
	var (
		title = utils.GinInput(c, "title")
		key   = utils.GinInput(c, "key")
	)
	if key == "" {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "KEY不能为空",
		})
		return
	}
	_, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(key))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "KEY解析失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	if title == "" {
		title = comment
	}
	// 保存KEY
	data := keyModel{
		Title: title,
		Key:   key,
	}
	result, err := json.Marshal(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "解析失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	err = utils.WriteFile(utils.RunDir("/.codepass/keys/default"), string(result))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "保存失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	//
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "保存成功",
	})
}
