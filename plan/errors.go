package plan

import "errors"

var (
	dirCreationNotAllowed    = errors.New("not allowed to create directory")
	multipleChoiceNotAllowed = errors.New("not allowed to replace file with directory")
)
