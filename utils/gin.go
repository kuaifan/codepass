package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

// GinInput Gin获取参数（优先POST、取Query）
func GinInput(c *gin.Context, key string) string {
	if c.PostForm(key) != "" {
		return strings.TrimSpace(c.PostForm(key))
	}
	return strings.TrimSpace(c.Query(key))
}

// GinGetCookie Gin获取Cookie
func GinGetCookie(c *gin.Context, name string) string {
	value, _ := c.Cookie(name)
	return value
}

// GinSetCookie Gin设置Cookie
func GinSetCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetCookie(name, value, maxAge, "/", "", false, false)
}

// GinRemoveCookie Gin删除Cookie
func GinRemoveCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "", false, false)
}

// GinResult 返回结果
func GinResult(c *gin.Context, code int, content string, values ...any) {
	c.Header("Expires", "-1")
	c.Header("Cache-Control", "no-cache")
	c.Header("Pragma", "no-cache")
	var data any
	if len(values) == 1 {
		data = values[0]
	} else if len(values) == 0 {
		data = gin.H{}
	} else {
		data = values
	}
	//
	if strings.Contains(c.GetHeader("Accept"), "application/json") {
		// 接口返回
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  content,
			"data": data,
		})
	} else {
		// 页面返回
		if code == http.StatusMovedPermanently {
			c.Redirect(code, content)
		} else {
			c.HTML(http.StatusOK, "/web/dist/index.html", gin.H{
				"CODE": code,
				"MSG":  url.QueryEscape(content),
			})
		}
	}
}
