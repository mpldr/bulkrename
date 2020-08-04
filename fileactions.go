package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlab.com/poldi1405/bulkrename/plan"
)

func listFiles(p *plan.Plan, files []string, recursive bool) {
	if p.AbsolutePaths {
		for i := range files {
			filep, err := filepath.Abs(files[i])
			if err != nil {
				fmt.Printf("\033[33m\033[1mWARNING!\033[0m %v. Defaulting to relative path.\n", err)
				continue
			}
			files[i] = filep
		}
	}

	if recursive {
		for _, v := range files {
			fmt.Println(len(listAllFiles(v)))
		}
	}

	p.InFiles = files
	fmt.Println(p.InFiles)
}

func listAllFiles(start string) []string {
	var done bool

	go func() {
		select {
		case <-time.After(2 * time.Second):
			if !done {
				fmt.Printf("\033[33m\033[1mWARNING!\033[0m Scanning %v takes a long time. Please be patient.\n", start)
			}
		}
	}()
	var files []string

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	done = true
	return files
}
