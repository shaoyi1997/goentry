package user

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"golang.org/x/crypto/bcrypt"
)

const (
	isUsingMd5 = true
)

type IPasswordHasher interface {
	Hash(password string) (string, error)
	ComparePasswords(hashedPassword string, plainPassword string) bool
}

func NewPasswordHasher() IPasswordHasher {
	if isUsingMd5 {
		return &md5PasswordHasher{}
	}

	return &bcryptPasswordHasher{}
}

type (
	bcryptPasswordHasher struct{}
	md5PasswordHasher    struct{}
)

func (*bcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		logger.WarningLogger.Println("Failed to generate hash from password", err)
	}

	return string(hash), err
}

func (*bcryptPasswordHasher) ComparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err == nil
}

func (*md5PasswordHasher) Hash(password string) (string, error) {
	salt := generateRandomSalt()
	passwordBytes := append([]byte(password), salt...)
	hash := md5.New()
	hash.Write(passwordBytes)
	hashedPassword := hash.Sum(nil)

	return base64.URLEncoding.EncodeToString(hashedPassword), nil
}

func generateRandomSalt() []byte {
	salt := make([]byte, 16)

	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	return salt
}

func (hasher *md5PasswordHasher) ComparePasswords(hashedPassword string, plainPassword string) bool {
	digest, _ := hasher.Hash(plainPassword)

	return digest == hashedPassword
}
