package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/poldi1405/bulkrename/plan"
	"gitlab.com/poldi1405/go-ansi"
)

// RemoveInvalidEntries checks every entry in files and removes it if there is an issue accessing it. Additionally an error message with additional information is shown.
func RemoveInvalidEntries(files []string) []string {
	for i, file := range files {
		l.Debug("trying file " + file)
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			l.Error(fmt.Sprintf("File %v does not exist", file))
			l.Trace("Error: " + err.Error())
		} else if os.IsPermission(err) {
			l.Error(fmt.Sprintf("Access to %v denied", file))
			l.Trace("Error: " + err.Error())
		} else if os.IsTimeout(err) {
			l.Error("Timeout while accessing " + file)
			l.Trace("Error: " + err.Error())
		} else if err != nil {
			l.Error("Error while accessing File")
			l.Trace("Error: " + err.Error())
		}
		if err != nil {
			l.Debug("an error occured, removing file from list")
			// switch with last element and remove the last
			files[i] = files[len(files)-1]
			files = files[:len(files)-1]
		}
	}
	l.Trace("Complete list of files:" + strings.Join(files, ":"))
	return files
}

func listFiles(p *plan.Plan, files []string, recursive bool) {
	if p.AbsolutePaths {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println("\033[33m\033[1mWARNING!\033[0m Unable to determine current directory. Defaulting to relative paths")
		} else {
			for i := range files {
				files[i] = pwd + string(os.PathSeparator) + files[i]
			}
		}
	}

	if recursive {
		for _, v := range files {
			listAllFiles(v)
		}
	}

	p.InFiles = files
	fmt.Println(p.InFiles)
}

func listAllFiles(path string) []string {
	var result []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(ansi.Red("ERROR!"), err.Error())
		}
		fmt.Printf("File Name: %s\n", info.Name())
		return nil
	})
	return result
}
