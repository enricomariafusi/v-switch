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

func FrameEncrypt(key []byte, text []byte) []byte {

	plaintext := text

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] problem %s", err.Error())
		return nil
	} else {
		log.Printf("[CRYPT][AES][ENC] Created NewCypher with blocksize %d", aes.BlockSize)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	if n, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	} else {
		log.Printf("[CRYPT][AES][ENC] Created IV[%d]", n)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	// return converted frame
	return append(iv, ciphertext...)
}

func FrameDecrypt(key []byte, cryptoText []byte) []byte {
	ciphertext := cryptoText

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}
