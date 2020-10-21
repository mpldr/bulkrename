package plan

import "errors"

var (
	dirCreationNotAllowed    = errors.New("not allowed to create directory")
	multipleChoiceNotAllowed = errors.New("job not possible with current restrictions")
)
