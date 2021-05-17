package user

import (
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	hash(password string) (string, error)
	comparePasswords(hashedPassword string, plainPassword string) bool
}

func newPasswordHasher() IPasswordHasher {
	return &bCryptPasswordHasher{}
}

type bCryptPasswordHasher struct {
}

func (_ *bCryptPasswordHasher) hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.WarningLogger.Println("Failed to generate hash from password", err)
	}
	return string(hash), err
}

func (_ *bCryptPasswordHasher) comparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
