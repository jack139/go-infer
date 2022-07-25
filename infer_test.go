package infer

import (
	"testing"
	"log"

	"github.com/jack139/go-infer/types"
	"github.com/jack139/go-infer/cli"
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

	log.Println("request data: ", *reqData)

	return reqData, nil
	//return &map[string]interface{}{"code":9999}, fmt.Errorf("parameters error test") // 错误返回： 错误代码，错误信息
}

func (x *EchoModel) Infer(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Model Infer()", x.ApiPath())

	log.Println("infer return data: ", reqData)

	return reqData, nil
	//return &map[string]interface{}{"code":9998}, fmt.Errorf("infer error test") // 错误返回： 错误代码，错误信息
}


func TestHttp(t *testing.T) {
	t.Log("test HTTP service")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 启动 http
	cli.HttpCmd.RunE(nil, nil)
}


func TestServer(t *testing.T) {
	t.Log("test Server")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 启动 分发服务
	cli.ServerCmd.RunE(nil, []string{"0"})
}
