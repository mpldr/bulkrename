package main

import (
	"fmt"
	"os"

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
