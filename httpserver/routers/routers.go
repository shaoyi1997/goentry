package routers

import (
	"github.com/fasthttp/router"
)

func InitRouter() *router.Router {
	router := router.New()
	initUserRouter(router)
	return router
}
