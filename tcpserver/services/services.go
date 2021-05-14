package services

import "git.garena.com/shaoyihong/go-entry-task/common/logger"

func Init() {
	logger.InitLogger()
	initDB()
}
func TearDown() {
	tearDownDB()
}
