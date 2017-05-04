package main

import (
	co "V-switch/conf"
	ta "V-switch/tap"
	to "V-switch/tools"
	"log"
	"os"
)

func init() {

	to.Log_Engine_Start()
	ta.TapEngineStart()
	co.StartConfig()

}

func main() {

	log.Println("[OMG] V-Switch starts now!")

	select {}

	os.Exit(0)

}
