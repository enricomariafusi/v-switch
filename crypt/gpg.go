package crypt

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/openpgp"
)

func init() {

	log.Printf("[CRYPT][PGP] Engine INIT")

}

// FrameEncrypt returns an encrypted frame, given the frame and the key
func FrameEncrypt(key []byte, text2encrypt []byte) []byte {

	if len(text2encrypt) == 0 {
		log.Printf("[ENCRYPT][PGP] Can't encrypt NULL payload")
		return nil
	}

	var encrypted bytes.Buffer
	foo := bufio.NewWriter(&encrypted)

	plaintext, err := openpgp.SymmetricallyEncrypt(foo, key, nil, &GpGconfig)
	if err != nil {
		log.Printf("[ENCRYPT][PGP] Can't encrypt payload: %s", err.Error())
		return nil
	}

	// we encode it to avoid EOF problems with \0x00

	tmp2Enc := base64.StdEncoding.EncodeToString(text2encrypt)

	plaintext.Write([]byte(tmp2Enc))

	plaintext.Close()
	foo.Flush()
	return encrypted.Bytes()

}

// FrameDecrypt returns the UNencrypted frame, given the encrypted frame and the key
func FrameDecrypt(key []byte, text2decrypt []byte) []byte {

	encrypted := bytes.NewReader(text2decrypt)

	log.Printf("[DECRYPT][PGP] Start decrypting %d long payload", len(text2decrypt))

	cleartextMd, err := openpgp.ReadMessage(encrypted, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		return key, nil
	}, &GpGconfig)
	if err != nil {
		log.Printf("[DECRYPT][PGP] Can't decrypt payload: %s", err.Error())
		return nil
	}

	plaintext, err := ioutil.ReadAll(cleartextMd.UnverifiedBody)
	if err != nil {
		log.Printf("[DECRYPT][PGP] Can't read cleartext: %s", err.Error())
		return nil
	}

	log.Printf("[DECRYPT][PGP] Decryption successful")

	tmpClr, _ := base64.StdEncoding.DecodeString(string(plaintext[:]))

	return tmpClr
}

//GPGEngineStart triggers the init function in the package tap
func GPGEngineStart() {

	log.Println("[CRYPT][GPG] Engine Init")

}
