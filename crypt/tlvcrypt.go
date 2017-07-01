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
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, len(plaintext))
	eiv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, eiv); err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}

	stream := cipher.NewCFBEncrypter(block, eiv)
	stream.XORKeyStream(ciphertext, plaintext)

	// return converted frame
	return append(eiv, ciphertext...)
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
	div := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, div)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}
