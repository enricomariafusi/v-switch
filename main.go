package main

import (
	ta "V-switch/tap"
	to "V-switch/tools"
	"log"
	"os"
)

func init() {

	to.Log_Engine_Start()
	ta.TapEngineStart()

}

func main() {

	log.Println("[OMG] V-Switch starts now!")

	select {}

	os.Exit(0)

}
