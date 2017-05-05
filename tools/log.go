package tools

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type vswitchlogfile struct {
	filename string
	logfile  *os.File
}

func init() {

	// just the first time
	var currentFolder = Hpwd()
	os.MkdirAll(filepath.Join(currentFolder, "logs"), 0755)
	//

	var mylogfile vswitchlogfile
	mylogfile.SetLogFolder()
	go mylogfile.RotateLogFolder()

}

// rotates the log folder
func (lf *vswitchlogfile) RotateLogFolder() {

	for {

		time.Sleep(1 * time.Hour)
		if lf.logfile != nil {
			err := lf.logfile.Close()
			log.Println("[LOG] close logfile returned: ", err)
		}

		lf.SetLogFolder()

	}

}

// sets the log folder
func (lf *vswitchlogfile) SetLogFolder() {

	const layout = "2006-01-02.15"

	orario := time.Now()

	var currentFolder = Hpwd()
	lf.filename = filepath.Join(currentFolder, "logs", "vswitch."+orario.Format(layout)+"00.log")
	log.Println("[LOG] Logfile is: " + lf.filename)

	lf.logfile, _ = os.Create(lf.filename)

	log.SetPrefix("V-SWITCH> ")
	log.SetOutput(lf.logfile)

}

//LogEngineStart just triggers the init for the package, and logs it.
func LogEngineStart() {

	log.Println("[LOG] LogRotation engine started")

}
