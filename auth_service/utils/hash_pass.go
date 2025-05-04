package utils

import (
	"golang.org/x/crypto/bcrypt"
)

//hash password
func HashPass(password string) (string, error)  {
	hashed, err:= bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashed), err
}

//if passhash == compatibles return true

func CheckPassHash(password, hash string) bool  {
	err:=bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err==nil
}