package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (model *ServiceModel) UserRepositories(c *gin.Context, user *githubUserModel) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", user.Name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		utils.GinResult(c, http.StatusBadRequest, "无法创建请求", gin.H{
			"err": err.Error(),
		})
		return
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("%s <%s>", user.Name, user.AccessToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		utils.GinResult(c, http.StatusBadRequest, "发送请求失败", gin.H{
			"err": err.Error(),
		})
		return
	}
	var list = make([]githubReposSimplify, 0)
	if err = json.NewDecoder(res.Body).Decode(&list); err != nil {
		utils.GinResult(c, http.StatusBadRequest, "解析数据失败", gin.H{
			"err": err.Error(),
		})
		return
	}

	utils.GinResult(c, http.StatusOK, "获取成功", gin.H{
		"list": list,
	})
}
