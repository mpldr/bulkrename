package plan

import (
	"os"
	"reflect"
	"testing"

	"github.com/mborders/logmatic"
)

func TestPrepareArguments(t *testing.T) {
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)
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

func TestEditEmpty(t *testing.T) {
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)

	p := NewPlan()
	p.Editor = "true"
	err := p.StartEditing()
	if err != nil {
		t.Fail()
	}
}
func TestEditSuccess(t *testing.T) {
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)

	p := NewPlan()
	p.Editor = "true"
	p.InFiles = []string{"ok"}
	err := p.StartEditing()
	if err != nil {
		t.Fail()
	}
}
func TestEditFail(t *testing.T) {
	L = logmatic.NewLogger()
	L.ExitOnFatal = false

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = os.NewFile(0, os.DevNull)

	p := NewPlan()
	p.Editor = "false"
	p.InFiles = []string{"ok"}
	err := p.StartEditing()
	if err == nil {
		t.Fail()
	}
}
