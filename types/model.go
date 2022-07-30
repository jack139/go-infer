// Package to provide base Model interface definition
package types

// base Model interface definition
type Model interface {
	// 返回 api 的 http路径
	ApiPath() (string)
	// 处理 api 参数的过程
	ApiEntry(*map[string]interface{}) (*map[string]interface{}, error)

	// 模型初始化，装入权重等
	Init() (error)
	// 模型推理的过程
	Infer(*map[string]interface{}) (*map[string]interface{}, error)
}

var (
	// Models list which been used in API call and inference
	ModelList []Model
)
