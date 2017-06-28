package main

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/schollz/cryptopasta"
)

// ENCRYPTION TAKEN DIRECTLY FROM COWYO (SO IT IS INTERCHANGEABLE)

func EncryptString(toEncrypt string, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	encrypted, err := cryptopasta.Encrypt([]byte(toEncrypt), &key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}

func DecryptString(toDecrypt string, password string) (string, error) {
	key := sha256.Sum256([]byte(password))
	contentData, err := hex.DecodeString(toDecrypt)
	if err != nil {
		return "", err
	}
	bDecrypted, err := cryptopasta.Decrypt(contentData, &key)
	return string(bDecrypted), err
}
