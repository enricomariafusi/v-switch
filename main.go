package main

import (
	co "V-switch/conf"
	ta "V-switch/tap"
	to "V-switch/tools"
	"log"
	"os"
)

func init() {

	to.LogEngineStart()
	ta.TapEngineStart()
	co.StartConfig()

}

func main() {

	log.Println("[OMG] End of bootstrap, V-Switch 100% operating")

	select {}

	os.Exit(0)

}
