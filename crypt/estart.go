package crypt

import (
	"V-switch/conf"
	"V-switch/tools"
	"log"
	"reflect"
)

func init() {

	log.Printf("[CRYPT][AES256] Engine INIT")

	// Now testing the engine.

	for i := 0; i < 16; i++ {

		l := 32

		originalText := []byte(tools.RandSeq(1000 + i - 8))
		key := []byte(conf.GetConfigItem("SWITCHID")) //at least as long as the MTU

		if len(key) != l {
			log.Printf("[CRYPT][AES256] Wrong key lenght (%d) ", len(key))
			log.Println("[CRYPT][AES256] AES256 must be 32Bytes. Generating a random one")
			log.Println("[CRYPT][AES256] PLEASE NOTICE THE SWITCH WILL BE ISOLATED")
			key = []byte(tools.RandSeq(l)) // key must be as long as the payload
			conf.SetConfigItem("SWITCHID", string(key[:]))
			log.Printf("[CRYPT][AES256] Your EXAMPLE safe key is: %s", key)
		}

		encrypted := FrameEncrypt(key, originalText)
		inverted := FrameDecrypt(key, encrypted)

		log.Println("[CRYPT][AES256] Originaltext Len: ", len(originalText))
		log.Println("[CRYPT][AES256] EncryptedText Len: ", len(encrypted))
		log.Println("[CRYPT][AES256] DEcryptedText Len: ", len(inverted))
		if reflect.DeepEqual(inverted, originalText) {
			log.Printf("[CRYPT][AES256] AES engine test #%d PASSED", i+1)
		} else {
			log.Printf("[CRYPT][AES256] AES engine test #%d FAILED", i+1)

		}

	}

}

// Just used to force init to run
func GPGEngineStart() {

	log.Println("[CRYPT][AES] Triggering Crypt Engine to start")

}
