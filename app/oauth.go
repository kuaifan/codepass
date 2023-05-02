package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

var (
	clientId     = "cbd2d3097323fbbdafaa"
	clientSecret = "ca27bc159978a87c7570a15ea760e39663af4fb8"
)

// OAuthGithub Github授权
func (model *ServiceModel) OAuthGithub(c *gin.Context) bool {
	urlPath := c.Request.URL.Path
	_, homeUrl := instanceDomain("")
	if strings.HasPrefix(urlPath, "/oauth/redirect") {
		code := c.Query("code")
		githubToken, err := githubGetToken(clientId, clientSecret, code)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("授权失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		if githubToken.AccessToken == "" {
			c.String(http.StatusBadRequest, "授权失败：bad_verification_code")
			return true
		}
		userInfo, err := githubGetUserInfo(githubToken.AccessToken)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("获取用户信息失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		apiToken := utils.GenerateString(32)
		userData, err := json.Marshal(&userInfo)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("解析用户信息失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		err = utils.WriteFile(utils.RunDir(fmt.Sprintf("/.codepass/users/%s", apiToken)), string(userData))
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("AccessToken 保存失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		c.SetCookie("apiToken", apiToken, 0, "/", homeUrl, false, true)
		c.Redirect(http.StatusMovedPermanently, homeUrl)
		return true
	}
	var apiFile string
	apiToken, _ := c.Cookie("apiToken")
	userInfo := &githubUserModel{}
	if apiToken != "" {
		apiFile = utils.RunDir(fmt.Sprintf("/.codepass/users/%s", apiToken))
		userData := utils.ReadFile(apiFile)
		if err := json.Unmarshal([]byte(userData), userInfo); err != nil {
			apiToken = ""
		}
		if userInfo.ID == 0 {
			apiToken = ""
		}
	}
	if apiToken == "" {
		location := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s/oauth/redirect", clientId, homeUrl)
		c.Redirect(http.StatusMovedPermanently, location)
		return true
	}
	if strings.HasPrefix(urlPath, "/oauth/user") {
		userInfo.AccessToken = "" // 防止前端泄露token
		c.JSON(http.StatusOK, gin.H{
			"ret":  1,
			"msg":  "获取成功",
			"data": userInfo,
		})
		return true
	} else if strings.HasPrefix(urlPath, "/oauth/logout") {
		if utils.IsFile(apiFile) {
			_ = os.Remove(apiFile)
		}
		c.String(http.StatusOK, "退出成功")
		return true
	}
	return false
}
