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

func (p *Plan) LoadFileList(files []string, recursive bool) {
	if recursive {
		var wg sync.WaitGroup
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
				go func(pth string) {
					wg.Add(1)
					defer wg.Done()
					if err := p.listAllFiles(pth); err != nil {
						fmt.Printf("%v Error while scanning paths of %v. %v\n", ansi.Red("ERROR!"), path, err)
					}
				}(abspath)
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

// TODO: implement function
func (p *Plan) listAllFiles(start string) error {
	var done bool

	go func() {
		select {
		case <-time.After(2 * time.Second):
			if !done {
				fmt.Printf("%v Scanning %v takes a long time. Please be patient.\n", ansi.Yellow("WARNING!"), start)
			}
		}
	}()
	var files []string

	fileList := make([]string, 0)
	e := filepath.Walk(start, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	if e != nil {
		panic(e)
	}

	for _, file := range fileList {
		fmt.Println(file)
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
