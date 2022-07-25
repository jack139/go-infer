## 开发指南

以下简述使用go-infer快速构建模型推理api、和部署的过程。具体代码可以参考[examples](../examples)目录下完整示例。



### 1. 基本设计思路

请求影响流程如下图。框架隐藏了HTTP服务、队列处理、Dispatcher分发服务等通用逻辑，用户开发时，只需要处理API请求所带的参数（API入口函数）和模型推理部分（模型推理函数）。

<img src="arch2.png" alt="调用流程" width="300" />



代码实现中，框架定义了一个简化的model interface，用户只要实现其下的方法函数，就可以由框架实现上图中的流程。model interface定义如下：

```go
// 模型接口定义
type Model interface {
	ApiPath() (string) // HTTP URL 路径
	ApiEntry(*map[string]interface{}) (*map[string]interface{}, error)  // 处理API参数的过程
	Init() (error)  // 模型初始化，装入权重等
	Infer(*map[string]interface{}) (*map[string]interface{}, error)  // 模型推理的过程
}
```



- ApiPath() 比较简单，返回API的URL路径字符串，在HTTP server启动时，会根据这个串注册URL服务。
- ApiEntry() 用于对API传入的参数进行处理，通常根据业务逻辑对参数进行合法性检查。入参为一个key-value map，包含API传入的data字段的内容（API入参结构请参考[API文档模板](API.md)），返回值也是一个key-value map，包含传给推理函数的内容。
- Init() 用于载入模型权重和模型初始化相关的工作，在Dispatcher server启动时，会被调用。
- Infer() 用于实现具体的推理服务，入参是ApiEntry()处理过的参数数据，出参是将在API返回结果中data字段返回的内容。

具体可以参考[代码示例](../examples/models/embedding/bert_embedding.go)



### 4. 命令行集成

框架提供命令行集成，可以集成到用户的命令行指令中：

```go
// 添加模型实例
types.ModelList = append(types.ModelList, &embedding.BertEMB{})

// 命令行设置
rootCmd.AddCommand(cli.HttpCmd)
rootCmd.AddCommand(cli.ServerCmd)
```

在添加命令行之前，要将上述实现的model interface添加到框架的ModelList中，框架会进行初始化和注册等工作。

具体可以参考[代码示例](../examples/main.go)



### 5. 配置文件

具体可以参考[配置文件示例](../examples/config/settings.yaml)



### 6. 模型权重导出

#### (1) Tensorflow权重导出

可参考examples/export路径下convert_bert_to_pb.py



#### (2) Keras权重导出（todo）



#### (3) PyTorch权重导出（todo）



### 7. 系统部署

#### (1) 本地测试



#### (2) 部署架构的设计