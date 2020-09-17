package plan

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

var wg sync.WaitGroup

// LoadFileList loads the list of files into the Plan-Type
func (p *Plan) LoadFileList(files []string, recursive bool) {
	if recursive {
		L.Debug("entering recursive mode")
		for _, path := range files {
			L.Debug("working with file " + path)
			abspath, err := filepath.Abs(path)
			if err != nil {
				L.Error("Unable to get absolute Path of " + path)
				L.Info("Error: " + err.Error())
				continue
			}

			s, err := os.Stat(abspath)
			if err != nil {
				L.Error("Unable to access " + path)
				L.Info("Error: " + err.Error())
				continue
			}

			if s.IsDir() {
				L.Debug(path, "is a directory, scanning for files")
				wg.Add(1)
				go p.listAllFiles(abspath)
			} else {
				L.Debug(path, "is a file, appending to files")
				p.inFilesMtx.Lock()
				p.InFiles = append(p.InFiles, filepath.Clean(abspath))
				p.inFilesMtx.Unlock()
			}
		}
		L.Debug("Waiting for directory scans to finish")
		wg.Wait()
	} else {
		for _, path := range files {
			L.Debug("working with file " + path)
			abspath, err := filepath.Abs(path)
			if err != nil {
				L.Error("Unable to get absolute Path of " + path)
				L.Info("Error: " + err.Error())
				continue
			}

			s, err := os.Stat(abspath)
			if err != nil {
				L.Error("Unable to access " + path)
				L.Info("Error: " + err.Error())
				continue
			}

			if s.IsDir() {
				L.Debug("is a directory, appending path separator")
				abspath += string(os.PathSeparator)
			}
			p.inFilesMtx.Lock()
			p.InFiles = append(p.InFiles, filepath.Clean(abspath))
			p.inFilesMtx.Unlock()
		}
	}
	p.inFilesMtx.Lock()
	L.Debug("sorting filelist")
	sort.Strings(p.InFiles)
	p.inFilesMtx.Unlock()
}

func (p *Plan) listAllFiles(start string) error {
	var done bool
	defer wg.Done()

	go func() {
		select {
		case <-time.After(2 * time.Second):
			if !done {
				L.Debug("2 seconds elapsed, issuing warning")
				L.Warn(fmt.Sprintf("Scanning %v takes a long time. Please be patient.", start))
			}
		}
	}()
	var files []string

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		L.Debug("Found " + path)
		if info == nil {
			L.Trace("dafuq @ " + path)
			return nil
		}
		if !info.IsDir() {
			L.Debug(path, "is a file")
			files = append(files, filepath.Clean(path))
			return nil
		}

		L.Debug("is a directory")

		f, err := os.Open(path)
		if err != nil {
			L.Error("Error opening " + path)
			L.Info("Error: " + err.Error())
			return nil
		}
		defer f.Close()
		_, err = f.Readdirnames(1) // Or f.Readdir(1)
		if err == io.EOF {
			// Directory is empty, append it
			L.Debug("is empty")
			files = append(files, filepath.Clean(path)+string(os.PathSeparator))
			return nil
		}
		if err != nil {
			L.Error("Error while scanning " + path)
			L.Info("Error:" + err.Error())
			return nil
		}

		// directory is not empty, ignoring

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
		L.Error("Unable to create temporary file")
		L.Info("Error: " + err.Error())
		return err
	}
	for _, v := range p.GetFileList() {
		fmt.Fprintln(f, v)
		if err != nil {
			L.Error("Error writing filelist to temporary file")
			L.Trace("Path: " + v)
			L.Trace("TempFile: " + p.TempFile())
			L.Info("Error: " + err.Error())
			return err
		}
	}

	return nil
}
