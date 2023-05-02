package cmd

import (
	"codepass/app"
	utils "codepass/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/unrolled/secure"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
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
		if !utils.IsFile(app.ServiceConf.Key) {
			utils.PrintError("SSL私钥路径错误")
			os.Exit(1)
		}
		if !utils.IsFile(app.ServiceConf.Crt) {
			utils.PrintError("SSL证书路径错误")
			os.Exit(1)
		}
		app.UpdateProxy()
	},
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		//
		router.Any("/*path", func(c *gin.Context) {
			domain := c.Request.Host
			urlPath := c.Request.URL.Path
			regFormat := fmt.Sprintf("^((\\d+)-)*([a-zA-Z][a-zA-Z0-9_]*)-code.%s", app.ServiceConf.Host)
			if utils.Test(domain, regFormat) {
				reg := regexp.MustCompile(regFormat)
				match := reg.FindStringSubmatch(domain)
				port := match[2]
				name := match[3]
				lose := true
				for _, entry := range app.ProxyList {
					if entry.Name == name {
						c.Request.Header.Set("X-Real-Ip", c.ClientIP())
						c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
						var targetUrl *url.URL
						if port == "" {
							targetUrl, _ = url.Parse(fmt.Sprintf("http://%s:55123", entry.Ip))
						} else {
							targetUrl, _ = url.Parse(fmt.Sprintf("http://%s:%s", entry.Ip, port))
						}
						proxy := httputil.NewSingleHostReverseProxy(targetUrl)
						proxy.ServeHTTP(c.Writer, c.Request)
						lose = false
						break
					}
				}
				if lose {
					if port == "" {
						c.String(http.StatusNotFound, fmt.Sprintf("%s not found", name))
					} else {
						c.String(http.StatusNotFound, fmt.Sprintf("%s(%s) not found", name, port))
					}
				}
			} else {
				if strings.HasPrefix(urlPath, "/api/workspaces/create/log") {
					app.ServiceConf.WorkspacesCreateLog(c)
				} else if strings.HasPrefix(urlPath, "/api/workspaces/create") {
					app.ServiceConf.WorkspacesCreate(c)
				} else if strings.HasPrefix(urlPath, "/api/workspaces/list") {
					app.ServiceConf.WorkspacesList(c)
				} else if strings.HasPrefix(urlPath, "/api/workspaces/info") {
					app.ServiceConf.WorkspacesInfo(c)
				} else if strings.HasPrefix(urlPath, "/api/workspaces/delete") {
					app.ServiceConf.WorkspacesDelete(c)
				} else if strings.HasPrefix(urlPath, "/api/keys/info") {
					app.ServiceConf.KeysInfo(c)
				} else if strings.HasPrefix(urlPath, "/api/keys/save") {
					app.ServiceConf.KeysSave(c)
				} else if strings.HasPrefix(urlPath, "/assets") {
					c.File(fmt.Sprintf("./web/dist%s", urlPath))
				} else {
					c.File("./web/dist/index.html")
				}
			}
		})
		//
		router.Use(tlsHandler())
		err := router.RunTLS(fmt.Sprintf(":%s", app.ServiceConf.Port), app.ServiceConf.Crt, app.ServiceConf.Key)
		if err != nil {
			utils.PrintError(err.Error())
		}
	},
}

func tlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     fmt.Sprintf("%s:%s", app.ServiceConf.Host, app.ServiceConf.Port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		c.Next()
	}
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.Flags().StringVar(&app.ServiceConf.Host, "host", "0.0.0.0", "主机地址或IP")
	serviceCmd.Flags().StringVar(&app.ServiceConf.Port, "port", "443", "服务端口")
	serviceCmd.Flags().StringVar(&app.ServiceConf.Key, "key", "", "SSL私钥路径(KEY)")
	serviceCmd.Flags().StringVar(&app.ServiceConf.Crt, "crt", "", "SSL证书路径(PEM格式)")
}
