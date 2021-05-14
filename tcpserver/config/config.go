package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
}

var config Config

func init()  {
	file := openConfigFile()
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal("Unable to decode config file")
	}
}

func openConfigFile() *os.File {
	file, err := os.Open("tcpserver/config/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func GetServerConfig() ServerConfig {
	return config.Server
}
