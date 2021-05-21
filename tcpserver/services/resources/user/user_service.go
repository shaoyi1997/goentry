package user

import (
	"errors"
	"strings"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/common/rpc"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services/resources/user/storage"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	repo           IUserRepository
	hasher         IPasswordHasher
	validator      IUserValidator
	sessionManager ISessionManager
	imageStorage   storage.IImageStorage
}

func NewUserService(database common.Database, redis *redis.Client) *Service {
	service := &Service{
		repo:           NewUserRepository(database, redis),
		hasher:         NewPasswordHasher(),
		validator:      newUserValidator(),
		sessionManager: newSessionManager(redis),
		imageStorage:   storage.NewImageStorage(),
	}

	return service
}

// GetByUsername retrieves a user by the username.
func (service *Service) GetByUsername(username string) (*pb.User, error) {
	return service.repo.GetByUsername(username, true)
}

// GetByToken retrieves a user by the token.
func (service *Service) GetByToken(token string) (*pb.User, error) {
	username, err := service.sessionManager.GetCacheUsername(token)
	if err != nil {
		return nil, err
	}

	return service.GetByUsername(username)
}

// GetUser retrieves a user by either its username or session token.
func (service *Service) GetUser(messageByte []byte) ([]byte, error) {
	user, errorCode := service.processGetUser(messageByte)
	response := &pb.GetUserResponse{
		User:  user,
		Error: errorCode,
	}

	return serializeResponse(pb.RpcRequest_GetUser, response)
}

func (service *Service) processGetUser(messageByte []byte) (*pb.User, *pb.GetUserResponse_ErrorCode) {
	var (
		args      pb.GetUserRequest
		errorCode = pb.GetUserResponse_InternalServerError
	)

	err := proto.Unmarshal(messageByte, &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)

		return nil, &errorCode
	}

	if token := args.GetToken(); token != "" {
		user, _ := service.GetByToken(token)
		if user != nil {
			return user, nil
		}
	}

	if username := args.GetUsername(); username != "" {
		user, _ := service.GetByUsername(username)
		if user != nil {
			return user, nil
		}
	}

	errorCode = pb.GetUserResponse_UserNotFound

	return nil, &errorCode
}

// Update a user via the provided `UpdateRequest`.
func (service *Service) Update(messageByte []byte) ([]byte, error) {
	user, errorCode := service.processUpdate(messageByte)
	response := &pb.UpdateResponse{
		User:  user,
		Error: errorCode,
	}

	return serializeResponse(pb.RpcRequest_Update, response)
}

func (service *Service) processUpdate(messageByte []byte) (*pb.User, *pb.UpdateResponse_ErrorCode) {
	var (
		args      pb.UpdateRequest
		errorCode pb.UpdateResponse_ErrorCode
	)

	err := proto.Unmarshal(messageByte, &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)

		errorCode = pb.UpdateResponse_InternalServerError

		return nil, &errorCode
	}

	username := args.GetUsername()

	isValidToken := service.checkValidSessionToken(username, args.GetToken())
	if !isValidToken {
		errorCode = pb.UpdateResponse_InvalidToken

		return nil, &errorCode
	}

	err = service.repo.UpdateNickname(username, args.GetNickname())
	if err != nil {
		logger.ErrorLogger.Println("Failed to update nickname:", err)

		errorCode = pb.UpdateResponse_InternalServerError

		return nil, &errorCode
	}

	updateImageError := service.updateImage(&args)
	if updateImageError != nil {
		return nil, updateImageError
	}

	user, err := service.repo.GetByUsername(username, false)
	if err != nil {
		logger.ErrorLogger.Println("Failed to get user after update:", err)

		errorCode = pb.UpdateResponse_InternalServerError

		return nil, &errorCode
	}

	return user, nil
}

// updateImage performs the img storage & updates profile url in db. It is a no-op if there is no given image data.
func (service *Service) updateImage(args *pb.UpdateRequest) *pb.UpdateResponse_ErrorCode {
	var errorCode pb.UpdateResponse_ErrorCode

	username := args.GetUsername()
	imageData := args.GetImageData()
	imageFileExtension := args.GetImageFileType()

	if imageData == "" && imageFileExtension == "" {
		return nil
	}

	switch strings.ToLower(imageFileExtension) {
	case ".jpg", ".jpeg", ".png", ".bmp":
		break
	default:
		errorCode = pb.UpdateResponse_InvalidImageFile

		return &errorCode
	}

	storedImageName := username + "-profileImage" + imageFileExtension

	imageURL, err := service.imageStorage.StoreImage(username, storedImageName, &imageData)
	if err != nil {
		logger.ErrorLogger.Println("Failed to store image:", err)

		errorCode = pb.UpdateResponse_InternalServerError

		return &errorCode
	}

	err = service.repo.UpdateProfileImage(username, imageURL)
	if err != nil {
		logger.ErrorLogger.Println("Failed to update profile image:", err)

		errorCode = pb.UpdateResponse_InternalServerError

		return &errorCode
	}

	return nil
}

func (service *Service) Register(messageByte []byte) ([]byte, error) {
	user, token, errorCode := service.processRegister(messageByte)

	return generateLoginRegisterResponse(pb.RpcRequest_Register, user, token, errorCode)
}

func (service *Service) processRegister(messageByte []byte) (*pb.User, string, *pb.LoginRegisterResponse_ErrorCode) {
	var (
		args      pb.RegisterRequest
		errorCode pb.LoginRegisterResponse_ErrorCode
	)

	err := proto.Unmarshal(messageByte, &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return nil, "", &errorCode
	}

	return service.executeRegister(&args)
}

func (service *Service) executeRegister(args *pb.RegisterRequest) (*pb.User, string,
	*pb.LoginRegisterResponse_ErrorCode) {
	username := args.GetUsername()
	password := args.GetPassword()
	nickname := args.GetNickname()

	var errorCode pb.LoginRegisterResponse_ErrorCode

	err := service.validator.ValidateRegister(username, password)
	if err != nil {
		if errors.Is(err, errEmptyUsername) {
			errorCode = pb.LoginRegisterResponse_MissingCredentials
		} else if errors.Is(err, errTooShortPassword) {
			errorCode = pb.LoginRegisterResponse_InvalidPassword
		}

		logger.ErrorLogger.Println("Failed to validate register:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return nil, "", &errorCode
	}

	hashedPassword, err := service.hasher.Hash(password)
	if err != nil {
		logger.ErrorLogger.Println("Failed to hash password:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return nil, "", &errorCode
	}

	user, err := service.repo.Insert(username, hashedPassword, nickname, "")
	if err != nil {
		if errors.Is(err, errUsernameAlreadyExists) {
			errorCode = pb.LoginRegisterResponse_InvalidUsername

			return nil, "", &errorCode
		}

		logger.ErrorLogger.Println("Failed to insert user:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return nil, "", &errorCode
	}

	token, err := service.sessionManager.SetCacheToken(user.GetUsername())
	if err != nil {
		logger.ErrorLogger.Println("Failed to set token:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return user, "", &errorCode
	}

	return user, token, nil
}

func (service *Service) Logout(messageByte []byte) ([]byte, error) {
	logoutErr := service.processLogout(messageByte)
	isSuccess := logoutErr == nil
	response := &pb.LogoutResponse{
		Success: &isSuccess,
		Error:   logoutErr,
	}

	return serializeResponse(pb.RpcRequest_Logout, response)
}

func (service *Service) processLogout(messageByte []byte) *pb.LogoutResponse_ErrorCode {
	var (
		args      pb.LogoutRequest
		errorCode pb.LogoutResponse_ErrorCode
	)

	err := proto.Unmarshal(messageByte, &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)

		errorCode = pb.LogoutResponse_InternalServerError

		return &errorCode
	}

	username := args.GetUsername()

	err = service.validator.ValidateLogout(username)
	if err != nil {
		errorCode = pb.LogoutResponse_MissingUsername

		return &errorCode
	}

	go func() {
		service.sessionManager.DeleteCacheToken(username)
	}()

	return nil
}

func (service *Service) Login(messageByte []byte) ([]byte, error) {
	user, token, errorCode := service.processLogin(messageByte)

	return generateLoginRegisterResponse(pb.RpcRequest_Login, user, token, errorCode)
}

func (service *Service) processLogin(messageByte []byte) (*pb.User, string, *pb.LoginRegisterResponse_ErrorCode) {
	var (
		args      pb.LoginRequest
		errorCode pb.LoginRegisterResponse_ErrorCode
	)

	err := proto.Unmarshal(messageByte, &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return nil, "", &errorCode
	}

	username := args.GetUsername()
	password := args.GetPassword()

	err = service.validator.ValidateNonEmptyUsernamePassword(username, password)
	if err != nil {
		errorCode = pb.LoginRegisterResponse_MissingCredentials

		return nil, "", &errorCode
	}

	user, err := service.GetByUsername(username)
	if err != nil {
		if errors.Is(err, errUsernameNotFound) {
			errorCode = pb.LoginRegisterResponse_InvalidUsername
		} else {
			errorCode = pb.LoginRegisterResponse_InternalServerError
		}

		return nil, "", &errorCode
	}

	isValidPassword := service.hasher.ComparePasswords(*user.Password, password)
	if !isValidPassword {
		errorCode = pb.LoginRegisterResponse_InvalidPassword

		return nil, "", &errorCode
	}

	token, err := service.sessionManager.SetCacheToken(*user.Username)
	if err != nil {
		logger.ErrorLogger.Println("Failed to set token:", err)

		errorCode = pb.LoginRegisterResponse_InternalServerError

		return user, "", &errorCode
	}

	return user, token, nil
}

func generateLoginRegisterResponse(method pb.RpcRequest_Method, user *pb.User,
	token string, errorCode *pb.LoginRegisterResponse_ErrorCode) ([]byte, error) {
	var response *pb.LoginRegisterResponse

	if errorCode != nil {
		response = &pb.LoginRegisterResponse{
			User:  nil,
			Token: nil,
			Error: errorCode,
		}
	} else {
		response = &pb.LoginRegisterResponse{
			User:  user,
			Token: &token,
			Error: nil,
		}
	}

	return serializeResponse(method, response)
}

func (service *Service) checkValidSessionToken(username, token string) bool {
	storedToken, err := service.sessionManager.GetCacheToken(username)
	if err != nil {
		return false
	}

	return storedToken == token
}

func serializeResponse(method pb.RpcRequest_Method, response interface{}) ([]byte, error) {
	responseMessage, err := rpc.SerializeMessage(method, response)
	if err != nil {
		logger.ErrorLogger.Println("Failed to serialize message:", err)

		return nil, err
	}

	return responseMessage, nil
}
