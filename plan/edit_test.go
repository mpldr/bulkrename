package plan

import (
	"os"
	"reflect"
	"testing"

	"git.sr.ht/~poldi1405/glog"
)

func TestPrepareArguments(t *testing.T) {
	glog.SetLevel(glog.Level(42))
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
	glog.SetLevel(glog.Level(42))

	p := NewPlan()
	p.Editor = "true"
	err := p.StartEditing()
	if err != nil {
		t.Fail()
	}
}

func TestEditSuccess(t *testing.T) {
	glog.SetLevel(glog.Level(42))

	p := NewPlan()
	p.Editor = "true"
	p.InFiles = []string{"ok"}
	err := p.StartEditing()
	if err != nil {
		t.Fail()
	}
}

func TestEditFail(t *testing.T) {
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
