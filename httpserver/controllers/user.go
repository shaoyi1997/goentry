package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/httpserver/rpc"
	"github.com/valyala/fasthttp"
)

const (
	tokenKey     = "sessionId"
	maxImageSize = 5000000
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

	responseError := response.GetError()
	if responseError != pb.LoginRegisterResponse_Ok {
		executeTemplate(ctx, controller.loginTemplate, nil)

		return
	}

	token := response.GetToken()
	if token == "" {
		ctx.Error(err.Error(), http.StatusInternalServerError)

		return
	}

	setSessionIDCookie(ctx, token)
	executeTemplate(ctx, controller.profileTemplate, response.User)
}

func setSessionIDCookie(ctx *fasthttp.RequestCtx, token string) {
	exp := time.Now().AddDate(0, 0, 1)
	cookie := fasthttp.Cookie{}
	cookie.SetKey(tokenKey)
	cookie.SetValue(token)
	cookie.SetExpire(exp)
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true)
	ctx.Response.Header.SetCookie(&cookie)
}

func (controller *UserController) LogoutHandler(ctx *fasthttp.RequestCtx) {
	username := string(ctx.FormValue("username"))
	token := extractToken(ctx)

	logoutRequest := &pb.LogoutRequest{
		Username: &username,
		Token:    &token,
	}

	response := new(pb.LogoutResponse)

	err := controller.client.CallMethod(pb.RpcRequest_Logout, logoutRequest, response)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)

		return
	}

	executeTemplate(ctx, controller.loginTemplate, nil)
}

func (controller *UserController) RegisterHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "RegisterHandler")
}

func (controller *UserController) UpdateUserHandler(ctx *fasthttp.RequestCtx) {
	updateRequest := controller.extractUpdateRequest(ctx)
	if updateRequest == nil {
		// ctx error is supplied in extractUpdateRequest
		return
	}

	response := new(pb.UpdateResponse)
	if err := controller.client.CallMethod(pb.RpcRequest_Update, updateRequest, response); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)

		return
	}

	responseErr := response.GetError()
	if responseErr == pb.UpdateResponse_InvalidToken {
		executeTemplate(ctx, controller.loginTemplate, nil)
	} else if responseErr != pb.UpdateResponse_Ok {
		executeTemplate(ctx, controller.loginTemplate, nil) // TODO: error in profile
	}

	executeTemplate(ctx, controller.profileTemplate, response.User)
}

func (controller *UserController) extractUpdateRequest(ctx *fasthttp.RequestCtx) *pb.UpdateRequest {
	form, err := ctx.MultipartForm()
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, fasthttp.ErrNoMultipartForm) {
			statusCode = http.StatusBadRequest
		}

		ctx.Error(err.Error(), statusCode)

		return nil
	}

	nickname := extractNickname(form)
	username := extractUsername(ctx)
	token := extractToken(ctx)
	encodedImageData, fileExtension := controller.extractImageDataAndExtension(ctx)

	updateRequest := &pb.UpdateRequest{
		Username:      &username,
		Token:         &token,
		Nickname:      &nickname,
		ImageData:     &encodedImageData,
		ImageFileType: &fileExtension,
	}

	return updateRequest
}

func extractNickname(form *multipart.Form) string {
	nicknames := form.Value["nickname"]

	var nickname string

	if len(nicknames) > 0 {
		nickname = nicknames[0]
	}

	return nickname
}

func (controller *UserController) extractImageDataAndExtension(ctx *fasthttp.RequestCtx) (string, string) {
	imageFileHeader, err := ctx.FormFile("profile_image")
	if err != nil {
		if errors.Is(err, fasthttp.ErrMissingFile) {
			return "", ""
		}

		logger.ErrorLogger.Println("Unexpected error:", err)

		ctx.Error("unexpected error", http.StatusInternalServerError)

		return "", ""
	}

	fileExtension := filepath.Ext(imageFileHeader.Filename)

	imageSize := imageFileHeader.Size
	if imageSize > maxImageSize {
		ctx.Error("image file is too huge", http.StatusBadRequest)

		return "", ""
	}

	imgData := make([]byte, imageSize)

	file, err := imageFileHeader.Open()
	if err != nil {
		logger.ErrorLogger.Println("Failed to open image file:", err)

		ctx.Error("failed to open image file", http.StatusInternalServerError)

		return "", ""
	}
	defer file.Close()

	n, err := file.Read(imgData)
	if err != nil || n <= 0 {
		ctx.Error("failed to read image file", http.StatusBadRequest)

		return "", ""
	}

	encodedImageData := base64.StdEncoding.EncodeToString(imgData)

	return encodedImageData, fileExtension
}

func (controller *UserController) UploadProfileImageHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "UploadProfileImageHandler")
}

func extractUsername(ctx *fasthttp.RequestCtx) string {
	return fmt.Sprintf("%v", ctx.UserValue("username"))
}

func extractToken(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(tokenKey))
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
