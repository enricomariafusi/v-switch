package crypt

import (
	"crypto/aes"
	"math/rand"
	"time"
)

const (
	// AES128 is AES-128 bit encryption
	AES128 = 16
	// AES192 is AES-192 bit encryption
	AES192 = 24
	// AES256 is AES-256 bit encryption
	AES256 = 32
)

// AESy struct for encryption
type AESy struct {
	AESKey     string
	Ciphertext []byte
	Plaintext  []byte
}

// NewAESy constructor for AESy
func NewAESy(key string) *AESy {
	aesy := new(AESy)
	aesy.AESKey = key

	return aesy
}

// Encrypt encrypts the string
func (aesy *AESy) Encrypt(plaintext string) error {
	aesy.Plaintext = []byte(plaintext)

	bc, err := aes.NewCipher([]byte(aesy.AESKey))

	if err != nil {
		return err
	}

	var src = []byte(plaintext)

	var ciphertext []byte

	// If the plaintext is greater than the blocksize
	if len(src) > bc.BlockSize() {
		srca := SplitBytesN(src, bc.BlockSize())

		for i := range srca {
			dst := make([]byte, bc.BlockSize())

			// No need to pad
			if len(srca[i]) == bc.BlockSize() {
				bc.Encrypt(dst, srca[i])

				ciphertext = append(ciphertext, dst...)
				// Need to pad
			} else {
				padSize := bc.BlockSize() - len(srca[i])
				padding := pad(padSize)

				newsrc := append(srca[i], padding...)

				bc.Encrypt(dst, newsrc)

				ciphertext = append(ciphertext, dst...)
			}
		}
		// If the plaintext is less than the blocksize
	} else {
		dst := make([]byte, bc.BlockSize())
		padSize := bc.BlockSize() - len(src)
		padding := pad(padSize)

		src = append(src, padding...)

		bc.Encrypt(dst, src)

		ciphertext = dst
	}

	aesy.Ciphertext = ciphertext

	return nil
}

// Decrypt decrypts the ciphertext
func (aesy *AESy) Decrypt(ciphertext string) error {
	aesy.Ciphertext = []byte(ciphertext)

	bc, err := aes.NewCipher([]byte(aesy.AESKey))

	if err != nil {
		return err
	}

	var src = []byte(ciphertext)

	var plaintext []byte

	srca := SplitBytesN(src, bc.BlockSize())

	for i := range srca {
		dst := make([]byte, bc.BlockSize())

		bc.Decrypt(dst, srca[i])

		plaintext = append(plaintext, dst...)
	}

	aesy.Plaintext = plaintext

	return nil
}

// KeyGen generates a key of given size (note that certain characters are omitted to make it easier to read)
func KeyGen(size int) []byte {
	rand.Seed(time.Now().UTC().UnixNano())

	const chars = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	result := make([]byte, size)

	for i := 0; i < size; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return result
}

// SplitBytesN chunks ciphertext into the proper blocksize for decryption
func SplitBytesN(src []byte, size int) [][]byte {
	var newb [][]byte

	for i := 0; i < len(src); i += size {
		newb = append(newb, src[i:i+size])
	}

	return newb
}

// pad generates padding to keep the cipher block the proper size
func pad(size int) []byte {
	result := make([]byte, size)

	for i := 0; i < size; i++ {
		result[i] = 0x00
	}

	return result
}
