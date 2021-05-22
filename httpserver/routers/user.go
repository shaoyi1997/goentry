package routers

import (
	"git.garena.com/shaoyihong/go-entry-task/httpserver/controllers"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var userController controllers.UserController

func initUserRouter(router *router.Router, rpcClient rpc.IRPCClient) {
	userController = controllers.NewUserController(rpcClient)
	router.GET("/login", userController.GetLoginHandler)
	router.POST("/login", func(ctx *fasthttp.RequestCtx) {
		userController.LoginRegisterHandler(ctx, true)
	})
	router.POST("/logout", userController.LogoutHandler)
	router.GET("/register", userController.GetRegisterHandler)
	router.POST("/register", func(ctx *fasthttp.RequestCtx) {
		userController.LoginRegisterHandler(ctx, false)
	})
	router.GET("/profile", userController.GetProfilePage)
	router.GET("/edit", userController.GetEditPage)
	router.POST("/edit/{username}/", userController.UpdateUserHandler)
}
