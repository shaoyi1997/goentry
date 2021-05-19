package user

import (
	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	Hash(password string) (string, error)
	ComparePasswords(hashedPassword string, plainPassword string) bool
}

func NewPasswordHasher() IPasswordHasher {
	return &bcryptPasswordHasher{}
}

type bcryptPasswordHasher struct {
}

func (_ *bcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.WarningLogger.Println("Failed to generate hash from password", err)
	}
	return string(hash), err
}

func (_ *bcryptPasswordHasher) ComparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
