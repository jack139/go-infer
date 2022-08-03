package infer

import (
	"testing"
	"log"
	"time"

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

func (x *EchoModel) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Model Infer()", x.ApiPath())

	log.Println("requestId", requestId, "infer return data: ", reqData)

	time.Sleep(1 * time.Second) // 延时，模拟推理业务

	return reqData, nil
	//return &map[string]interface{}{"code":9998}, fmt.Errorf("infer error test") // 错误返回： 错误代码，错误信息
}


func TestHttp(t *testing.T) {
	t.Log("test HTTP service")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 启动 http
	cli.HttpCmd.SetArgs([]string{"--yaml=config/settings.yaml"})
	cli.HttpCmd.Execute()

}


func TestServer(t *testing.T) {
	t.Log("test Server")

	types.ModelList = append(types.ModelList, &EchoModel{})

	// 启动 分发服务
	cli.ServerCmd.SetArgs([]string{"0", "--yaml=config/settings.yaml"})
	cli.ServerCmd.Execute()
}
