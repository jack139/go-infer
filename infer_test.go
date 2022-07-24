package infer

import (
	"testing"
	"log"
	"fmt"

	"github.com/jack139/go-infer/types"
	"github.com/jack139/go-infer/http"
	"github.com/jack139/go-infer/server"
)

/*  定义模型相关参数和方法  */
type EchoModel struct{}

func (x *EchoModel) ApiPath() string {
	return "/api/echo"
}

func (x *EchoModel) Init() error {
	log.Println("Model Init()", x.ApiPath())
	return nil
}

func (x *EchoModel) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Model ApiEntry()", x.ApiPath())

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"api": x.ApiPath(),
		"params": map[string]interface{}{
			"data": *reqData,
		},
	}

	log.Println("request data: ", reqDataMap)

	return &reqDataMap, nil
	//return &map[string]interface{}{"code":9999}, fmt.Errorf("parameters error test") // 错误返回： 错误代码，错误信息
}

func (x *EchoModel) Infer(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Model Infer()", x.ApiPath())
	retData, ok := (*reqData)["data"].(map[string]interface{})
	if ok {
		return &retData, nil
	} else {
		return nil, fmt.Errorf("retrieve response data fail") // 错误返回： 错误代码，错误信息	
	}
}


func TestHttp(t *testing.T) {
	t.Log("test HTTP service")

	types.ModelList = append(types.ModelList, &EchoModel{})

	http.RunServer()
}


func TestServer(t *testing.T) {
	t.Log("test Server")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 启动 分发服务
	server.RunServer("0")
}
