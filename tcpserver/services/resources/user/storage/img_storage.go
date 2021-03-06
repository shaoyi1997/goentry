package storage

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"git.garena.com/shaoyihong/go-entry-task/tcpserver/config"
)

type IImageStorage interface {
	StoreImage(username, fileName string, data *string) (string, error)
}

type ImageStorage struct {
	storageDir string
}

func NewImageStorage() IImageStorage {
	config := config.GetFileServerConfig()

	return &ImageStorage{storageDir: config.StorageDir}
}

func (storage *ImageStorage) StoreImage(username, fileName string, data *string) (string, error) {
	imgData, err := base64.StdEncoding.DecodeString(*data)
	if err != nil {
		logger.ErrorLogger.Println("Failed to base64 decode image file: ", err)

		return "", err
	}

	imgFolder := storage.storageDir + "/" + username
	if _, err := os.Stat(imgFolder); os.IsNotExist(err) {
		err = os.Mkdir(imgFolder, 0766)
		if err != nil {
			logger.ErrorLogger.Panicln("Failed to create image folder:", err)
		}
	}

	imgPath := imgFolder + "/" + fileName

	err = ioutil.WriteFile(imgPath, imgData, 0600)
	if err != nil {
		logger.ErrorLogger.Println("Failed to write image to file system: ", err)

		return "", err
	}

	return imgPath, nil
}
