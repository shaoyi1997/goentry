package config

import (
	"os"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port       string `yaml:"port"`
	TCPAddress string `yaml:"tcp_address"`
}

type PoolConfig struct {
	InitialCapacity int `yaml:"initial_capacity"`
	MaxCapacity     int `yaml:"max_capacity"`
	WaitTimeout     int `yaml:"wait_timeout"`
	IdleTimeout     int `yaml:"idle_timeout"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Pool   PoolConfig   `yaml:"pool"`
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
	file, err := os.Open("httpserver/config/config.yml")
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	return file
}

func GetServerConfig() ServerConfig {
	return config.Server
}

func GetPoolConfig() PoolConfig {
	return config.Pool
}
