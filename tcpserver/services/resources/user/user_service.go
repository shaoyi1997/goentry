package user

import (
	"strings"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/rpc"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services/resources/user/storage"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
)

type UserService struct {
	repo           IUserRepository
	hasher         IPasswordHasher
	validator      IUserValidator
	sessionManager ISessionManager
	imageStorage   storage.IImageStorage
}

func NewUserService(database common.Database, redis *redis.Client) *UserService {
	service := &UserService{
		repo:           NewUserRepository(database, redis),
		hasher:         NewPasswordHasher(),
		validator:      newUserValidator(),
		sessionManager: newSessionManager(redis),
		imageStorage:   storage.NewImageStorage(),
	}

	return service
}

func (service *UserService) GetByUsername(username string) (*pb.User, error) {
	return service.repo.GetByUsername(username, true)
}

func (service *UserService) Update(messageByte []byte) ([]byte, error) {
	user, errorCode := service.processUpdate(messageByte)
	response := &pb.UpdateResponse{
		User:  user,
		Error: errorCode,
	}

	return serializeResponse(pb.RpcRequest_Update, response)
}

func (service *UserService) processUpdate(messageByte []byte) (*pb.User, *pb.UpdateResponse_ErrorCode) {
	var args pb.UpdateRequest
	var errorCode pb.UpdateResponse_ErrorCode

	err := proto.Unmarshal(messageByte[:], &args)
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

// updateImage performs the image storage & updates profile url in storage. It is a no-op if there is no given image data
func (service *UserService) updateImage(args *pb.UpdateRequest) *pb.UpdateResponse_ErrorCode {
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

func (service *UserService) Register(messageByte []byte) ([]byte, error) {
	user, token, errorCode := service.processRegister(messageByte)
	return generateLoginRegisterResponse(pb.RpcRequest_Register, user, token, errorCode)
}

func (service *UserService) processRegister(messageByte []byte) (*pb.User, string, *pb.LoginRegisterResponse_ErrorCode) {
	var args pb.RegisterRequest
	var errorCode pb.LoginRegisterResponse_ErrorCode

	err := proto.Unmarshal(messageByte[:], &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)
		errorCode = pb.LoginRegisterResponse_InternalServerError
		return nil, "", &errorCode
	}

	return service.executeRegister(&args)
}

func (service *UserService) executeRegister(args *pb.RegisterRequest) (*pb.User, string, *pb.LoginRegisterResponse_ErrorCode) {
	username := args.GetUsername()
	password := args.GetPassword()
	nickname := args.GetNickname()

	var errorCode pb.LoginRegisterResponse_ErrorCode
	err := service.validator.ValidateRegister(username, password)
	if err != nil {
		if err == emptyUsernameError {
			errorCode = pb.LoginRegisterResponse_MissingCredentials
		} else if err == tooShortPasswordError {
			errorCode = pb.LoginRegisterResponse_InvalidPassword
		}
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
		if err == usernameAlreadyExistsError {
			errorCode = pb.LoginRegisterResponse_InvalidUsername
			return nil, "", &errorCode
		}
	}

	token, err := service.sessionManager.SetCacheToken(user.GetUsername())
	if err != nil {
		logger.ErrorLogger.Println("Failed to set token:", err)
		errorCode = pb.LoginRegisterResponse_InternalServerError
		return user, "", &errorCode
	}

	return user, token, nil
}

func (service *UserService) Logout(messageByte []byte) ([]byte, error) {
	logoutErr := service.processLogout(messageByte)
	isSuccess := logoutErr == nil
	response := &pb.LogoutResponse{
		Success: &isSuccess,
		Error:   logoutErr,
	}

	return serializeResponse(pb.RpcRequest_Logout, response)
}

func (service *UserService) processLogout(messageByte []byte) *pb.LogoutResponse_ErrorCode {
	var args pb.LogoutRequest
	var errorCode pb.LogoutResponse_ErrorCode

	err := proto.Unmarshal(messageByte[:], &args)
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

func (service *UserService) Login(messageByte []byte) ([]byte, error) {
	user, token, errorCode := service.processLogin(messageByte)
	return generateLoginRegisterResponse(pb.RpcRequest_Login, user, token, errorCode)
}

func (service *UserService) processLogin(messageByte []byte) (*pb.User, string, *pb.LoginRegisterResponse_ErrorCode) {
	var args pb.LoginRequest
	var errorCode pb.LoginRegisterResponse_ErrorCode

	err := proto.Unmarshal(messageByte[:], &args)
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
		if err == usernameNotFoundError {
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

func generateLoginRegisterResponse(method pb.RpcRequest_Method, user *pb.User, token string, errorCode *pb.LoginRegisterResponse_ErrorCode) ([]byte, error) {
	response := &pb.LoginRegisterResponse{}

	if errorCode != nil {
		response = &pb.LoginRegisterResponse{
			Error: errorCode,
		}
	} else {
		response = &pb.LoginRegisterResponse{
			User:  user,
			Token: &token,
		}
	}

	return serializeResponse(method, response)
}

func (service *UserService) checkValidSessionToken(username, token string) bool {
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
