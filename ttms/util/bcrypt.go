package utils

import "golang.org/x/crypto/bcrypt"

func GetPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}
func ComparePwd(savehash, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(savehash), []byte(pwd)); err != nil {
		return false
	}
	return true
}
