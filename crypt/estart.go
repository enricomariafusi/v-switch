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

	passed := 0
	failed := 0

	for i := 0; i < 100; i++ {

		l := 32

		originalText := []byte(tools.RandSeq(50 * i))
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

		if reflect.DeepEqual(inverted, originalText) {
			passed++

		} else {
			failed++

		}

	}

	log.Printf("[CRYPT][AES256][TEST] Passed %d, Failed %d", passed, failed)

}

// Just used to force init to run
func GPGEngineStart() {

	log.Println("[CRYPT][AES] Triggering Crypt Engine to start")

}
