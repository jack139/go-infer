package types

// 处理函数入口类型
type funcType func(*map[string]interface{}) (*map[string]interface{}, error)

// 模型接口定义
type Model interface {
    ApiPath() (string)
    ApiEntry(*map[string]interface{}) (*map[string]interface{}, error)  // 处理 api 参数的过程

    Init() (error)  // 模型初始化，装入权重等
    Infer(*map[string]interface{}) (*map[string]interface{}, error)  // 模型推理的过程
}

var (
	// api 入口 与 处理过程 映射
	EntryMap = map[string]funcType{}

	// 模型列表
	ModelList []Model
)
