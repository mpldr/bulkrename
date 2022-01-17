package plan

import (
	"fmt"
	"os"
	"os/exec"

	"git.sr.ht/~poldi1405/glog"
)

// StartEditing launches the editor and loads the required file for editing
func (p *Plan) StartEditing() error {
	if len(p.InFiles) == 0 {
		glog.Info("No files for editing left, so no editing necessary")
		return nil
	}
	err := p.writeTempFile()
	if err != nil {
		return err
	}

	p.prepareArguments()

	c := exec.Command(p.Editor, p.EditorArgs...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Env = os.Environ()

	err = c.Run()
	if err != nil {
		glog.Trace("Executable:" + p.Editor)
		glog.Trace("Arguments:" + fmt.Sprint(p.EditorArgs))
		glog.Info("Error:" + err.Error())
		glog.Fatal("Cannot start editor!")
		return err
	}

	return p.CreatePlan(p.TempFile())
}

func (p *Plan) prepareArguments() {
	replace := map[string]string{
		"{}": p.TempFile(),
	}

	for i, arg := range p.EditorArgs {
		if val, found := replace[arg]; found {
			glog.Debug("Replacing " + p.EditorArgs[i] + " with " + val)
			p.EditorArgs[i] = val
		}
	}
}
