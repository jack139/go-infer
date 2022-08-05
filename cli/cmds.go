// Package to provide command-line settings
package cli

import (
	"fmt"
	"log"
	"strconv"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/http"
	"github.com/jack139/go-infer/server"
	"github.com/jack139/go-infer/helper"
)

var (
	// Command to start a HTTP server
	HttpCmd = &cobra.Command{
		Use:   "http",
		Short: "start http service",
		PersistentPreRunE: preRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 启动 http 服务
			http.RunServer()

			return nil
		},
	}

	// Command to start a Dispatcher server
	ServerCmd = &cobra.Command{
		Use:   "server <queue No.>",
		Short: "start dispatcher service",
		PersistentPreRunE: preRun,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("queue number needed")
			}

			_, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("queue number should be a integer")
			}

			// 启动 分发服务
			server.RunServer(args[0])

			return nil
		},
	}
)

func init(){
	HttpCmd.Flags().String("yaml", "config/settings.yaml", "yaml file path")
	ServerCmd.Flags().String("yaml", "config/settings.yaml", "yaml file path")
}

func preRun(cmd *cobra.Command, args []string) error {
	// 取得参数
	yaml, err := cmd.Flags().GetString("yaml")
	if err!=nil {
		log.Fatal(err)
	}
	// 载入配置文件，初始化
	helper.InitSettings(yaml)

	return nil
}
