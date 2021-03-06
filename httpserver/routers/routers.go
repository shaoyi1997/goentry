package routers

import (
	"net/http"

	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func InitRouter(rpcClient rpc.IRPCClient) *router.Router {
	router := router.New()

	router.GET("/", rootHandler)
	registerCSSFileServer(router)
	registerJSFileServer(router)
	initUserRouter(router, rpcClient)

	return router
}

func registerCSSFileServer(router *router.Router) {
	router.ServeFiles("/css/{filepath:*}", "./httpserver/view/css")
}

func registerJSFileServer(router *router.Router) {
	router.ServeFiles("/js/{filepath:*}", "./httpserver/view/js")
}

func rootHandler(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/login", http.StatusFound)
}
