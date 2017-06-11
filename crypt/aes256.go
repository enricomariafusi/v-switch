package crypt

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"reflect"
	"strconv"

	"crypto"

	"golang.org/x/crypto/openpgp/packet"
)

var GpGconfig packet.Config

func init() {

	log.Printf("[CRYPT][PGP] Engine INIT")

	GpGconfig.DefaultHash = crypto.SHA256

	log.Printf("[CRYPT][PGP] Configured ENGINE: hash is SHA256")

	GpGconfig.DefaultCipher = packet.CipherAES256
	log.Printf("[CRYPT][PGP] Configured ENGINE: Crypto is AES256")

	GpGconfig.DefaultCompressionAlgo = packet.CompressionZLIB
	log.Printf("[CRYPT][PGP] Configured ENGINE: Compression is GZIP")

	GpGconfig.CompressionConfig = &packet.CompressionConfig{
		Level: 9,
	}

	log.Printf("[CRYPT][PGP] Configured ENGINE: Gzip Compress level 9")

	// Now testing the engine.

	for i := 0; i < 10; i++ {

		l, _ := strconv.Atoi(conf.GetConfigItem("MTU"))

		originalText := []byte(tools.RandSeq(l + i - 5))
		key := []byte(conf.GetConfigItem("SWITCHID")) //at least as long as the MTU

		if len(key) < l {
			log.Printf("[CRYPT][GPG] Wrong key lenght (%d) ", len(key))
			log.Println("[CRYPT][GPG] AES256 cannot be shorter than MTU. Generating a random one")
			log.Println("[CRYPT][GPG] PLEASE NOTICE THE SWITCH WILL BE ISOLATED")
			key = []byte(tools.RandSeq(l)) // key must be as long as the payload
			conf.SetConfigItem("SWITCHID", string(key[:]))
			log.Printf("[CRYPT][GPG] Your EXAMPLE safe key is: %s", key)
		}

		encrypted := FrameEncrypt(key, originalText)
		inverted := FrameDecrypt(key, encrypted)

		log.Println("[CRYPT][GPG] Originaltext Len: ", len(originalText))
		log.Println("[CRYPT][GPG] EncryptedText Len: ", len(encrypted))

		if reflect.DeepEqual(inverted, originalText) {
			log.Printf("[CRYPT][GPG] AES engine test %d PASSED", i+1)
		} else {
			log.Printf("[CRYPT][GPG] AES engine test %d FAILED", i+1)

		}

	}

}
