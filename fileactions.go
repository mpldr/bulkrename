package main

import (
	"os"
	"strings"

	"git.sr.ht/~poldi1405/glog"
)

// RemoveInvalidEntries checks every entry in files and removes it if there is an issue accessing it. Additionally an error message with additional information is shown.
func RemoveInvalidEntries(files []string) []string {
	dropped := 0
	for i, file := range files {
		glog.Debug("trying file " + file)
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			glog.Errorf("File %v does not exist", file)
			glog.Info("Error: " + err.Error())
		} else if os.IsPermission(err) {
			glog.Errorf("Access to %v denied", file)
			glog.Info("Error: " + err.Error())
		} else if err != nil {
			glog.Error("Error while accessing File")
			glog.Info("Error: " + err.Error())
		}
		if err != nil {
			glog.Debug("an error occurred, removing file from list")
			// switch with last element and remove the last
			files[i-dropped] = files[len(files)-1]
			files = files[:len(files)-1]
			dropped++
		}
	}
	glog.Tracef("Complete list of files: %s", strings.Join(files, ":"))
	return files
}
