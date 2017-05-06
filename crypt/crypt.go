package crypt

import (
	"V-switch/conf"
	"V-switch/tools"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"reflect"
)

func init() {

	originalText := []byte(tools.RandSeq(1500))
	key := []byte(conf.VConfig["SWITCHID"]) //32 BYTE fpr AES256

	if len(key) < 32 {
		log.Printf("[AES] Key too short (%d) ", len(key))
		log.Println("[AES] AES256 cannot be shorter than 32 bytes. Generating a random one")
		log.Println("[AES] PLEASE NOTICE THE SWITCH WILL BE ISOLATED")
		key = []byte(tools.RandSeq(32)) // 32 because of yes.
		conf.VConfig["SWITCHID"] = string(key[:])
		log.Printf("[AES] Your new safe key is: %s", key)
	}

	encrypted := FrameEncrypt(key, originalText)
	inverted := FrameDecrypt(key, encrypted)

	if reflect.DeepEqual(inverted, originalText) {
		log.Printf("[CRYPT] AES engine test PASSED")
	} else {
		log.Printf("[CRYPT] AES engine test FAILED")

	}

}

// encrypt string to base64 crypto using AES
func FrameEncrypt(key []byte, text []byte) []byte {

	plaintext := text

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT] AES problem %s", err.Error())
		return nil
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println("[CRYPT] AES problem %s", err.Error())
		return nil

	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// return converted frame
	return ciphertext
}

// decrypt from base64 to decrypted string
func FrameDecrypt(key []byte, cryptoText []byte) []byte {
	ciphertext := cryptoText

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("[CRYPT] AES problem %s", err.Error())
		return nil
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Println("[CRYPT] AES problem %s", err.Error())
		return nil
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}

//EngineStart triggers the init function in the package tap
func EngineStart() {

	log.Println("[CRYPT] AES Engine Starting...")

}
