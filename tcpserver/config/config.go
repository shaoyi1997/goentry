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

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	PoolSize int    `yaml:"pool_size"`
}

type FileServerConfig struct {
	Address    string `yaml:"address"`
	StorageDir string `yaml:"storage_dir"`
}

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	FileServer FileServerConfig `yaml:"file_server"`
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

func GetRedisConfig() RedisConfig {
	return config.Redis
}

func GetFileServerConfig() FileServerConfig {
	return config.FileServer
}
