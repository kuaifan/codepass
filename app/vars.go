package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
	"time"
)

type ServiceModel struct {
	Mode               string
	Conf               string
	Host               string
	Port               string
	SslCrt             string
	SslKey             string
	GithubClientId     string
	GithubClientSecret string
	GithubUserInfo     *githubUserModel
}

type ProxyModel struct {
	Name string
	Ip   string
}

type infoModel struct {
	Errors any            `json:"errors"`
	Info   map[string]any `json:"info"`
}

type listModel struct {
	List []*instanceModel `json:"list"`
}

type instanceModel struct {
	Ipv4    []string `json:"ipv4"`
	Ip      string   `json:"ip"`
	Name    string   `json:"name"`
	Release string   `json:"release"`
	State   string   `json:"state"` // 状态（实例）

	Status   string `json:"status"` // 状态（自定义）
	Password string `json:"password"`

	OwnerName  string `json:"owner_name"`
	ReposOwner string `json:"repos_owner"`
	ReposName  string `json:"repos_name"`
	ReposUrl   string `json:"repos_url"`
	CreatedAt  string `json:"created_at"`

	Domain string `json:"domain"`
	Url    string `json:"url"`
}

type githubTokenModel struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type githubUserModel struct {
	AccessToken string    `json:"access_token"`
	Login       string    `json:"login"`
	ID          int       `json:"id"`
	AvatarURL   string    `json:"avatar_url"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Company     string    `json:"company"`
	Blog        string    `json:"blog"`
	Location    string    `json:"location"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SaveAt      time.Time `json:"save_at"`
}

type githubReposSimplify struct {
	CreatedAt     time.Time `json:"created_at"`
	DefaultBranch string    `json:"default_branch"`
	Description   string    `json:"description"`
	Disabled      bool      `json:"disabled"`
	FullName      string    `json:"full_name"`
	HTMLURL       string    `json:"html_url"`
	ID            int       `json:"id"`
	Language      string    `json:"language"`
	Name          string    `json:"name"`
	NodeID        string    `json:"node_id"`
	Private       bool      `json:"private"`
	SSHURL        string    `json:"ssh_url"`
	Visibility    string    `json:"visibility"`
}

var (
	ServiceConf ServiceModel
	ProxyList   []ProxyModel
)

// 获取工作区列表
func workspacesList() []*instanceModel {
	result, err := utils.Cmd("-c", "multipass list --format json")
	if err != nil {
		return nil
	}
	var data listModel
	if err = json.Unmarshal([]byte(result), &data); err != nil {
		return nil
	}
	var list []*instanceModel
	for _, entry := range data.List {
		instanceBase(entry)
		if utils.Test(entry.Name, "^([^\\/\\.]+)-([^\\/\\.]+)-([^\\/\\.]+)$") {
			list = append(list, entry)
		}
	}
	//
	dirEntry, err := os.ReadDir(utils.RunDir("/.codepass/workspaces"))
	if err != nil {
		return nil
	}
	for _, entry := range dirEntry {
		if entry.IsDir() {
			name := entry.Name()
			exist := false
			for _, exists := range list {
				if name == exists.Name {
					exist = true
					break
				}
			}
			if !exist {
				list = append(list, instanceBase(&instanceModel{
					Name: name,
				}))
			}
		}
	}
	return list
}

// 获取工作区基本信息（通过工作区名称）
func instanceInfo(name string, oneself bool) (*instanceModel, error) {
	dirPath := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s", name))
	if !utils.IsDir(dirPath) {
		return nil, fmt.Errorf("工作区不存在")
	}
	info := instanceBase(&instanceModel{
		Name: name,
	})
	if oneself && info.OwnerName != ServiceConf.GithubUserInfo.Login {
		return nil, fmt.Errorf("无权操作：此工作区不属于你")
	}
	return info, nil
}

// 获取工作区基本信息
func instanceBase(entry *instanceModel) *instanceModel {
	name := entry.Name
	statusFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/status", name))
	viper.SetConfigFile(utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/config/code-server/config.yaml", name)))
	_ = viper.ReadInConfig()
	entry.Password = viper.GetString("password")
	viper.SetConfigFile(utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/config/info.yaml", name)))
	_ = viper.ReadInConfig()
	entry.OwnerName = viper.GetString("owner_name")
	entry.ReposOwner = viper.GetString("repos_owner")
	entry.ReposName = viper.GetString("repos_name")
	entry.ReposUrl = viper.GetString("repos_url")
	entry.CreatedAt = viper.GetString("created_at")
	if len(entry.Ipv4) > 0 {
		entry.Ip = entry.Ipv4[0]
	}
	entry.Status = strings.TrimSpace(utils.ReadFile(statusFile))
	entry.Domain, entry.Url = instanceDomain(name)
	return entry
}

// 获取[工作区]域名
func instanceDomain(name string) (string, string) {
	domainAddr := ServiceConf.Host
	if name != "" {
		domainAddr = fmt.Sprintf("%s.%s", name, domainAddr)
	}
	if ServiceConf.Port == "443" {
		return domainAddr, fmt.Sprintf("https://%s", domainAddr)
	} else {
		return domainAddr, fmt.Sprintf("https://%s:%s", domainAddr, ServiceConf.Port)
	}
}

// 获取 GitHub token
func githubGetToken(cid, sid, code string) (*githubTokenModel, error) {
	var url = fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", cid, sid, code)
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}
	token := &githubTokenModel{}
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	return token, nil
}

// 获取 GitHub UserInfo
func githubGetUserInfo(accessToken string) (*githubUserModel, error) {
	var userInfoUrl = "https://api.github.com/user"
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}
	userInfo := &githubUserModel{}
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	userInfo.AccessToken = accessToken
	userInfo.SaveAt = time.Now()
	return userInfo, nil
}

// 去除关键信息
func removeCriticalInformation(str string) string {
	if str == "" {
		return ""
	}
	str = strings.Replace(str, ServiceConf.GithubClientId, "********", -1)
	str = strings.Replace(str, ServiceConf.GithubClientSecret, "********", -1)
	return str
}

// UpdateProxy 更新代理地址
func UpdateProxy() {
	ProxyList = []ProxyModel{}
	for _, entry := range workspacesList() {
		ProxyList = append(ProxyList, ProxyModel{
			Name: entry.Name,
			Ip:   entry.Ip,
		})
	}
}
