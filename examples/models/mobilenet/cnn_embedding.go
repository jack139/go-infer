package mobilenet

import (
	"fmt"
	"log"
	"encoding/base64"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/jack139/go-infer/helper"
)

/* 训练好的模型权重 */
var (
	m *tf.SavedModel
)

/* 初始化模型 */
func initModel() error {
	var err error

	m, err = tf.LoadSavedModel(helper.Settings.Customer["MobilenetModelPath"], []string{"train"}, nil)
	if err != nil {
		return err
	}

	return nil
}


/*  定义模型相关参数和方法  */
type Mobilenet struct{}

func (x *Mobilenet) Init() error {
	return initModel()
}

func (x *Mobilenet) ApiPath() string {
	return "/api/mobile"
}

func (x *Mobilenet) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("ApiEntry_Mobilenet")

	// 检查参数
	imageBase64, ok := (*reqData)["image"].(string)
	if !ok {
		return &map[string]interface{}{"code":1001}, fmt.Errorf("need image")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"image": imageBase64,
	}

	return &reqDataMap, nil
}


// Mobilenet 推理
func (x *Mobilenet) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_Mobilenet")

	imageBase64 := (*reqData)["image"].(string)

	image, err  := base64.StdEncoding.DecodeString(imageBase64)
	if err!=nil {
		return &map[string]interface{}{"code":2001}, err
	}

	tensor, err := makeTensorFromBytes(image, 224, 224, 127.5, 127.5, false)
	if err!=nil {
		return &map[string]interface{}{"code":2002}, err
	}

	res, err := m.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.Graph.Operation("input_1").Output(0):      tensor,
		},
		[]tf.Output{
			m.Graph.Operation("out_relu/Relu6").Output(0),
		},
		nil,
	)
	if err != nil {
		return &map[string]interface{}{"code":2003}, err
	}

	ret := res[0].Value().([][][][]float32)
	return &map[string]interface{}{"embeddings":ret[0]}, nil
}
