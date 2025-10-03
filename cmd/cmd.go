package cmd

import (
	"context"
	C "scunet-auto-login/pkg/constant"

	"github.com/spf13/cobra"
)

var (
	globalCxt context.Context
)

var mainCommand = &cobra.Command{
	Use: "scunet-auto-login",
	Short: "SCU 网络自动登录 CLI\n" +
		"Version: " + C.Version,
}

func init() {
	runCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "设置配置文件路径 (默认: 当前工作目录 config.json)")
	runCmd.PersistentFlags().StringVarP(&workPath, "work", "w", "", "设置工作目录 (用于寻找默认 config.json)")
	loginCommand.PersistentFlags().StringVarP(&username, "username", "u", "", "用户名")
	loginCommand.PersistentFlags().StringVarP(&password, "password", "p", "", "密码")
	loginCommand.PersistentFlags().StringVarP(&service, "service", "s", "", "服务标识")

	mainCommand.AddCommand(loginCommand)
	mainCommand.AddCommand(runCmd)
}

// Execute 供 main 调用
func Execute(ctx context.Context) error {
	globalCxt = ctx
	return mainCommand.Execute()
}

// perRun 预留
func perRun(cmd *cobra.Command, args []string) {}
