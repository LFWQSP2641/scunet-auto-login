package cmd

import (
	"fmt"

	A "github.com/LFWQSP2641/scunet-auto-login/pkg/adapter"
	scu "github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/auth"

	"github.com/spf13/cobra"
)

// login 子命令
var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "登录",
	Run: func(cmd *cobra.Command, args []string) {
		// 校验
		if username == "" || password == "" || service == "" {
			fmt.Println("用户名/密码/服务三项必须全部提供")
			return
		}

		var auth A.Authenticator
		auth = scu.NewSCUAuthenticator()

		extra := map[string]string{
			"service": service,
		}

		if err := auth.Login(globalCxt, username, password, extra); err != nil {
			fmt.Println("登录失败:", err)
			return
		}
		fmt.Println("登录流程完成")
	},
}
