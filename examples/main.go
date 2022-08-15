package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"

	"examples/models/embedding"
	"examples/models/mobilenet"
	"examples/models/facedet"
)


var (
	rootCmd = &cobra.Command{
		Use:   "go-embedding",
		Short: "go-embedding examples",
	}
)

func init() {
	// 添加模型实例
	types.ModelList = append(types.ModelList, &embedding.BertEMB{})
	types.ModelList = append(types.ModelList, &mobilenet.Mobilenet{})
	types.ModelList = append(types.ModelList, &facedet.FaceDet{})

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
