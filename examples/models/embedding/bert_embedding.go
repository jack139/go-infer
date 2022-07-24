package embedding

import (
	"fmt"
	"log"

	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/tokenize/vocab"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/jack139/go-infer/helper"
)

/* 训练好的模型权重 */
var (
	m *tf.SavedModel
	voc vocab.Dict
)

/* 初始化模型 */
func initModel() error {
	var err error
	voc, err = vocab.FromFile(helper.Settings.Customer["BertVocabPath"])
	if err != nil {
		return err
	}
	m, err = tf.LoadSavedModel(helper.Settings.Customer["BertModelPath"], []string{"train"}, nil)
	if err != nil {
		return err
	}

	return nil
}

/* 判断是否是英文字符 */
func isAlpha(c byte) bool {
	return (c>=65 && c<=90) || (c>=97 && c<=122)
}


/*  定义模型相关参数和方法  */
type BertEMB struct{}

func (x *BertEMB) Init() error {
	return initModel()
}

func (x *BertEMB) ApiPath() string {
	return "/api/embedding"
}

func (x *BertEMB) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("ApiEntry_BertEMB")

	// 检查参数
	text, ok := (*reqData)["text"].(string)
	if !ok {
		return &map[string]interface{}{"code":1001}, fmt.Errorf("need text")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"api": x.ApiPath(),
		"params": map[string]interface{}{
			"text": text,
		},
	}

	return &reqDataMap, nil
}


// Bert 推理
func (x *BertEMB) Infer(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_BertQA")

	const MaxSeqLength = 512

	text := (*reqData)["text"].(string)

	tkz := tokenize.NewTokenizer(voc)
	ff := tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: MaxSeqLength}
	// 获取 token 向量
	f := ff.Feature(text)

	tids, err := tf.NewTensor([][]int32{f.TokenIDs})
	if err != nil {
		return &map[string]interface{}{"code":2001}, err
	}
	mask, err := tf.NewTensor([][]int32{f.Mask})
	if err != nil {
		return &map[string]interface{}{"code":2002}, err
	}
	sids, err := tf.NewTensor([][]int32{f.TypeIDs})
	if err != nil {
		return &map[string]interface{}{"code":2003}, err
	}

	res, err := m.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.Graph.Operation("input_ids").Output(0):      tids,
			m.Graph.Operation("input_mask").Output(0):     mask,
			m.Graph.Operation("segment_ids").Output(0):    sids,
		},
		[]tf.Output{
			m.Graph.Operation("bert/pooler/Squeeze").Output(0),
		},
		nil,
	)
	if err != nil {
		return &map[string]interface{}{"code":2004}, err
	}

	ret := res[0].Value().([][]float32)
	return &map[string]interface{}{"embeddings":ret[0]}, nil
}
