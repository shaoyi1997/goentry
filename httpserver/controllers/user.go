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
	"git.garena.com/shaoyihong/go-entry-task/httpserver/view"
	"github.com/valyala/fasthttp"
)

const (
	tokenKey     = "sessionId"
	maxImageSize = 5000000
)

var tmpl = template.Must(template.ParseGlob("httpserver/view/*.html"))

type UserController struct {
	client rpc.IRPCClient
}

func NewUserController(rpcClient rpc.IRPCClient) UserController {
	return UserController{
		client: rpcClient,
	}
}

func (controller *UserController) LoginRegisterHandler(ctx *fasthttp.RequestCtx, isLogin bool) {
	username := string(ctx.FormValue("username"))
	password := string(ctx.FormValue("password"))

	var (
		request interface{}
		method  pb.RpcRequest_Method
	)

	if isLogin {
		request = &pb.LoginRequest{
			Username: &username,
			Password: &password,
		}
		method = pb.RpcRequest_Login
	} else {
		request = &pb.RegisterRequest{
			Username: &username,
			Password: &password,
		}
		method = pb.RpcRequest_Register
	}

	response := new(pb.LoginRegisterResponse)

	err := controller.client.CallMethod(method, request, response)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)

		return
	}

	responseError := response.GetError()
	if responseError != pb.LoginRegisterResponse_Ok {
		data := &view.LoginRegisterPageData{}

		switch responseError { //nolint:exhaustive
		case pb.LoginRegisterResponse_InvalidPassword:
			data.IsLoginPage = true
			data.InvalidPassword = true
		case pb.LoginRegisterResponse_InvalidUsername:
			data.IsLoginPage = true
			data.InvalidUsername = true
		case pb.LoginRegisterResponse_UsernameTaken:
			data.IsLoginPage = false
			data.UsernameTaken = true
		default:
			data.IsLoginPage = true
		}

		executeTemplate(ctx, view.Templates.LoginRegister, data)

		return
	}

	token := response.GetToken()
	if token == "" {
		ctx.Error(err.Error(), http.StatusInternalServerError)

		return
	}

	setSessionIDCookie(ctx, token)
	redirectToProfilePage(ctx)
}

func (controller *UserController) getProfile(ctx *fasthttp.RequestCtx, isEdit bool) {
	token := extractToken(ctx)
	if token == "" {
		redirectToLoginPage(ctx)

		return
	}

	getUserRequest := &pb.GetUserRequest{Token: &token}

	response := new(pb.GetUserResponse)

	err := controller.client.CallMethod(pb.RpcRequest_GetUser, getUserRequest, response)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)

		return
	}

	responseErr := response.GetError()
	if responseErr != pb.GetUserResponse_Ok {
		if responseErr == pb.GetUserResponse_UserNotFound {
			removeSessionIDCookie(ctx)
			redirectToLoginPage(ctx)

			return
		}
	}

	var template view.TemplateString

	if isEdit {
		template = view.Templates.Edit
	} else {
		template = view.Templates.Profile
	}

	user := response.GetUser()
	executeTemplate(ctx, template, user)
}

// GetProfilePage renders the profile page by the session token.
func (controller *UserController) GetProfilePage(ctx *fasthttp.RequestCtx) {
	controller.getProfile(ctx, false)
}

// GetEditPage renders the profile page by the session token, allowing the user to edit the page.
func (controller *UserController) GetEditPage(ctx *fasthttp.RequestCtx) {
	controller.getProfile(ctx, true)
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

	removeSessionIDCookie(ctx)
	redirectToLoginPage(ctx)
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
		removeSessionIDCookie(ctx)
		redirectToLoginPage(ctx)
	} else if responseErr != pb.UpdateResponse_Ok {
		data := &view.LoginRegisterPageData{
			IsLoginPage: true,
		}

		executeTemplate(ctx, view.Templates.LoginRegister, data)
	}

	redirectToProfilePage(ctx)
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

func (controller *UserController) GetLoginHandler(ctx *fasthttp.RequestCtx) {
	if isRedirected := controller.redirectIfCtxHasValidToken(ctx); isRedirected {
		return
	}

	data := &view.LoginRegisterPageData{
		IsLoginPage: true,
	}
	executeTemplate(ctx, view.Templates.LoginRegister, data)
}

func (controller *UserController) GetRegisterHandler(ctx *fasthttp.RequestCtx) {
	if isRedirected := controller.redirectIfCtxHasValidToken(ctx); isRedirected {
		return
	}

	data := &view.LoginRegisterPageData{
		IsLoginPage: false,
	}
	executeTemplate(ctx, view.Templates.LoginRegister, data)
}

func (controller *UserController) redirectIfCtxHasValidToken(ctx *fasthttp.RequestCtx) bool {
	if token := extractToken(ctx); token == "" {
		return false
	}

	redirectToProfilePage(ctx)

	return true
}

func executeTemplate(ctx *fasthttp.RequestCtx, template view.TemplateString, data interface{}) {
	err := tmpl.ExecuteTemplate(ctx, string(template), data)
	if err != nil {
		logger.ErrorLogger.Println("Failed to execute profile template:", err)
		ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	ctx.SetContentType("text/html; charset=utf-8")
}

func extractUsername(ctx *fasthttp.RequestCtx) string {
	return fmt.Sprintf("%v", ctx.UserValue("username"))
}

func extractToken(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(tokenKey))
}

func setSessionIDCookie(ctx *fasthttp.RequestCtx, token string) {
	exp := time.Now().AddDate(0, 0, 1)
	addCookie(ctx, exp, tokenKey, token)
}

func removeSessionIDCookie(ctx *fasthttp.RequestCtx) {
	exp := time.Now().AddDate(0, 0, -1)
	addCookie(ctx, exp, tokenKey, "")
}

func addCookie(ctx *fasthttp.RequestCtx, exp time.Time, key, value string) {
	cookie := fasthttp.Cookie{}
	cookie.SetKey(key)
	cookie.SetValue(value)
	cookie.SetExpire(exp)
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true)
	ctx.Response.Header.SetCookie(&cookie)
}

func redirectToLoginPage(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/login", http.StatusFound)
}

func redirectToProfilePage(ctx *fasthttp.RequestCtx) {
	ctx.Redirect("/profile", http.StatusFound)
}
