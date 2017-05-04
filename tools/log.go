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
func (this *vswitchlogfile) RotateLogFolder() {

	for {

		time.Sleep(1 * time.Hour)
		if this.logfile != nil {
			err := this.logfile.Close()
			log.Println("[LOG] close logfile returned: ", err)
		}

		this.SetLogFolder()

	}

}

// sets the log folder
func (this *vswitchlogfile) SetLogFolder() {

	const layout = "2006-01-02.15"

	orario := time.Now()

	var currentFolder = Hpwd()
	this.filename = filepath.Join(currentFolder, "logs", "vswitch."+orario.Format(layout)+"00.log")
	log.Println("[LOG] Logfile is: " + this.filename)

	this.logfile, _ = os.Create(this.filename)

	log.SetPrefix("V-SWITCH> ")
	log.SetOutput(this.logfile)

}

func Log_Engine_Start() {

	log.Println("[LOG] LogRotation engine started")

}
