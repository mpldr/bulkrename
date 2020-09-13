// Package plan provides functions that are associated with the Plan and the
// type of the same name.
package plan

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "gitlab.com/poldi1405/bulkrename/plan/jobdescriptor"
	"gitlab.com/poldi1405/go-ansi"
)

// CreatePlan reads the new filenames from the temporary file
func (p *Plan) CreatePlan(planfile string) error {
	f, err := os.Open(planfile)
	if err != nil {
		L.Error("Unable to open temporary file")
		L.Trace("Error:", err)
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		scan := scanner.Text()
		L.Debug("Read", scan)
		var path string
		if scan != "" {
			path, err = filepath.Abs(scan)
			if err != nil {
				L.Error("Unable to get absolute path")
				L.Trace("Path:", scan)
				L.Trace("Error:", err)
				return err
			}
		}
		p.OutFiles = append(p.OutFiles, path)
	}

	for i, file := range p.InFiles {
		if i >= len(p.OutFiles) || p.OutFiles[i] == "" { // the line is empty
			if p.DeleteEmpty {
				p.jobs = append(p.jobs, JobDescriptor{Action: -1, SourcePath: file})
			}
		} else if file != p.OutFiles[i] { // the line is changed
			p.jobs = append(p.jobs, JobDescriptor{Action: 1, SourcePath: file, DstPath: p.OutFiles[i]})
		}
	}

	return nil
}

// PrepareExecution creates a set of prerules that need to be executed in order
// to execute the actual plan.
func (p *Plan) PrepareExecution() error {
	var prerules []JobDescriptor

	assumeExisting := make(map[string]bool)

	for _, job := range p.jobs {
		L.Debug("From:", job.SourcePath)
		L.Debug("To  :", job.DstPath)
		f, err := os.Open(job.SourcePath)
		if err != nil {
			f.Close()
			L.Error("Cannot access sourcefile")
			L.Trace("Path:", job.SourcePath)
			L.Trace("Error:", err)
			return err
		}

		fi, err := f.Stat()
		f.Close()
		if err != nil {
			L.Error("Cannot stat sourcefile")
			L.Trace("Path:", job.SourcePath)
			L.Trace("Error:", err)
			return err
		}

		if !fi.IsDir() {
			dir := filepath.Dir(job.DstPath)
			dir = strings.TrimSuffix(dir, fi.Name())
			if _, exists := assumeExisting[dir]; exists {
				continue
			}

			// if the containing folder doesn't exist, create it
			d, err := os.Open(dir)
			if os.IsNotExist(err) && p.CreateDirs {
				prerules = append(prerules, JobDescriptor{Action: 2, DstPath: dir + string(os.PathSeparator)})
				d.Close()
				continue
			} else if err != nil {
				d.Close()
				return err
			}

			dfi, err := d.Stat()
			d.Close()
			if err != nil {
				continue
			}

			// if it is not a directory but a file, delete (overwrite) it and remake it as a directory
			if !dfi.IsDir() {
				prerules = append(prerules, JobDescriptor{Action: -1, SourcePath: dir})
				prerules = append(prerules, JobDescriptor{Action: 2, SourcePath: dir})
			}
			assumeExisting[dir] = true
		} else {
			dst := job.DstPath
			//dst := strings.TrimSuffix(job.DstPath, string(os.PathSeparator))
			dst = strings.TrimSuffix(dst, filepath.Base(dst))
			if _, exists := assumeExisting[dst]; exists {
				continue
			}

			_, err := os.Stat(dst)
			if os.IsNotExist(err) && p.CreateDirs {
				prerules = append(prerules, JobDescriptor{Action: 2, DstPath: dst})
			} else if os.IsNotExist(err) {
				L.Error("Destination does not exist")
				L.Trace("Destination:", dst)
				L.Trace("Error:", err)
			} else if err != nil {
				L.Error("There is an issue with the destination directory")
				L.Trace("Destination:", dst)
				L.Trace("Error:", err)
				return err
			}
			L.Debug("assume that", dst, "does exist from now on")
			assumeExisting[dst] = true
		}
	}
	p.jobs = append(prerules, p.jobs...)
	return nil
}

// PreviewPlan prints a preview of the plan that is to be executed
func (p *Plan) PreviewPlan() {
	if len(p.jobs) == 0 {
		fmt.Println("There is nothing to do.")
		os.Exit(0)
	}
	for _, job := range p.jobs {
		switch job.Action {
		case -1:
			fmt.Printf(ansi.Yellow("delete:")+" %v\n", job.SourcePath)
		case 1:
			fmt.Printf(ansi.Yellow("move  :")+" %v "+ansi.Blue("⮕")+" %v\n", job.SourcePath, job.DstPath)
		case 2:
			fmt.Printf(ansi.Yellow("mkdir :")+" %v\n", job.DstPath)
		}
	}
}
