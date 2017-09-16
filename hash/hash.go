package hash

import "golang.org/x/crypto/bcrypt"

func Make(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Check(plainText string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainText))
	return err == nil
}
