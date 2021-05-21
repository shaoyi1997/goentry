package user

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
	"golang.org/x/crypto/bcrypt"
)

const (
	isUsingMd5                    = true
	lengthOfSeparatedPasswordHash = 2
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

	return md5Hash(password, salt)
}

func md5Hash(password string, salt []byte) (string, error) {
	passwordBytes := append([]byte(password), salt...)
	hash := md5.New()

	hash.Write(passwordBytes)

	encodedPassword := hex.EncodeToString(hash.Sum(nil))
	encodedSalt := hex.EncodeToString(salt)
	encodedPasswordWithSalt := fmt.Sprintf("%s:%s", encodedPassword, encodedSalt)

	return encodedPasswordWithSalt, nil
}

func generateRandomSalt() []byte {
	salt := make([]byte, 8)

	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

	return salt
}

func (*md5PasswordHasher) ComparePasswords(hashedPassword string, plainPassword string) bool {
	separatedPasswordSalt := strings.Split(hashedPassword, ":")
	if len(separatedPasswordSalt) != lengthOfSeparatedPasswordHash {
		return false
	}

	encodedSalt := separatedPasswordSalt[1]

	decodedSalt, err := hex.DecodeString(encodedSalt)
	if err != nil {
		return false
	}

	digest, _ := md5Hash(plainPassword, decodedSalt)

	return digest == hashedPassword
}
