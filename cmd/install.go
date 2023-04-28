package cmd

import (
	utils "codepass/util"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var cmdFile string

var updateCmd = &cobra.Command{
	Use:   "install",
	Short: "安装服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		cmdFile = utils.RunDir("/.codepass/install/cmd")
		//
		if !utils.CheckOs() {
			utils.PrintError("暂不支持的操作系统")
			os.Exit(1)
		}
		err := utils.WriteFile(cmdFile, utils.TemplateContent(utils.InstallExecContent, map[string]any{}))
		if err != nil {
			utils.PrintError(fmt.Sprintf("保存安装文件失败：%s", err.Error()))
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = utils.Cmd("-c", fmt.Sprintf("chmod +x %s", cmdFile))
		cmdString := exec.Command("/bin/sh", cmdFile)
		utils.PrintCmdOutput(cmdString)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
