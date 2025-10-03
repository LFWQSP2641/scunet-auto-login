package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"scunet-auto-login/pkg/adapter"
	scu "scunet-auto-login/pkg/schools/scu/auth"
	o "scunet-auto-login/pkg/schools/scu/option"

	"github.com/spf13/cobra"
)

var (
	configPath string
	workPath   string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "执行",
	Run: func(cmd *cobra.Command, args []string) {
		// 确定配置文件路径：显式 -c > 默认当前目录 config.json
		pathToUse := configPath
		if pathToUse == "" {
			pathToUse = filepath.Join(workOrCwd(workPath), "config.json")
		}
		fileCfg, err := loadConfig(pathToUse)
		if err != nil {
			fmt.Println("加载配置失败:", err)
			return
		}
		// 允许命令行覆盖配置文件中的对应字段
		user := fileCfg.Username
		pass := fileCfg.Password
		svc := fileCfg.Service

		// 最终校验
		if user == "" || pass == "" || svc == "" {
			fmt.Println("用户名/密码/服务三项必须全部提供 (命令行或配置文件)")
			return
		}

		var auth adapter.Authenticator
		auth = scu.NewSCUAuthenticator()

		var extra map[string]string
		extra = map[string]string{
			"service": svc,
		}

		if err := auth.Login(globalCxt, user, pass, extra); err != nil {
			fmt.Println("登录失败:", err)
			return
		}
		fmt.Println("登录流程完成")

		return
	},
}

func loadConfig(path string) (*o.UserData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c o.UserData
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func workOrCwd(w string) string {
	if w != "" {
		return w
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
