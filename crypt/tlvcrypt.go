package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

func init() {

	log.Printf("[CRYPT][AES] Engine INIT")

}

func FrameEncrypt(key []byte, plaintext []byte) []byte {

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("[CRYPT][AES][NONCE] Problem %s", err.Error())
		return nil
	}

	return gcm.Seal(nonce, nonce, plaintext, nil)

}

func FrameDecrypt(key []byte, ciphertext []byte) []byte {

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[DECRYPT][AES] Problem %s", err.Error())
		return nil
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("[DECRYPT][AES] Problem %s", err.Error())
		return nil
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Println("[DECRYPT][AES] Problem %s", "Cyphertext too short")
		return nil
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	cleartext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return cleartext

}
