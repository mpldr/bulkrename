package plan

import (
	"fmt"
	"math/rand"
	"time"
)

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
	// CreateDirs indicates whether non-existant directories should be created as needed
	CreateDirs bool
	// StopToShow indicates whether an overview of the applied actions should be shown and confirmation requested
	StopToShow bool
	// DeleteEmpty indicates whether files corresponding with empty lines should be deleted
	DeleteEmpty bool
}

// NewPlan returns a pointer to a new Plan
func NewPlan() *Plan {
	return &Plan{
		TempID: fmt.Sprintf("%X", rand.Uint64()),
	}
}

func init() {
	rand.Seed(time.Since(time.Now()).Nanoseconds())
}
