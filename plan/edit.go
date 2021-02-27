package plan

import (
	"fmt"
	"os"
	"os/exec"
)

// StartEditing launches the editor and loads the required file for editing
func (p *Plan) StartEditing() error {
	if len(p.InFiles) == 0 {
		L.Info("No files for editing left, so no editing necessary")
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
		L.Trace("Executable:" + p.Editor)
		L.Trace("Arguments:" + fmt.Sprint(p.EditorArgs))
		L.Info("Error:" + err.Error())
		L.Fatal("Cannot start editor!")
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
			L.Debug("Replacing " + p.EditorArgs[i] + " with " + val)
			p.EditorArgs[i] = val
		}
	}
}
