package app

import (
	utils "codepass/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type MultipassModel struct {
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

	Domain   string `json:"domain"`
	KeyPath  string `json:"key_path"`
	CertPath string `json:"cert_path"`
}

var (
	MConf MultipassModel
)

// 获取实例列表
func instancesList() []*instanceModel {
	result, err := utils.Cmd("-c", "multipass list --format json")
	if err != nil {
		return nil
	}
	var data listModel
	if err = json.Unmarshal([]byte(result), &data); err != nil {
		return nil
	}
	for _, entry := range data.List {
		instancesBase(entry)
	}
	//
	dirEntry, err := os.ReadDir("/tmp/.codepass/instances")
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
				data.List = append(data.List, instancesBase(&instanceModel{
					Name: name,
				}))
			}
		}
	}
	return data.List
}

// 获取实例基本信息
func instancesBase(entry *instanceModel) *instanceModel {
	name := entry.Name
	domainAddr := ""
	domainKey := ""
	domainCert := ""
	domainFile := fmt.Sprintf("/tmp/.codepass/instances/%s/certs/domain", name)
	if utils.IsFile(domainFile) {
		domainAddr = strings.TrimSpace(utils.ReadFile(domainFile))
		domainKey = fmt.Sprintf("/tmp/.codepass/instances/%s/certs/%s.key", name, name)
		domainCert = fmt.Sprintf("/tmp/.codepass/instances/%s/certs/%s.crt", name, name)
	}
	createFile := fmt.Sprintf("/tmp/.codepass/instances/%s/create", name)
	passFile := fmt.Sprintf("/tmp/.codepass/instances/%s/pass", name)
	if len(entry.Ipv4) > 0 {
		entry.Ip = entry.Ipv4[0]
	}
	entry.Create = strings.TrimSpace(utils.ReadFile(createFile))
	entry.Pass = strings.TrimSpace(utils.ReadFile(passFile))
	entry.Domain = domainAddr
	entry.KeyPath = domainKey
	entry.CertPath = domainCert
	return entry
}

// 更新域名
func updateDomain() error {
	var list []string
	list = append(list, utils.NginxDefaultConf)
	for _, entry := range instancesList() {
		if entry.Ip != "" && entry.Domain != "" {
			list = append(list, utils.FromTemplateContent(utils.NginxDomainConf, map[string]any{
				"NAME":   entry.Name,
				"DOMAIN": entry.Domain,
				"IP":     entry.Ip,
			}))
		}
	}
	err := utils.WriteFile("/tmp/.codepass/nginx/default.conf", strings.Join(list, "\n"))
	if err != nil {
		return err
	}
	err = utils.WriteFile("/tmp/.codepass/docker/docker-compose.yml", utils.DockerComposeContent)
	if err != nil {
		return err
	}
	_, _ = utils.Cmd("-c", "docker-compose -f /tmp/.codepass/docker/docker-compose.yml down")
	_, err = utils.Cmd("-c", "docker-compose -f /tmp/.codepass/docker/docker-compose.yml up -d --remove-orphans")
	if err != nil {
		return err
	}
	return nil
}
