package main

import (
	"fmt"
	"os"
	"strings"
)

// RemoveInvalidEntries checks every entry in files and removes it if there is an issue accessing it. Additionally an error message with additional information is shown.
func RemoveInvalidEntries(files []string) []string {
	dropped := 0
	for i, file := range files {
		l.Debug("trying file " + file)
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			l.Error(fmt.Sprintf("File %v does not exist", file))
			l.Info("Error: " + err.Error())
		} else if os.IsPermission(err) {
			l.Error(fmt.Sprintf("Access to %v denied", file))
			l.Info("Error: " + err.Error())
		} else if os.IsTimeout(err) {
			l.Error("Timeout while accessing " + file)
			l.Info("Error: " + err.Error())
		} else if err != nil {
			l.Error("Error while accessing File")
			l.Info("Error: " + err.Error())
		}
		if err != nil {
			l.Debug("an error occured, removing file from list")
			// switch with last element and remove the last
			files[i-dropped] = files[len(files)-1]
			files = files[:len(files)-1]
			dropped++
		}
	}
	l.Trace("Complete list of files:" + strings.Join(files, ":"))
	return files
}
