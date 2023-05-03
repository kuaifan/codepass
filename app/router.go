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

// OAuth 授权
func (model *ServiceModel) OAuth(c *gin.Context) {
	urlPath := c.Request.URL.Path
	// 静态资源
	if strings.HasPrefix(urlPath, "/assets") {
		c.File(fmt.Sprintf("./web/dist%s", urlPath))
		return
	}
	// 退出登录
	if strings.HasPrefix(urlPath, "/oauth/logout") {
		apiToken, _ := c.Cookie("api_token")
		if apiToken != "" {
			apiFile := utils.RunDir(fmt.Sprintf("/.codepass/users/%s", apiToken))
			if utils.IsFile(apiFile) {
				_ = os.Remove(apiFile)
			}
		}
		c.SetCookie("api_token", "", -1, "/", c.Request.Host, false, false)
		utils.GinResult(c, http.StatusOK, "退出成功")
		return
	}
	// 授权回应
	if strings.HasPrefix(urlPath, "/oauth/redirect") {
		code := c.Query("code")
		githubToken, err := githubGetToken(clientId, clientSecret, code)
		if err != nil {
			utils.GinResult(c, http.StatusOK, fmt.Sprintf("授权失败：%s", removeCriticalInformation(err.Error())))
			return
		}
		if githubToken.AccessToken == "" {
			utils.GinResult(c, http.StatusOK, "授权失败：bad_verification_code")
			return
		}
		userInfo, err := githubGetUserInfo(githubToken.AccessToken)
		if err != nil {
			utils.GinResult(c, http.StatusOK, fmt.Sprintf("获取用户信息失败：%s", removeCriticalInformation(err.Error())))
			return
		}
		apiToken := utils.GenerateString(32)
		userData, err := json.Marshal(&userInfo)
		if err != nil {
			utils.GinResult(c, http.StatusOK, fmt.Sprintf("解析用户信息失败：%s", removeCriticalInformation(err.Error())))
			return
		}
		err = utils.WriteFile(utils.RunDir(fmt.Sprintf("/.codepass/users/%s", apiToken)), string(userData))
		if err != nil {
			utils.GinResult(c, http.StatusOK, fmt.Sprintf("AccessToken 保存失败：%s", removeCriticalInformation(err.Error())))
			return
		}
		c.SetCookie("api_token", apiToken, 0, "/", c.Request.Host, false, false)
		utils.GinResult(c, http.StatusMovedPermanently, "/")
		return
	}
	// 读取身份
	apiFile := ""
	userInfo := &githubUserModel{}
	apiToken, _ := c.Cookie("api_token")
	if apiToken != "" {
		apiFile = utils.RunDir(fmt.Sprintf("/.codepass/users/%s", apiToken))
		userData := utils.ReadFile(apiFile)
		_ = json.Unmarshal([]byte(userData), userInfo)
		if userInfo.ID == 0 {
			apiToken = ""
		}
	}
	// 发起授权
	if apiToken == "" {
		_, homePage := instanceDomain("")
		var items []map[string]any
		items = append(items, gin.H{
			"type":  "github",
			"label": "使用GitHub登录",
			"href":  fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s/oauth/redirect", clientId, homePage),
		})
		content, _ := json.Marshal(&items)
		utils.GinResult(c, http.StatusUnauthorized, string(content))
		return
	}
	// 身份信息
	if strings.HasPrefix(urlPath, "/api/user/info") {
		userInfo.AccessToken = "" // 清空防止前端泄露AccessToken
		utils.GinResult(c, http.StatusOK, "获取成功", userInfo)
		return
	}
	// 工作区接口
	if strings.HasPrefix(urlPath, "/api/workspaces/create/log") {
		ServiceConf.WorkspacesCreateLog(c)
		return
	}
	if strings.HasPrefix(urlPath, "/api/workspaces/create") {
		ServiceConf.WorkspacesCreate(c)
		return
	}
	if strings.HasPrefix(urlPath, "/api/workspaces/list") {
		ServiceConf.WorkspacesList(c)
		return
	}
	if strings.HasPrefix(urlPath, "/api/workspaces/info") {
		ServiceConf.WorkspacesInfo(c)
		return
	}
	if strings.HasPrefix(urlPath, "/api/workspaces/delete") {
		ServiceConf.WorkspacesDelete(c)
		return
	}
	// 页面输出
	c.SetCookie("result_code", "", -1, "/", c.Request.Host, false, false)
	c.SetCookie("result_msg", "", -1, "/", c.Request.Host, false, false)
	c.File("./web/dist/index.html")
}
