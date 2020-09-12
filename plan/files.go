package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"gitlab.com/poldi1405/go-ansi"
)

var wg sync.WaitGroup

func (p *Plan) LoadFileList(files []string, recursive bool) {
	if recursive {
		for _, path := range files {
			abspath, err := filepath.Abs(path)
			if err != nil {
				fmt.Printf("%v Unable to get absolute path of %v. %v\n", ansi.Red("ERROR!"), path, err)
				continue
			}
			fmt.Println(abspath)

			s, err := os.Stat(abspath)
			if err != nil {
				fmt.Printf("%v Unable to stat %v. %v\n", ansi.Red("ERROR!"), path, err)
				continue
			}

			if s.IsDir() {
				fmt.Println("isdir")
				wg.Add(1)
				go p.listAllFiles(abspath)
			} else {
				p.inFilesMtx.Lock()
				p.InFiles = append(p.InFiles, filepath.Clean(abspath))
				p.inFilesMtx.Unlock()
			}
		}
		wg.Wait()
	} else {
		for _, path := range files {
			abspath, err := filepath.Abs(path)
			if err != nil {
				fmt.Printf("%v Unable to get absolute path of %v. %v\n", ansi.Red("ERROR!"), path, err)
				continue
			}

			s, err := os.Stat(abspath)
			if err != nil {
				fmt.Printf("%v Unable to stat %v. %v\n", ansi.Red("ERROR!"), path, err)
				continue
			}

			if s.IsDir() {
				abspath += string(os.PathSeparator)
			}
			p.inFilesMtx.Lock()
			p.InFiles = append(p.InFiles, filepath.Clean(abspath))
			p.inFilesMtx.Unlock()
		}
	}
	p.inFilesMtx.Lock()
	sort.Strings(p.InFiles)
	p.inFilesMtx.Unlock()
}

func (p *Plan) listAllFiles(start string) error {
	var done bool
	defer wg.Done()

	fmt.Println("here")
	go func() {
		select {
		case <-time.After(2 * time.Second):
			if !done {
				fmt.Printf("%v Scanning %v takes a long time. Please be patient.\n", ansi.Yellow("WARNING!"), start)
			}
		}
	}()
	var files []string

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		files = append(files, filepath.Clean(path))
		return nil
	})
	if err != nil {
		return err
	}
	done = true

	p.inFilesMtx.Lock()
	p.InFiles = append(p.InFiles, files...)
	p.inFilesMtx.Unlock()

	return nil
}

func (p *Plan) writeTempFile() error {
	f, err := os.Create(p.TempFile())
	defer f.Close()
	if err != nil {
		return err
	}
	for _, v := range p.GetFileList() {
		fmt.Fprintln(f, v)
		if err != nil {
			return err
		}
	}

	return nil
}
