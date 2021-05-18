package user

import (
	"fmt"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/common/rpc"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"

	"git.garena.com/shaoyihong/go-entry-task/common/pb"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
)

type UserService struct {
	repo      IUserRepository
	hasher    IPasswordHasher
	validator IUserValidator
}

func NewUserService(database common.Database, redis *redis.Client) *UserService {
	return &UserService{
		repo:      NewUserRepository(database, redis),
		hasher:    newPasswordHasher(),
		validator: newUserValidator(),
	}
}
func (service *UserService) GetByUsername(username string) (*pb.User, error) {
	return service.repo.GetByUsername(username)
}

func (service *UserService) UpdateNickname(username, nickname string) error {
	return service.repo.UpdateNickname(username, nickname)
}

func (service *UserService) UpdateProfileImage(username, imageUrl string) error {
	return service.repo.UpdateProfileImage(username, imageUrl)
}

func (service *UserService) Register(username, password, nickname, imageUrl string) (*pb.User, error) {
	err := service.validator.ValidateLogin(username, password)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := service.hasher.hash(password)
	if err != nil {
		return nil, err
	}
	return service.repo.Insert(username, hashedPassword, nickname, imageUrl)
}

func (service *UserService) Login(messageByte []byte) ([]byte, error) {
	user, errorCode := service.processLogin(messageByte)

	response := &pb.LoginResponse{}

	if errorCode != nil {
		response = &pb.LoginResponse{
			Error: errorCode,
		}
	} else {
		response = &pb.LoginResponse{
			User:  user,
			Token: func(i string) *string { return &i }("token"),
		}
	}

	responseMessage, err := rpc.SerializeMessage(pb.RpcRequest_Login, response)
	if err != nil {
		logger.ErrorLogger.Println("Failed to serialize message:", err)
		return nil, err
	}
	return responseMessage, nil
}

func (service *UserService) processLogin(messageByte []byte) (*pb.User, *pb.LoginResponse_ErrorCode) {
	var args pb.LoginRequest
	var errorCode pb.LoginResponse_ErrorCode

	err := proto.Unmarshal(messageByte[:], &args)
	if err != nil {
		logger.ErrorLogger.Println("Failed to unmarshal message:", err)
		errorCode = pb.LoginResponse_InternalServerError
		return nil, &errorCode
	}

	username := *args.Username
	password := *args.Password

	_, err = service.Register(username, password, "", "")
	fmt.Println(err)

	err = service.validator.ValidateLogin(username, password)
	if err != nil {
		errorCode = pb.LoginResponse_MissingCredentials
		return nil, &errorCode
	}

	user, err := service.GetByUsername(username)
	if err != nil {
		if err == usernameNotFoundError {
			errorCode = pb.LoginResponse_InvalidUsername
		} else {
			errorCode = pb.LoginResponse_InternalServerError
		}
		return nil, &errorCode
	}

	isValidPassword := service.hasher.comparePasswords(*user.Password, password)
	if !isValidPassword {
		errorCode = pb.LoginResponse_InvalidPassword
		return nil, &errorCode
	}

	return user, nil
}
