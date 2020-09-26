// Package plan provides functions that are associated with the Plan and the
// type of the same name.
package plan

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "gitlab.com/poldi1405/bulkrename/plan/jobdescriptor"
	"gitlab.com/poldi1405/go-ansi"
)

// CreatePlan reads the new filenames from the temporary file
func (p *Plan) CreatePlan(planfile string) error {
	f, err := os.Open(planfile)
	if err != nil {
		L.Error("Unable to open temporary file")
		L.Trace("Path:" + planfile)
		L.Info("Error:" + err.Error())
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		scan := scanner.Text()
		L.Debug("Read " + scan)
		var path string
		if scan != "" {
			path, err = filepath.Abs(scan)
			if err != nil {
				L.Error("Unable to get absolute path")
				L.Trace("Path:" + scan)
				L.Info("Error:" + err.Error())
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

	assumeExisting := make(map[string]bool)

	L.Debug("checking for circular file-movement")
	prerules := p.findCollisions()

	for _, job := range p.jobs {
		L.Debug("From:", job.SourcePath)
		L.Debug("To  :", job.DstPath)
		if job.Action == 3 { // this file was moved by the ringdetection
			L.Debug("ignoring this job, it was generated as collision prevention")
			continue
		}
		f, err := os.Open(job.SourcePath)
		if err != nil {
			f.Close()
			L.Error("Cannot access sourcefile")
			L.Trace("Path:" + job.SourcePath)
			L.Info("Error:" + err.Error())
			return err
		}

		fi, err := f.Stat()
		f.Close()
		if err != nil {
			L.Error("Cannot stat sourcefile")
			L.Trace("Path:" + job.SourcePath)
			L.Info("Error:" + err.Error())
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
			if !dfi.IsDir() && p.CreateDirs {
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
				L.Trace("Destination:" + dst)
				L.Info("Error:" + err.Error())
			} else if err != nil {
				L.Error("There is an issue with the destination directory")
				L.Trace("Destination:" + dst)
				L.Info("Error:" + err.Error())
				return err
			}
			L.Debug("assume that " + dst + " does exist from now on")
			assumeExisting[dst] = true
		}
	}
	p.jobs = append(prerules, p.jobs...)
	return nil
}

// findCollisions scans for file-switching. If there is a loop, break it.
func (p *Plan) findCollisions() []JobDescriptor {
	var prerules []JobDescriptor

	destinations := make(map[string]struct{})

	L.Debug("setting up map of destinationpaths")
	for _, j := range p.jobs {
		destinations[j.DstPath] = struct{}{}
	}

	for i := range p.jobs {
		L.Debug("From:", p.jobs[i].SourcePath)
		L.Debug("To  :", p.jobs[i].DstPath)
		_, match := destinations[p.jobs[i].SourcePath]
		if match { // this sourcefile is also a destination
			rand.Seed(time.Now().UnixNano())

			var safePath string
			for {
				safePath = p.jobs[i].SourcePath + "_" + strconv.Itoa(rand.Int())
				if _, err := os.Stat(safePath); os.IsNotExist(err) { // file does not exist, we may continue
					break
				}
			}
			L.Debug("Collision found, moving from " + p.jobs[i].SourcePath + " to " + safePath)

			moveToSafety := JobDescriptor{
				Action:     1,
				SourcePath: p.jobs[i].SourcePath,
				DstPath:    safePath,
			}

			prerules = append(prerules, moveToSafety)
			p.jobs[i].SourcePath = safePath
			p.jobs[i].Action = 3
		}
	}

	return prerules
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
		case 3: // rescued from being wrongfully overwritten
			fmt.Printf(ansi.Yellow("rcvr  :")+" %v "+ansi.Blue("⮕")+" %v\n", job.SourcePath, job.DstPath)
		}
	}
}
