package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/poldi1405/bulkrename/plan"
	"gitlab.com/poldi1405/go-ansi"
)

// RemoveInvalidEntries checks every entry in files and removes it if there is an issue accessing it. Additionally an error message with additional information is shown.
func RemoveInvalidEntries(files []string) []string {
	for i, file := range files {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			fmt.Printf(ansi.Red("ERROR!")+" Unable to find %v\n", file)
		} else if os.IsPermission(err) {
			fmt.Printf("%v Unable to access %v\n", ansi.Red("ERROR!"), file)
		} else if os.IsTimeout(err) {
			fmt.Printf("%v Timeout while trying to access %v\n", ansi.Red("ERROR!"), file)
		} else if err != nil {
			fmt.Printf("%v An unknown error occured while trying to find %v. Error: %v\n", ansi.Red("ERROR!"), file, err)
		}
		if err != nil {
			// switch with last element and remove the last
			files[i] = files[len(files)-1]
			files = files[:len(files)-1]
		}
	}
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
