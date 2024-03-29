package plan

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"git.sr.ht/~poldi1405/glog"
)

var wg sync.WaitGroup

// LoadFileList loads the list of files into the Plan-Type
func (p *Plan) LoadFileList(files []string, recursive bool) {
	glog.Debug("entering recursive mode")
	for _, path := range files {
		glog.Debug("working with file " + path)
		abspath, err := filepath.Abs(path)
		if err != nil {
			glog.Error("Unable to get absolute Path of " + path)
			glog.Info("Error: " + err.Error())
			continue
		}

		s, err := os.Stat(abspath)
		if err != nil {
			glog.Error("Unable to access " + path)
			glog.Info("Error: " + err.Error())
			continue
		}

		if s.IsDir() && recursive {
			glog.Debug(path, "is a directory, scanning for files")
			wg.Add(1)
			go p.listAllFiles(abspath)
		} else if s.IsDir() { // no recursion
			glog.Debug("is a directory, appending path separator")
			abspath += string(os.PathSeparator)
			p.inFilesMtx.Lock()
			p.InFiles = append(p.InFiles, filepath.Clean(abspath))
			p.inFilesMtx.Unlock()
		} else {
			glog.Debug(path, "is a file, appending to files")
			p.inFilesMtx.Lock()
			p.InFiles = append(p.InFiles, filepath.Clean(abspath))
			p.inFilesMtx.Unlock()
		}
	}
	glog.Debug("Waiting for directory scans to finish")
	wg.Wait()
	p.inFilesMtx.Lock()
	glog.Debug("sorting filelist")
	sort.Strings(p.InFiles)
	p.inFilesMtx.Unlock()
}

func (p *Plan) listAllFiles(start string) {
	var done bool
	defer wg.Done()

	go func() {
		<-time.After(2 * time.Second)
		if !done {
			glog.Debug("2 seconds elapsed, issuing warning")
			glog.Warn(fmt.Sprintf("Scanning %v takes a long time. Please be patient.", start))
		}
	}()
	var files []string

	err := filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		glog.Debug("Found " + path)
		if err != nil {
			glog.Debug("Error passed " + err.Error())
		}

		if info == nil {
			glog.Trace("dafuq @ " + path)
			return nil
		}
		if !info.IsDir() {
			glog.Debug(path, "is a file")
			files = append(files, filepath.Clean(path))
			return nil
		}

		glog.Debug("is a directory")

		f, err := os.Open(path)
		if err != nil {
			glog.Error("Error opening " + path)
			glog.Info("Error: " + err.Error())
			return nil
		}
		defer f.Close()
		_, err = f.Readdirnames(1) // Or f.Readdir(1)
		if err == io.EOF {
			// Directory is empty, append it
			glog.Debug("is empty")
			files = append(files, filepath.Clean(path)+string(os.PathSeparator))
			return nil
		}
		if err != nil {
			glog.Error("Error while scanning " + path)
			glog.Info("Error:" + err.Error())
			return nil
		}

		// directory is not empty, ignoring

		return nil
	})
	if err != nil {
		glog.Debug("error occurred: " + err.Error())
	}
	done = true

	p.inFilesMtx.Lock()
	p.InFiles = append(p.InFiles, files...)
	p.inFilesMtx.Unlock()
}

func (p *Plan) writeTempFile() error {
	f, err := os.Create(p.TempFile())
	if err != nil {
		glog.Error("Unable to create temporary file")
		glog.Info("Error: " + err.Error())
		return err
	}
	defer f.Close()
	for _, v := range p.GetFileList() {
		fmt.Fprintln(f, v)
		if err != nil {
			glog.Error("Error writing filelist to temporary file")
			glog.Trace("Path: " + v)
			glog.Trace("TempFile: " + p.TempFile())
			glog.Info("Error: " + err.Error())
			return err
		}
	}

	return nil
}
