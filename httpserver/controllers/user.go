package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"

	"github.com/valyala/fasthttp"
)

type UserController struct {
	client rpc.IRPCClient
}

func NewUserController() UserController {
	client, err := rpc.NewRPCClient()
	if err != nil {
		logger.ErrorLogger.Panicln("Failed to create rpc client:", err)
	}
	return UserController{client: client}
}

func (controller *UserController) LoginHandler(ctx *fasthttp.RequestCtx) {
	username := string(ctx.FormValue("username"))
	password := string(ctx.FormValue("password"))

	loginRequest := &pb.LoginRequest{
		Username: &username,
		Password: &password,
	}

	response := new(pb.LoginResponse)

	err := controller.client.CallMethod(pb.RpcRequest_Login, loginRequest, response)

	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
		//} else if response.ErrorCode != uint32(protos.Constant_ERROR_CODE_OK) {
		//	http.Error(w, "something wrong happened", http.StatusOK)
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
	jsonResponse, err := json.Marshal(response)
	ctx.Response.SetBody(jsonResponse)
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
