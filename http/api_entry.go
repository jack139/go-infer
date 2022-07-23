package http

import (
	"github.com/valyala/fasthttp"

	"antigen-go/go-infer/helper"
	"antigen-go/go-infer/types"
)


/* 空接口, 只进行签名校验 */
func apiEntry(ctx *fasthttp.RequestCtx) {
	// POST 的数据
	content := ctx.PostBody()

	// 验签
	data, err := helper.CheckSign(content)
	if err != nil {
		helper.RespError(ctx, 9000, err.Error())
		return
	}

	for path := range types.EntryMap {
		if path == string(ctx.Path()) {
			ret, err := (types.EntryMap[path])(data)
			if err==nil {
				helper.RespJson(ctx, ret) // 正常返回
			} else {
				helper.RespError(ctx, (*ret)["code"].(int), err.Error()) 
			}
			return
		}
	}

	helper.RespError(ctx, 9900, "unknow path") 
}
