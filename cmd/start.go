package cmd

import (
	"codepass/app"
	utils "codepass/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
)

var startCmd = &cobra.Command{
	Use:   "start",
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/create", func(c *gin.Context) {
			app.MConf.Create(c)
		})
		r.GET("/create/log", func(c *gin.Context) {
			app.MConf.CreateLog(c)
		})
		r.GET("/info", func(c *gin.Context) {
			app.MConf.Info(c)
		})
		r.GET("/delete", func(c *gin.Context) {
			app.MConf.Delete(c)
		})
		_ = r.Run(fmt.Sprintf("%s:%s", app.MConf.Ip, app.MConf.Port))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&app.MConf.Ip, "ip", "0.0.0.0", "启动服务的IP")
	startCmd.Flags().StringVar(&app.MConf.Port, "port", "8080", "启动服务的端口")
}