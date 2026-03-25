package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("password can not be empty")
	}

	hash,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		return "",err
	}
	return string(hash),nil
}

func VerifyPassword(hashedPassword, plainPassword string) bool{
	err:= bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(plainPassword))
	return err == nil
}