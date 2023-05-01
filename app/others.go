package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// OthersDomainUpdate 更新域名
func (model *ServiceModel) OthersDomainUpdate(c *gin.Context) {
	err := UpdateDomain()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ret": 0,
			"msg": "更新失败",
			"data": gin.H{
				"err": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": 1,
		"msg": "更新成功",
	})
}
