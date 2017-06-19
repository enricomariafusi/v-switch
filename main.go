package main

import (
	co "V-switch/conf"
	aes "V-switch/crypt"
	pl "V-switch/plane"
	tap "V-switch/tap"
	to "V-switch/tools"
	udp "V-switch/udp"
	"fmt"
	"log"
	"os"
)

func init() {

	if (os.Getuid() != 0) && (os.Getgid() != 0) {
		fmt.Println("[MAIN] YOU MUST BE ROOT TO RUN THIS PROGRAM , BECAUSE ONLY ROOT IS ALLOWED TO CREATE TAP DEVICES")
		os.Exit(1)
	}

	to.LogEngineStart()
	co.StartConfig()
	pl.PlaneInit()
	aes.GPGEngineStart()
	tap.EngineStart()
	udp.UDPEngineStart()

}

func main() {

	if co.GetConfigItem("DEBUG") != "TRUE" {
		log.Println("[MAIN] End of bootstrap.")
		to.VSlogfile.DisableLog()
	}

	select {}

	os.Exit(0)

}
