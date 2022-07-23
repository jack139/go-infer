package http

import (
	"log"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"antigen-go/go-infer/helper"
	"antigen-go/go-infer/types"
)


/* 入口 */
func RunServer() {

	/* router */
	r := router.New()
	r.GET("/", index)
	for path := range types.EntryMap {
		r.POST(path, apiEntry)
		log.Println("router added: ", path)
	}

	host := fmt.Sprintf("%s:%d", helper.Settings.Api.Addr, helper.Settings.Api.Port)
	log.Printf("start HTTP server at %s\n", host)

	/* 启动server */
	s := &fasthttp.Server{
		Handler: helper.Combined(r.Handler),
		Name: "FastHttpLogger",
	}
	log.Fatal(s.ListenAndServe(host))
}

/* 根返回 */
func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}
