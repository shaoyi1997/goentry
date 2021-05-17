package routers

import (
	"git.garena.com/shaoyihong/go-entry-task/httpserver/controllers"
	"github.com/fasthttp/router"
)

func initUserRouter(router *router.Router) {
	router.POST("/user/login", controllers.LoginHandler)
	router.POST("/user/logout", controllers.LogoutHandler)
	router.POST("/user/register", controllers.RegisterHandler)
	router.PATCH("/user/{user_id}/", controllers.UpdateUserHandler)
	router.PATCH("/user/{user_id}/uploadProfileImage", controllers.UploadProfileImageHandler)
}
