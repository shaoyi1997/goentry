package controllers

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func LoginHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "LoginHandler")
}

func LogoutHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "LogoutHandler")
}

func RegisterHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "RegisterHandler")
}

func UpdateUserHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "UpdateUserHandler")
}

func UploadProfileImageHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "UploadProfileImageHandler")
}

func extractUserId(ctx *fasthttp.RequestCtx) string {
	return fmt.Sprintf("%v", ctx.UserValue("user_id"))
}
