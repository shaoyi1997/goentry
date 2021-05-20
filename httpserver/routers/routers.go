package routers

import (
	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"
	"github.com/fasthttp/router"
)

func InitRouter(rpcClient rpc.IRPCClient) *router.Router {
	router := router.New()
	initUserRouter(router, rpcClient)

	return router
}
