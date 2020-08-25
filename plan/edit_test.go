package plan

import (
	"reflect"
	"testing"
)

func TestPrepareArguments(t *testing.T) {
	var p Plan
	p.EditorArgs = []string{
		"{}",
		"123456",
		"aasddnew",
	}

	p.prepareArguments()

	if !reflect.DeepEqual(p.EditorArgs, []string{p.TempFile(), "123456", "aasddnew"}) {
		t.Fail()
	}
}
