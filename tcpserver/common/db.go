package common

import (
	"database/sql"
	"fmt"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"

	// initialises mysql driver.
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Database struct {
	DB           *sql.DB
	DatabaseName string
}

func InitDB() *Database {
	databaseConfig := config.GetDatabaseConfig()
	var err error
	db, err = sql.Open(databaseConfig.Driver, databaseConfig.ConnectionURL)
	validateDBConnection(err)
	logger.InfoLogger.Println("Database connection initialised successfully")
	configDB(databaseConfig)
	createDB(databaseConfig)

	return &Database{
		DB:           db,
		DatabaseName: databaseConfig.Name,
	}
}

func validateDBConnection(err error) {
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		logger.ErrorLogger.Fatal(err)
	}
}

func configDB(databaseConfig config.DatabaseConfig) {
	db.SetMaxIdleConns(databaseConfig.MaxIdleConn)
	db.SetMaxOpenConns(databaseConfig.MaxOpenConn)
}

func createDB(databaseConfig config.DatabaseConfig) {
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_bin",
		databaseConfig.Name)
	if _, err := db.Exec(createDBQuery); err != nil {
		logger.ErrorLogger.Panicln("Failed to create database")
	}
}

func TearDownDB() {
	logger.InfoLogger.Println("Closing DB connection")

	if err := db.Close(); err != nil {
		logger.ErrorLogger.Println(err)
	}
}
