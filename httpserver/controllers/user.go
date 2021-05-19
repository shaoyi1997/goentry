package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"

	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"

	"github.com/valyala/fasthttp"
)

type UserController struct {
	client          rpc.IRPCClient
	profileTemplate *template.Template
	loginTemplate   *template.Template
}

func NewUserController(rpcClient rpc.IRPCClient) UserController {
	return UserController{
		client:          rpcClient,
		profileTemplate: template.Must(template.ParseFiles("./httpserver/view/profile.html")),
		loginTemplate:   template.Must(template.ParseFiles("./httpserver/view/login.html")),
	}
}

func (controller *UserController) LoginHandler(ctx *fasthttp.RequestCtx) {
	username := string(ctx.FormValue("username"))
	password := string(ctx.FormValue("password"))

	loginRequest := &pb.LoginRequest{
		Username: &username,
		Password: &password,
	}

	response := new(pb.LoginRegisterResponse)

	err := controller.client.CallMethod(pb.RpcRequest_Login, loginRequest, response)

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	token := response.Token
	if token != nil {
		exp := time.Now().AddDate(0, 0, 1)
		cookie := fasthttp.Cookie{}
		cookie.SetKey(username)
		cookie.SetValue(*token)
		cookie.SetExpire(exp)
		ctx.Response.Header.SetCookie(&cookie)
	}
	executeTemplate(ctx, controller.profileTemplate, response.User)
}

func (controller *UserController) LogoutHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "LogoutHandler")
}

func (controller *UserController) RegisterHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "RegisterHandler")
}

func (controller *UserController) UpdateUserHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "UpdateUserHandler")
}

func (controller *UserController) UploadProfileImageHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "UploadProfileImageHandler")
}

func extractUserId(ctx *fasthttp.RequestCtx) string {
	return fmt.Sprintf("%v", ctx.UserValue("user_id"))
}

func extractToken(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie("token"))
}

func (controller *UserController) GetLoginHandler(ctx *fasthttp.RequestCtx) {
	executeTemplate(ctx, controller.loginTemplate, nil)
}

func executeTemplate(ctx *fasthttp.RequestCtx, template *template.Template, data interface{}) {
	err := template.Execute(ctx, data)
	if err != nil {
		logger.ErrorLogger.Println("Failed to execute profile template:", err)
		ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	ctx.SetContentType("text/html; charset=utf-8")
}
