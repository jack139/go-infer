package facedet

//  LD_LIBRARY_PATH=/usr/local/lib build/go-embedding server 0


import (
	"fmt"
	"log"
	"encoding/base64"
	"image"
	"image/color"
	"bytes"

	"github.com/disintegration/imaging"
	"github.com/ivansuteja96/go-onnxruntime"

	"github.com/jack139/go-infer/helper"
)

const (
	det_model_input_size = 224
)

/* 训练好的模型权重 */
var (
	//m *tf.SavedModel
	m *onnxruntime.ORTSession
)

/* 初始化模型 */
func initModel() error {
	var err error

	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_WARNING, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	m, err = onnxruntime.NewORTSession(ortEnvDet, helper.Settings.Customer["FacedetModelPath"], ortDetSO)
	if err != nil {
		return err
	}

	return nil
}


/*  定义模型相关参数和方法  */
type FaceDet struct{}

func (x *FaceDet) Init() error {
	return initModel()
}

func (x *FaceDet) ApiPath() string {
	return "/api/facedet"
}

func (x *FaceDet) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("ApiEntry_FaceDet")

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


// FaceDet 推理
func (x *FaceDet) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceDet")

	imageBase64 := (*reqData)["image"].(string)

	imgByte, err  := base64.StdEncoding.DecodeString(imageBase64)
	if err!=nil {
		return &map[string]interface{}{"code":2001}, err
	}


	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		return &map[string]interface{}{"code":2002}, err
	}


	shape1 := []int64{1, 3, det_model_input_size, det_model_input_size}
	input1 := preprocessImage(img, det_model_input_size)

	//fmt.Println(input1[:100])

	res, err := m.Predict([]onnxruntime.TensorValue{
		{
			Value: input1,
			Shape: shape1,
		},
	})
	if err != nil {
		return &map[string]interface{}{"code":2003}, err
	}

	if len(res) == 0 {
		return &map[string]interface{}{"code":2004}, fmt.Errorf("Failed get result")
	}

	for i:=0;i<len(res);i++ {
		log.Printf("Success do predict, shape : %+v, result : %+v\n", 
			res[i].Shape, 
			res[i].Value.([]float32)[:res[i].Shape[1]], // only show one value
		)
	}

	ret := res[2].Value.([]float32)
	return &map[string]interface{}{"result3":ret}, nil
}



func preprocessImage(src image.Image, inputSize int) []float32 {
	var newHeight, newWidth int
	im_ratio := float32(src.Bounds().Dx()) / float32(src.Bounds().Dy())
	if im_ratio > 1 { // width > height
		newWidth = inputSize
		newHeight = int(float32(newWidth) / im_ratio)
	} else {
		newHeight = inputSize
		newWidth = int(float32(newHeight) * im_ratio)		
	}

	fmt.Println(newWidth, newHeight)

	result := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)
	fmt.Println("resize: ", result.Rect)
	result = padBox(result)

	rgbs := make([]float32, inputSize*inputSize*3)

	j := 0
	for i := range result.Pix {
		if (i+1)%4 != 0 {
			rgbs[j] = float32(result.Pix[i])
			j++
		}
	}

	//fmt.Println(rgbs[:100])

	rgbs = Transpose(rgbs)

	//fmt.Println(rgbs[:100])

	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		rgbs[i] = normalize(rgbs[i], 127.5, 128.0)
		rgbs[i+channelLength] = normalize(rgbs[i+channelLength], 127.5, 128.0)
		rgbs[i+channelLength*2] = normalize(rgbs[i+channelLength*2], 127.5, 128.0)
	}
	return rgbs
}

func normalize(in float32, m float32, s float32) float32 {
	return (in - m) / s
}


// 调整为方形，黑色填充
func padBox(src image.Image) *image.NRGBA {
	var maxW int

	if src.Bounds().Dx() > src.Bounds().Dy() {
		maxW = src.Bounds().Dx()
	} else {
		maxW = src.Bounds().Dy()
	}

	dst := imaging.New(maxW, maxW, color.Black)
	dst = imaging.Paste(dst, src, image.Point{0,0})

	_ = imaging.Save(dst, "/tmp/test2.jpg")

	return dst
}


func Transpose(rgbs []float32) []float32 {
	out := make([]float32, len(rgbs))
	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		out[i] = rgbs[i*3]
		out[i+channelLength] = rgbs[i*3+1]
		out[i+channelLength*2] = rgbs[i*3+2]
	}
	return out
}

