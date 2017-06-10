package tools

import (
	"log"
	"os"
)

//Hpwd behaves like the unix pwd command, returning the current path
func Hpwd() string {

	tmpLoc, err := os.Getwd()

	if err != nil {
		tmpLoc = "/tmp"
		log.Printf("[TOOLS][FS] Problem getting unix pwd: %s", err.Error())

	}

	return tmpLoc

}
