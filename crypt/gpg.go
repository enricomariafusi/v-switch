package crypt

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/openpgp"
)

func init() {

	log.Printf("[CRYPT][PGP] Engine INIT")

}

// FrameEncrypt returns an encrypted frame, given the frame and the key
func FrameEncrypt(key []byte, text2encrypt []byte) []byte {
	var encrypted bytes.Buffer
	foo := bufio.NewWriter(&encrypted)

	plaintext, _ := openpgp.SymmetricallyEncrypt(foo, key, nil, &GpGconfig)
	plaintext.Write(text2encrypt)

	plaintext.Close()
	foo.Flush()
	return encrypted.Bytes()

}

// FrameDecrypt returns the UNencrypted frame, given the encrypted frame and the key
func FrameDecrypt(key []byte, text2decrypt []byte) []byte {

	encrypted := bytes.NewReader(text2decrypt)

	log.Printf("[CRYPT][PGP] Start decrypting %d long payload", len(text2decrypt))

	cleartextMd, err := openpgp.ReadMessage(encrypted, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return key, nil
	}, &GpGconfig)
	if err != nil {
		log.Printf("[CRYPT][PGP] Can't decrypt payload: %s", err.Error())
		return nil
	}

	plaintext, err := ioutil.ReadAll(cleartextMd.UnverifiedBody)
	if err != nil {
		log.Printf("[CRYPT][PGP] Can't read cleartext: %s", err.Error())
		return nil
	}

	log.Printf("[CRYPT][PGP] Decryption successful")

	return plaintext
}

//GPGEngineStart triggers the init function in the package tap
func GPGEngineStart() {

	log.Println("[CRYPT][GPG] Engine Init")

}
