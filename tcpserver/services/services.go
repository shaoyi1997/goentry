package services

import (
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/common"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/services/resources/user"
)

func Init() {
	db := common.InitDB()
	user.NewUserService(*db)
}

func TearDown() {
	common.TearDownDB()
}
