package app

import (
	utils "codepass/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CertsInfo 证书详情
func (model *ServiceModel) CertsInfo(c *gin.Context) {
	domainFile := utils.RunDir("/.codepass/nginx/cert/domain")
	keyFile := utils.RunDir("/.codepass/nginx/cert/key")
	crtFile := utils.RunDir("/.codepass/nginx/cert/crt")
	if !utils.IsFile(domainFile) {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "未设置",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "获取成功",
		"data": gin.H{
			"domain": utils.ReadFile(domainFile),
			"key":    utils.ReadFile(keyFile),
			"crt":    utils.ReadFile(crtFile),
		},
	})
}

// CertsSave 保存证书
func (model *ServiceModel) CertsSave(c *gin.Context) {
	var (
		domain = utils.GinInput(c, "domain")
		key    = utils.GinInput(c, "key")
		crt    = utils.GinInput(c, "crt")
	)
	if domain == "" {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "域名不能为空",
		})
		return
	}
	if key == "" {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "私钥不能为空",
		})
		return
	}
	if crt == "" {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "证书不能为空",
		})
		return
	}
	// 保存证书
	err1 := utils.WriteFile(utils.RunDir("/.codepass/nginx/cert/domain"), domain)
	err2 := utils.WriteFile(utils.RunDir("/.codepass/nginx/cert/key"), key)
	err3 := utils.WriteFile(utils.RunDir("/.codepass/nginx/cert/crt"), crt)
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "保存证书失败",
			"data": gin.H{
				"err1": err1.Error(),
				"err2": err2.Error(),
				"err3": err3.Error(),
			},
		})
		return
	}
	_ = UpdateDomain()
	//
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "保存证书成功",
	})
}
