package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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
	_, homePage := instanceDomain("")
	if strings.HasPrefix(urlPath, "/oauth/redirect") {
		code := c.Query("code")
		githubToken, err := githubGetToken(clientId, clientSecret, code)
		if err != nil {
			utils.GinResponse200(c, 0, fmt.Sprintf("授权失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		if githubToken.AccessToken == "" {
			utils.GinResponse200(c, 0, "授权失败：bad_verification_code")
			return true
		}
		userInfo, err := githubGetUserInfo(githubToken.AccessToken)
		if err != nil {
			utils.GinResponse200(c, 0, fmt.Sprintf("获取用户信息失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		apiToken := utils.GenerateString(32)
		userData, err := json.Marshal(&userInfo)
		if err != nil {
			utils.GinResponse200(c, 0, fmt.Sprintf("解析用户信息失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		err = utils.WriteFile(utils.RunDir(fmt.Sprintf("/.codepass/users/%s", apiToken)), string(userData))
		if err != nil {
			utils.GinResponse200(c, 0, fmt.Sprintf("AccessToken 保存失败：%s", removeCriticalInformation(err.Error())))
			return true
		}
		c.SetCookie("apiToken", apiToken, 0, "/", homePage, false, true)
		utils.GinResponse301(c, homePage)
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
		location := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s/oauth/redirect", clientId, homePage)
		utils.GinResponse301(c, location)
		return true
	}
	if strings.HasPrefix(urlPath, "/oauth/user") {
		userInfo.AccessToken = "" // 清空防止前端泄露AccessToken
		utils.GinResponse200(c, 1, "获取成功", userInfo)
		return true
	} else if strings.HasPrefix(urlPath, "/oauth/logout") {
		if utils.IsFile(apiFile) {
			_ = os.Remove(apiFile)
		}
		utils.GinResponse200(c, 1, "退出成功")
		return true
	}
	return false
}
