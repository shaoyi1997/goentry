package config

import (
	"os"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Driver        string `yaml:"driver"`
	Name          string `yaml:"name"`
	ConnectionURL string `yaml:"connection_url"`
	MaxIdleConn   int    `yaml:"max_idle_conn"`
	MaxOpenConn   int    `yaml:"max_open_conn"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

var config Config

func InitConfig() {
	file := openConfigFile()
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logger.ErrorLogger.Fatal("Unable to decode config file")
	}
}

func openConfigFile() *os.File {
	file, err := os.Open("tcpserver/config/config.yml")
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	return file
}

func GetServerConfig() ServerConfig {
	return config.Server
}

func GetDatabaseConfig() DatabaseConfig {
	return config.Database
}
