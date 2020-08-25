package plan

import (
	"fmt"
	"os"
	"os/exec"

	"gitlab.com/poldi1405/go-ansi"
)

func (p *Plan) StartEditing() error {
	if len(p.InFiles) == 0 {
		return nil
	}
	err := p.writeTempFile()
	if err != nil {
		fmt.Printf("%v Cannot write temporary file!\n\n %v\n", ansi.Red("PANIC!"), err)
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
			p.EditorArgs[i] = val
		}
	}
}
