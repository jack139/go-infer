package infer

import (
	"testing"
	"log"

	//"github.com/jack139/go-infer/helper"
	"github.com/jack139/go-infer/types"
	"github.com/jack139/go-infer/http"
	"github.com/jack139/go-infer/server"
)

/*  定义模型相关参数和方法  */
type EchoModel struct{}

func (x *EchoModel) Init() error {
	log.Println("Model Init()")
	return nil
}

func (x *EchoModel) ApiPath() string {
	return "/api/echo"
}

func (x *EchoModel) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Model ApiEntry()")
	return reqData, nil
	//return &map[string]interface{}{"code":9999}, fmt.Errorf("error test") // 错误返回： 错误代码，错误信息
}

func (x *EchoModel) Infer(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Model Infer()")
	return reqData, nil
}


func TestHttp(t *testing.T) {
	t.Log("test HTTP service")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 添加 api 入口
	for m := range types.ModelList {
		types.EntryMap[types.ModelList[m].ApiPath()] = types.ModelList[m].ApiEntry
	}

	http.RunServer()
}


func TestServer(t *testing.T) {
	t.Log("test Server")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 初始化模型
	for m := range types.ModelList {
		err := types.ModelList[m].Init()
		if err != nil {
			t.Log("Init deep model fail: ", types.ModelList[m].ApiPath())
			return
		}
	}

	// 启动 分发服务
	server.RunServer("0")
}
