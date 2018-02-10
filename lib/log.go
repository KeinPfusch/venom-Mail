package lib

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type venomlogfile struct {
	filename string
	logfile  *os.File
	active   bool
}

var VSlogfile venomlogfile

func init() {

	// just the first time
	var currentFolder = Hpwd()
	os.MkdirAll(filepath.Join(currentFolder, "logs"), 0755)
	//

	VSlogfile.active = true
	VSlogfile.SetLogFolder()
	go VSlogfile.RotateLogFolder()

}

// rotates the log folder
func (lf *venomlogfile) RotateLogFolder() {

	for {

		time.Sleep(1 * time.Hour)
		if lf.logfile != nil {
			err := lf.logfile.Close()
			log.Println("[TOOLS][LOG] close logfile returned: ", err)
		}

		lf.SetLogFolder()

	}

}

// sets the log folder
func (lf *venomlogfile) SetLogFolder() {

	if lf.active {

		const layout = "2006-01-02.15"

		orario := time.Now().UTC()

		var currentFolder = Hpwd()
		lf.filename = filepath.Join(currentFolder, "logs", "venom."+orario.Format(layout)+"00.log")

		lf.logfile, _ = os.Create(lf.filename)

		log.Println("[TOOLS][LOG] Logfile is: " + lf.filename)
		log.SetOutput(lf.logfile)
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	} else {
		log.SetOutput(ioutil.Discard)
	}

}

// enables logging
func (lf *venomlogfile) EnableLog() {

	lf.active = true

}

func (lf *venomlogfile) DisableLog() {

	lf.active = false
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

}

//LogEngineStart just triggers the init for the package, and logs it.
func LogEngineStart() {

	log.Println("[TOOLS][LOG] LogRotation Init")

}
