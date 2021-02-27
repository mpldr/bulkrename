package plan

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	j "git.sr.ht/~poldi1405/bulkrename/plan/jobdescriptor"
	"github.com/mborders/logmatic"
)

// Plan stores all information on a rename-job. and provides related funtions.
type Plan struct {
	// TempID contains the id of the Plan (i.e. the unique identifier of the job.
	TempID string
	// InFiles contains the paths of the files before renaming
	InFiles []string
	// OutFiles contains the paths of the files after renaming
	OutFiles []string

	// AbsolutePaths indicates whether to use absolute paths or not
	AbsolutePaths bool
	// Overwrite indicates whether existing files shall be overwritten
	Overwrite bool
	// Editor contains the Editor to use for editing
	Editor string
	// EditorArgs contains the arguments that will be passed to the editor. {} will be replaced by the absolute path to the plan-file.
	EditorArgs []string
	// CreateDirs indicates whether non-existent directories should be created as needed
	CreateDirs bool
	// StopToShow indicates whether an overview of the applied actions should be shown and confirmation requested
	StopToShow bool
	// DeleteEmpty indicates whether files corresponding with empty lines should be deleted
	DeleteEmpty bool

	// inFilesMtx is a Mutex to ensure that there are no issues when appending to the filelist
	inFilesMtx sync.Mutex
	// jobs contains the tasks that have to be executed in order for the target state to be reached
	jobs []j.JobDescriptor
}

// L contains the Logger used to log stuff
var L *logmatic.Logger

// NewPlan returns a pointer to a new Plan
func NewPlan() *Plan {
	return &Plan{
		TempID: fmt.Sprintf("%X", rand.Uint64()),
	}
}

// GetFileList returns a list of the files to edit
func (p *Plan) GetFileList() []string {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}

	if pwd != "" && !p.AbsolutePaths {
		var result []string
		for _, path := range p.InFiles {
			result = append(result, strings.TrimPrefix(path, pwd+string(os.PathSeparator)))
		}
		return result
	}

	return p.InFiles
}

// TempFile returns the path to the temporary file
func (p *Plan) TempFile() string {
	return os.TempDir() + string(os.PathSeparator) + "br_" + p.TempID
}

func init() {
	rand.Seed(time.Since(time.Now()).Nanoseconds())
}

// Execute iterates over the jobs and executes them
func (p *Plan) Execute() (errOccured bool, errorDescs []string, errs []error) {
	for _, job := range p.jobs {
		switch job.Action {
		case -1:
			err := os.Remove(job.SourcePath)
			if err != nil {
				errOccured = true
				errorDescs = append(errorDescs, "Error while deleting "+job.SourcePath)
				errs = append(errs, err)
				continue
			}
		case 3: // the same as 1 but special
			fallthrough
		case 1:
			err := os.Rename(job.SourcePath, job.DstPath)
			if err != nil {
				errOccured = true
				errorDescs = append(errorDescs, "Error while moving "+job.SourcePath+" to "+job.DstPath)
				errs = append(errs, err)
				continue
			}
		case 2:
			err := os.MkdirAll(job.DstPath, 0o777)
			if err != nil {
				errOccured = true
				errorDescs = append(errorDescs, "Error while creating "+job.SourcePath)
				errs = append(errs, err)
				continue
			}
		}
	}
	return
}
