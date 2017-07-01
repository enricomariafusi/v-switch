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

	eblock, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] problem %s", err.Error())
		return nil
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	eciphertext := make([]byte, len(plaintext))
	eiv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, eiv); err != nil {
		log.Println("[CRYPT][AES][IV] Problem %s", err.Error())
		return nil
	}

	stream := cipher.NewCFBEncrypter(eblock, eiv)
	stream.XORKeyStream(eciphertext, plaintext)

	// return converted frame
	return append(eiv, eciphertext...)
}

func FrameDecrypt(key []byte, cryptoText []byte) []byte {

	dblock, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(cryptoText) < aes.BlockSize {
		log.Println("[CRYPT][AES][ENC] Problem %s", err.Error())
		return nil
	}
	div := cryptoText[:aes.BlockSize]
	dciphertext := cryptoText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(dblock, div)

	dresult := make([]byte, len(dciphertext))

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(dresult, dciphertext)

	return dresult
}
