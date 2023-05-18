package cmd

import (
	"codepass/app"
	"codepass/assets"
	"codepass/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/unrolled/secure"
	"html/template"
	"io"
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
			utils.PrintError("未安装 multipass，请使用 ./codepass install 命令安装或手动安装")
			os.Exit(1)
		}
		err = utils.WriteFile(utils.WorkDir("/service"), utils.FormatYmdHis(time.Now()))
		if err != nil {
			utils.PrintError("无法写入文件")
			os.Exit(1)
		}
		if app.ServiceConf.Conf == "" {
			yamls := [...]string{
				"/custom.yaml",
				"/custom.yml",
				"/config.yaml",
				"/config.yml",
			}
			for _, yaml := range yamls {
				if utils.IsFile(utils.RunDir(yaml)) {
					app.ServiceConf.Conf = utils.RunDir(yaml)
					break
				}
			}
		}
		if !utils.IsFile(app.ServiceConf.Conf) {
			utils.PrintError("配置文件不存在")
			os.Exit(1)
		}
		viper.SetConfigFile(app.ServiceConf.Conf)
		_ = viper.ReadInConfig()
		app.ServiceConf.Host = viper.GetString("host")
		app.ServiceConf.Port = viper.GetString("port")
		app.ServiceConf.SslKey = viper.GetString("ssl_key")
		app.ServiceConf.SslCrt = viper.GetString("ssl_crt")
		app.ServiceConf.GithubClientId = viper.GetString("github_client_id")
		app.ServiceConf.GithubClientSecret = viper.GetString("github_client_secret")
		if app.ServiceConf.Host == "" {
			app.ServiceConf.Host = "0.0.0.0"
		}
		if app.ServiceConf.Port == "" {
			app.ServiceConf.Host = "8443"
		}
		if !utils.IsFile(app.ServiceConf.SslKey) {
			utils.PrintError("SSL 私钥配置错误")
			os.Exit(1)
		}
		if !utils.IsFile(app.ServiceConf.SslCrt) {
			utils.PrintError("SSL 证书配置错误")
			os.Exit(1)
		}
		if app.ServiceConf.GithubClientId == "" || app.ServiceConf.GithubClientSecret == "" {
			utils.PrintError("GitHub 配置必须填写")
			os.Exit(1)
		}
		app.UpdateProxy()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if app.ServiceConf.Mode == "debug" {
			gin.SetMode(gin.DebugMode)
		} else if app.ServiceConf.Mode == "test" {
			gin.SetMode(gin.TestMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router := gin.Default()
		templates, err := loadWebTemplate()
		if err != nil {
			utils.PrintError(err.Error())
			os.Exit(1)
		}
		router.SetHTMLTemplate(templates)
		//
		router.Any("/*path", func(c *gin.Context) {
			urlHost := c.Request.Host
			regFormat := fmt.Sprintf("^((\\d+)-)*(([^\\/\\.]+)-([^\\/\\.]+)-([^\\/\\.]+)).%s", app.ServiceConf.Host)
			if utils.Test(urlHost, regFormat) {
				// 工作区实例
				reg := regexp.MustCompile(regFormat)
				match := reg.FindStringSubmatch(urlHost)
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
				// 接口、页面
				app.ServiceConf.OAuth(c)
			}
		})
		//
		router.Use(tlsHandler())
		err = router.RunTLS(fmt.Sprintf(":%s", app.ServiceConf.Port), app.ServiceConf.SslCrt, app.ServiceConf.SslKey)
		if err != nil {
			utils.PrintError(err.Error())
		}
	},
}

func loadWebTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range assets.Web.Files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(name, "/web/dist/assets/") {
			h, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			err = utils.WriteByte(utils.WorkDir("%s", name), h)
			if err != nil {
				return nil, err
			}
		}
		if strings.HasSuffix(name, ".html") {
			h, err := io.ReadAll(file)
			if err != nil {
				return nil, err
			}
			t, err = t.New(name).Parse(string(h))
			if err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}

func tlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     fmt.Sprintf(":%s", app.ServiceConf.Port),
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
	serviceCmd.Flags().StringVar(&app.ServiceConf.Conf, "conf", "", "配置文件路径")
	serviceCmd.Flags().StringVar(&app.ServiceConf.Mode, "mode", "release", "运行模式：debug/test/release")
}
