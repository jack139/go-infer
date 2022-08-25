// Package to run a Http server to provide API service
// the http server is implemented by valyala/fasthttp
package http

import (
	"log"
	"fmt"
	"strings"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/jack139/go-infer/helper"
	"github.com/jack139/go-infer/types"
)


// start a HTTP service of API
func RunServer() {
	// 初始化SM2的密钥
	initSM2()

	/* router */
	r := router.New()
	r.GET("/", index)
	/* 添加模型定义的api入口 */
	for m := range types.ModelList {
		uri := types.ModelList[m].ApiPath()
		if strings.HasPrefix(uri, "__noapi__") {
			// 不注册 noapi 的入口
			log.Println("router skipped: ", uri)
		} else {
			// 注册 api 入口
			r.POST(uri, apiEntry)
			log.Println("router added: ", uri)
		}
	}

	host := fmt.Sprintf("%s:%d", helper.Settings.Api.Addr, helper.Settings.Api.Port)
	log.Printf("start HTTP server at %s\n", host)

	/* 启动server */
	s := &fasthttp.Server{
		Handler: combined(r.Handler),
		Name: "FastHttpLogger",
	}
	log.Fatal(s.ListenAndServe(host))
}

/* 根返回 */
func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}
