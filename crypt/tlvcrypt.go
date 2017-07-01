package crypt

import (
	"V-switch/tools"
	"log"
)

func init() {

	log.Printf("[CRYPT][PGP] Engine INIT")

}

// FrameEncrypt returns an encrypted frame, given the frame and the key
func FrameEncrypt(key []byte, text2encrypt []byte) []byte {

	ace1 := NewAESy(string(key))

	err := ace1.Encrypt(string(text2encrypt))

	if err != nil {
		log.Printf("[AES256][ENCRYPT] Cannot encrypt frame: %s", err.Error())
		return nil
	}

	return ace1.Ciphertext

}

// FrameDecrypt returns the UNencrypted frame, given the encrypted frame and the key
func FrameDecrypt(key []byte, text2decrypt []byte) []byte {

	ace1 := NewAESy(string(key))

	err := ace1.Decrypt(string(text2decrypt))

	if err != nil {
		log.Printf("[AES256][DECRYPT] Cannot decrypt frame: %s", err.Error())
		return nil
	}

	return tools.CleanFrame(ace1.Plaintext)

}
