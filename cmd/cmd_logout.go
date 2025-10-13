package cmd

import (
	"fmt"

	A "github.com/LFWQSP2641/scunet-auto-login/pkg/adapter"
	scu "github.com/LFWQSP2641/scunet-auto-login/pkg/schools/scu/auth"

	"github.com/spf13/cobra"
)

// logout 子命令
var logoutCommand = &cobra.Command{
	Use:   "logout",
	Short: "登出",
	Run: func(cmd *cobra.Command, args []string) {
		//// 校验
		//if username == "" || password == "" || service == "" {
		//	fmt.Println("用户名/密码/服务三项必须全部提供")
		//	return
		//}

		var auth A.Authenticator
		auth = scu.NewSCUAuthenticator()

		extra := map[string]string{
			"service": service,
		}

		if err := auth.Logout(globalCxt, username, password, extra); err != nil {
			fmt.Println("登出失败:", err)
			return
		}
		fmt.Println("登出流程完成")
	},
}
