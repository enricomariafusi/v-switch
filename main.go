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
	"os/signal"
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
	} else {
		log.Println("[MAIN] End of bootstrap.")
	}

	// Just a nice way to wait until the Operating system sends a kill signal.
	// select{} was just horrible.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	os.Exit(0)

}
