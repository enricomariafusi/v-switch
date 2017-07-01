package crypt

import (
	"V-switch/conf"
	"V-switch/tools"
	"crypto/rand"
	"io"
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

		originalText := make([]byte, 50*(i+1))
		if _, err := io.ReadFull(rand.Reader, originalText); err != nil {
			log.Println("[CRYPT][RAND][TEST] Problem %s", err.Error())
			return
		}

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
			log.Printf("[CRYPT][AES256][FAIL] Test #%d failed to encrypt %q to %q ", i, inverted, originalText)

		}

	}

	log.Printf("[CRYPT][AES256][TEST] Passed %d, Failed %d", passed, failed)

}

// Just used to force init to run
func GPGEngineStart() {

	log.Println("[CRYPT][AES] Triggering Crypt Engine to start")

}
