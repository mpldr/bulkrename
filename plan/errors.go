package plan

import "errors"

var (
	errDirCreationNotAllowed    = errors.New("not allowed to create directory")
	errMultipleChoiceNotAllowed = errors.New("job not possible with current restrictions")
)
