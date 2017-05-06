package main

import (
	co "V-switch/conf"
	aes "V-switch/crypt"
	pl "V-switch/plane"
	tap "V-switch/tap"
	to "V-switch/tools"
	"log"
	"os"
)

func init() {

	to.LogEngineStart()
	co.StartConfig()
	tap.EngineStart()
	pl.PlaneInit()
	aes.EngineStart()

}

func main() {

	log.Println("[OMG] End of bootstrap.")

	select {}

	os.Exit(0)

}
