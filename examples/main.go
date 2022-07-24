package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"

	"examples/models/echo"
	"examples/models/qa"
	"examples/models/embedding"
)


var (
	rootCmd = &cobra.Command{
		Use:   "go-infer",
		Short: "go-infer examples",
	}
)

func init() {
	// 添加模型实例
	types.ModelList = append(types.ModelList, &embedding.BertEMB{})

	// 添加 api 入口
	for m := range types.ModelList {
		types.EntryMap[types.ModelList[m].ApiPath()] = types.ModelList[m].ApiEntry
	}

	// 命令行设置
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cli.HttpCmd)
	rootCmd.AddCommand(cli.ServerCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
