package tools

import (
	"os"
)

//Hpwd behaves like the unix pwd command, returning the current path
func Hpwd() string {

	tmpLoc, err := os.Getwd()

	if err != nil {
		tmpLoc = "/tmp"
	}

	return tmpLoc

}
