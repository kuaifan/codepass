package cmd

import (
	"codepass/app"
	utils "codepass/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "启动服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		_, err := utils.Cmd("-c", "multipass version")
		if err != nil {
			utils.PrintError("未安装 multipass")
			os.Exit(1)
		}
		err = utils.WriteFile(utils.RunDir("/.codepass/service"), utils.FormatYmdHis(time.Now()))
		if err != nil {
			utils.PrintError("无法写入文件")
			os.Exit(1)
		}
		err = app.UpdateDomain()
		if err != nil {
			utils.PrintError(fmt.Sprintf("无法启动服务: %s", err.Error()))
			os.Exit(1)
		}
		go time.AfterFunc(1*time.Second, func() {
			_, url := app.InstanceDomain("")
			utils.PrintSuccess(fmt.Sprintf("\n服务地址: %s\n", url))
		})
	},
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/certs/info", func(c *gin.Context) {
			app.ServiceConf.CertsInfo(c)
		})
		r.POST("/certs/save", func(c *gin.Context) {
			app.ServiceConf.CertsSave(c)
		})
		r.GET("/workspaces/create", func(c *gin.Context) {
			app.ServiceConf.WorkspacesCreate(c)
		})
		r.GET("/workspaces/create/log", func(c *gin.Context) {
			app.ServiceConf.WorkspacesCreateLog(c)
		})
		r.GET("/workspaces/list", func(c *gin.Context) {
			app.ServiceConf.WorkspacesList(c)
		})
		r.GET("/workspaces/info", func(c *gin.Context) {
			app.ServiceConf.WorkspacesInfo(c)
		})
		r.GET("/workspaces/delete", func(c *gin.Context) {
			app.ServiceConf.WorkspacesDelete(c)
		})
		r.GET("/others/domain/update", func(c *gin.Context) {
			app.ServiceConf.OthersDomainUpdate(c)
		})
		_ = r.Run(fmt.Sprintf("%s:%s", app.ServiceConf.Ip, app.ServiceConf.Port))
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVar(&app.ServiceConf.Ip, "ip", "0.0.0.0", "启动服务的IP")
	serviceCmd.Flags().StringVar(&app.ServiceConf.Port, "port", "8080", "启动服务的端口")
}
