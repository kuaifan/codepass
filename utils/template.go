package utils

import (
	"bytes"
	"codepass/assets"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// Assets 从模板中获取内容
func Assets(name string, envMap map[string]interface{}) string {
	content := ""
	for key, file := range assets.Shell.Files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(key, name) {
			h, err := io.ReadAll(file)
			if err == nil {
				content = strings.ReplaceAll(string(h), "\t", "    ")
				break
			}
		}
	}
	return Template(content, envMap)
}

// Template 从模板中获取内容
func Template(templateContent string, envMap map[string]interface{}) string {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			PrintError(fmt.Sprintf("模板分析失败: %s", err))
		}
	}()
	if err != nil {
		panic(1)
	}
	envMap["RUN_PATH"] = RunDir("")
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}
