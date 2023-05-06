package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// OAuth 授权
func (model *ServiceModel) OAuth(c *gin.Context) {
	urlPath := c.Request.URL.Path
	// 静态资源
	if strings.HasPrefix(urlPath, "/assets") {
		c.File(utils.RunDir(fmt.Sprintf("/.codepass/web/dist%s", urlPath)))
		return
	}
	// 退出登录
	if strings.HasPrefix(urlPath, "/oauth/logout") {
		userToken := utils.GinGetCookie(c, "result_token")
		if userToken != "" {
			apiFile := utils.RunDir(fmt.Sprintf("/.codepass/users/%s", userToken))
			if utils.IsFile(apiFile) {
				_ = os.Remove(apiFile)
			}
		}
		utils.GinRemoveCookie(c, "result_token")
		utils.GinResult(c, http.StatusOK, "退出成功")
		return
	}
	// 授权回应
	if strings.HasPrefix(urlPath, "/oauth/redirect") {
		code := c.Query("code")
		githubToken, err := githubGetToken(ServiceConf.GithubClientId, ServiceConf.GithubClientSecret, code)
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
		userToken := utils.GenerateString(32)
		userData, err := json.Marshal(&userInfo)
		if err != nil {
			utils.GinResult(c, http.StatusOK, fmt.Sprintf("解析用户信息失败：%s", removeCriticalInformation(err.Error())))
			return
		}
		err = utils.WriteFile(utils.RunDir(fmt.Sprintf("/.codepass/users/%s", userToken)), string(userData))
		if err != nil {
			utils.GinResult(c, http.StatusOK, fmt.Sprintf("AccessToken 保存失败：%s", removeCriticalInformation(err.Error())))
			return
		}
		utils.GinSetCookie(c, "result_token", userToken)
		utils.GinResult(c, http.StatusMovedPermanently, "/")
		return
	}
	// 读取身份
	apiFile := ""
	userInfo := &githubUserModel{}
	userToken := utils.GinGetCookie(c, "result_token")
	if userToken != "" {
		apiFile = utils.RunDir(fmt.Sprintf("/.codepass/users/%s", userToken))
		userData := utils.ReadFile(apiFile)
		_ = json.Unmarshal([]byte(userData), userInfo)
		if userInfo.ID == 0 {
			userToken = ""
		}
	}
	// 发起授权
	if userToken == "" {
		_, homePage := instanceDomain("")
		redirectUri := url.QueryEscape(fmt.Sprintf("%s/oauth/redirect", homePage))
		var items []map[string]any
		items = append(items, gin.H{
			"type":  "github",
			"label": "使用GitHub登录",
			"href":  fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", ServiceConf.GithubClientId, redirectUri),
		})
		content, _ := json.Marshal(&items)
		utils.GinResult(c, http.StatusUnauthorized, string(content))
		return
	}
	ServiceConf.GithubUserInfo = userInfo
	// 身份信息
	if strings.HasPrefix(urlPath, "/api/user/info") {
		userInfo.AccessToken = "" // 清空防止前端泄露AccessToken
		utils.GinResult(c, http.StatusOK, "获取成功", userInfo)
		return
	}
	// 用户存储库
	if strings.HasPrefix(urlPath, "/api/user/repos") {
		ServiceConf.UserRepositories(c)
		return
	}
	// 工作区接口
	if strings.HasPrefix(urlPath, "/api/workspaces/create") {
		ServiceConf.WorkspacesCreate(c) // post
		return
	}
	if strings.HasPrefix(urlPath, "/api/workspaces/log") {
		ServiceConf.WorkspacesLog(c)
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
	if strings.HasPrefix(urlPath, "/api/workspaces/modify") {
		ServiceConf.WorkspacesModify(c)
		return
	}
	if strings.HasPrefix(urlPath, "/api/workspaces/delete") {
		ServiceConf.WorkspacesDelete(c)
		return
	}
	// 页面输出
	utils.GinRemoveCookie(c, "result_code")
	utils.GinRemoveCookie(c, "result_msg")
	c.HTML(http.StatusOK, "/web/dist/index.html", gin.H{})
}
