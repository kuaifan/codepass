package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type ServiceModel struct {
	Host string
	Port string
	Crt  string
	Key  string
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
	State   string   `json:"state"`

	Create string `json:"create"`
	Pass   string `json:"pass"`

	Domain string `json:"domain"`
	Url    string `json:"url"`
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
	for _, entry := range data.List {
		instanceBase(entry)
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
			for _, exists := range data.List {
				if name == exists.Name {
					exist = true
					break
				}
			}
			if !exist {
				data.List = append(data.List, instanceBase(&instanceModel{
					Name: name,
				}))
			}
		}
	}
	return data.List
}

// 获取工作区基本信息
func instanceBase(entry *instanceModel) *instanceModel {
	name := entry.Name
	createFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/create", name))
	viper.SetConfigFile(utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/config/code-server/config.yaml", name)))
	_ = viper.ReadInConfig()
	if len(entry.Ipv4) > 0 {
		entry.Ip = entry.Ipv4[0]
	}
	entry.Create = strings.TrimSpace(utils.ReadFile(createFile))
	entry.Pass = viper.GetString("password")
	entry.Domain, entry.Url = instanceDomain(name)
	return entry
}

// 获取[工作区]域名
func instanceDomain(name string) (string, string) {
	domainAddr := ServiceConf.Host
	if name != "" {
		domainAddr = fmt.Sprintf("%s-code.%s", name, domainAddr)
	}
	if ServiceConf.Port == "443" {
		return domainAddr, fmt.Sprintf("https://%s", domainAddr)
	} else {
		return domainAddr, fmt.Sprintf("https://%s:%s", domainAddr, ServiceConf.Port)
	}
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
