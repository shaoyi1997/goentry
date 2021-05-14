package services

import (
	"database/sql"
	"fmt"
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	databaseConfig := config.GetDatabaseConfig()
	var err error
	db, err = sql.Open(databaseConfig.Driver, databaseConfig.ConnectionUrl)
	validateDBConnection(err)
	logger.InfoLogger.Println("Database connection initialised successfully")
	createDB(databaseConfig)
}

func validateDBConnection(err error) {
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		logger.ErrorLogger.Fatal(err)
	}
}

func createDB(databaseConfig config.DatabaseConfig) {
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_bin", databaseConfig.Name)
	if _, err := db.Exec(sql); err != nil {
		logger.ErrorLogger.Panicln("Failed to create database")
	}
}