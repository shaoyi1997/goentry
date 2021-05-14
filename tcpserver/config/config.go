package config

import (
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"gopkg.in/yaml.v2"
	"os"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Driver        string `yaml:"driver"`
	Name          string `yaml:"name"`
	ConnectionUrl string `yaml:"connection_url"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

var config Config

func init() {
	file := openConfigFile()
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
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
