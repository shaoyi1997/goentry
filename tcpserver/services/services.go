package services

import (
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services/resources/user"
)

type Services struct {
	User *user.UserService
}

func Init() *Services {
	db := common.InitDB()
	redis := common.InitRedis()
	userService := user.NewUserService(*db, redis)
	return &Services{User: userService}
}

func TearDown() {
	common.TearDownDB()
	common.TearDownRedis()
}
