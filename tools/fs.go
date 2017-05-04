package tools

// Hpwd: the UNIX pwd
import (
	"os"
)

func Hpwd() string {

	tmpLoc, err := os.Getwd()

	if err != nil {
		tmpLoc = "/tmp"
	}

	return tmpLoc

}
