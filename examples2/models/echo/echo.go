package echo

import (
	//"fmt"
	"log"
)

/*  定义模型相关参数和方法  */
type EchoModel struct{}

func (x *EchoModel) Init() error {
	return nil
}

func (x *EchoModel) ApiPath() string {
	return "/api/echo"
}

func (x *EchoModel) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_EchoModel")

	return reqData, nil
	//return &map[string]interface{}{"code":9999}, fmt.Errorf("error test") // 错误返回： 错误代码，错误信息
}


// 这个不会被执行
func (x *EchoModel) Infer(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	return reqData, nil
}
