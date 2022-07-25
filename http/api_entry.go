package http

import (
	"github.com/valyala/fasthttp"

	"github.com/jack139/go-infer/helper"
	"github.com/jack139/go-infer/types"
)


/* API 默认入口 */
func apiEntry(ctx *fasthttp.RequestCtx) {
	// POST 的数据
	content := ctx.PostBody()

	// 验签
	data, err := checkSign(content)
	if err != nil {
		code, _ := (*data)["code"].(int) // data() 有带回错误代码
		respError(ctx, code, err.Error())
		return
	}

	for mIndex := range types.ModelList {
		if types.ModelList[mIndex].ApiPath() == string(ctx.Path()) {
			// 处理api请求参数
			reqDataMap, err := types.ModelList[mIndex].ApiEntry(data)
			if err!=nil {
				if reqDataMap!=nil {
					if code, ok := (*reqDataMap)["code"].(int); ok { // ApiEntry() 有带回错误代码
						respError(ctx, code, err.Error()) 
						return
					}
				}
				respError(ctx, 9001, err.Error()) 
				return
			}

			requestId := generateRequestId()

			// 注册消息队列，在发redis消息前注册, 防止消息漏掉
			pubsub := helper.Redis_subscribe(requestId)
			defer pubsub.Close()

			// 发 请求消息
			err = helper.Redis_publish_request(requestId, reqDataMap)
			if err!=nil {
				respError(ctx, 9002, err.Error())
				return
			}

			// 收 结果消息
			respData, err := helper.Redis_sub_receive(pubsub)
			if err!=nil {
				respError(ctx, 9003, err.Error())
				return
			}

			// code==0 提交成功
			if (*respData)["code"].(float64)!=0 { 
				respError(ctx, int((*respData)["code"].(float64)), (*respData)["msg"].(string))
				return
			}

			delete(*respData, "code")

			respJson(ctx, respData) // 正常返回

			return
		}
	}

	respError(ctx, 9009, "unknow path") 
}
