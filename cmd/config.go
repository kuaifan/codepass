package cmd

import (
	utils "codepass/util"
	"fmt"
	"github.com/spf13/cobra"
)

type model struct {
	path         string
	host         string
	port         string
	sslCert      string
	sslKey       string
	clientId     string
	clientSecret string
}

var config model

var command = &cobra.Command{
	Use:   "config",
	Short: "生成配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		if config.path == "" {
			config.path = "./config.yaml"
		}
		// 将配置保存到yaml文件
		content := `host: ` + config.host + `
port: ` + config.port + `

ssl_key: ` + config.sslKey + `
ssl_crt: ` + config.sslCert + `

github_client_id: ` + config.clientId + `
github_client_secret: ` + config.clientSecret
		err := utils.WriteFile(config.path, content)
		if err != nil {
			utils.PrintError(fmt.Sprintf("无法写入文件: %s, error: %s", config.path, err.Error()))
			return
		}
		utils.PrintSuccess(fmt.Sprintf("配置文件已保存到: %s", config.path))
	},
}

func init() {
	rootCmd.AddCommand(command)
	command.Flags().StringVar(&config.path, "path", "", "保存配置文件路径")
	command.Flags().StringVar(&config.host, "host", "", "域名")
	command.Flags().StringVar(&config.port, "port", "", "端口")
	command.Flags().StringVar(&config.sslCert, "ssl-cert", "", "SSL 证书路径")
	command.Flags().StringVar(&config.sslKey, "ssl-key", "", "SSL 私钥路径")
	command.Flags().StringVar(&config.clientId, "client-id", "", "GitHub Client ID")
	command.Flags().StringVar(&config.clientSecret, "client-secret", "", "GitHub Client Secret")
}
