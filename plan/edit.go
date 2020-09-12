package plan

import (
	"fmt"
	"os"
	"os/exec"

	"gitlab.com/poldi1405/go-ansi"
)

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
		fmt.Printf(ansi.Red("PANIC!")+" Cannot start editor!\n\tExecutable: %v\n\tArguments: %v\n%v\n\nOutput:\n%v", p.Editor, p.EditorArgs, err)
		L.Trace("Executable:", p.Editor)
		L.Trace("Arguments:", p.EditorArgs)
		L.Trace("Error:", err)
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
			L.Debug("Replacing", p.EditorArgs[i], "with", val)
			p.EditorArgs[i] = val
		}
	}
}
