package routers

import (
	"git.garena.com/shaoyihong/go-entry-task/httpserver/controllers"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"
	"github.com/fasthttp/router"
)

var (
	userController controllers.UserController
)

func initUserRouter(router *router.Router, rpcClient rpc.IRPCClient) {
	userController = controllers.NewUserController(rpcClient)
	router.POST("/user/login", userController.LoginHandler)
	router.POST("/user/logout", userController.LogoutHandler)
	router.POST("/user/register", userController.RegisterHandler)
	router.PATCH("/user/{user_id}/", userController.UpdateUserHandler)
	router.PATCH("/user/{user_id}/uploadProfileImage", userController.UploadProfileImageHandler)
}
