package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type CryptoInterface interface {
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}

type crypto struct {
	key string
}

func New(key string) crypto {
	return crypto{key}
}

func GenerateKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}

func (c crypto) Encrypt(data string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(c.key)
	if err != nil {
		return "", errors.New("crypto key could not be decoded")
	}

	cb, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(data), nil)), nil
}

func (c crypto) Decrypt(data string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(c.key)
	if err != nil {
		return "", errors.New("crypto key could not be decoded")
	}

	cb, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cb)
	if err != nil {
		return "", err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(dataBytes) < nonceSize {
		return "", errors.New("data is too short")
	}

	nonce := dataBytes[:nonceSize]
	dataBytes = dataBytes[nonceSize:]

	decrypted, err := gcm.Open(nil, []byte(nonce), dataBytes, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
