package app

import (
	utils "codepass/util"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ServiceModel struct {
	Ip   string
	Port string
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
	passFile := utils.RunDir(fmt.Sprintf("/.codepass/workspaces/%s/pass", name))
	if len(entry.Ipv4) > 0 {
		entry.Ip = entry.Ipv4[0]
	}
	entry.Create = strings.TrimSpace(utils.ReadFile(createFile))
	entry.Pass = strings.TrimSpace(utils.ReadFile(passFile))
	entry.Domain, entry.Url = instanceDomain(name)
	return entry
}

// 获取工作区域名
func instanceDomain(name string) (string, string) {
	domainFile := utils.RunDir("/.codepass/nginx/cert/domain")
	if utils.IsFile(domainFile) {
		domainAddr := fmt.Sprintf("%s-code.%s", name, strings.TrimSpace(utils.ReadFile(domainFile)))
		_, httpsPort := utils.GetProtsConfig()
		if httpsPort == "443" {
			return domainAddr, fmt.Sprintf("https://%s", domainAddr)
		} else {
			return domainAddr, fmt.Sprintf("https://%s:%s", domainAddr, httpsPort)
		}
	}
	return "", ""
}

// 更新工作区域名
func updateDomain() error {
	domainFile := utils.RunDir("/.codepass/nginx/cert/domain")
	domainAddr := ""
	if utils.IsFile(domainFile) {
		domainAddr = strings.TrimSpace(utils.ReadFile(domainFile))
	}
	if domainAddr == "" {
		return errors.New("no domain")
	}
	err := utils.WriteFile(utils.RunDir("/.codepass/nginx/lua/upsteam.lua"), utils.TemplateContent(utils.NginxUpsteamLua, map[string]any{}))
	if err != nil {
		return err
	}
	var list []string
	list = append(list, utils.TemplateContent(utils.NginxDefaultConf, map[string]any{
		"MAIN_DOMAIN":  domainAddr,
		"SERVICE_PORT": ServiceConf.Port,
	}))
	for _, entry := range workspacesList() {
		if entry.Ip != "" && entry.Domain != "" {
			list = append(list, utils.TemplateContent(utils.NginxDomainConf, map[string]any{
				"DOMAIN": entry.Domain,
				"IP":     entry.Ip,
			}))
		}
	}
	err = utils.WriteFile(utils.RunDir("/.codepass/nginx/conf.d/default.conf"), strings.Join(list, "\n"))
	if err != nil {
		return err
	}
	//
	httpPort, httpsPort := utils.GetProtsConfig()
	err = utils.WriteFile(utils.RunDir("/.codepass/docker/docker-compose.yml"), utils.TemplateContent(utils.DockerComposeContent, map[string]any{
		"HTTP_PORT":  httpPort,
		"HTTPS_PORT": httpsPort,
	}))
	if err != nil {
		return err
	}
	_, _ = utils.Cmd("-c", fmt.Sprintf("docker-compose -f %s down", utils.RunDir("/.codepass/docker/docker-compose.yml")))
	_, err = utils.Cmd("-c", fmt.Sprintf("docker-compose -f %s up -d --remove-orphans", utils.RunDir("/.codepass/docker/docker-compose.yml")))
	if err != nil {
		return err
	}
	return nil
}
