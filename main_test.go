package main

import (
	"flag"
	"reflect"
	. "testing"

	"git.sr.ht/~poldi1405/bulkrename/plan"

	cli "github.com/jawher/mow.cli"
)

func TestCLISetup(t *T) {
	validCalls := [][]string{
		{"jt", "abc", "def", "ghi"},
		{"jt", "-r", "-a", "Makefile"},
		{"jt", "-r", "Makefile"},
		{"jt", "--editor", "word.exe", "Makefile"},
		{"jt", "--editor", "word.exe", "-a", "Makefile"},
		{"jt", "--check", "--no-mkdir", "--no-overwrite", "Makefile"},
	}
	/*invalidCalls := [][]string{
		[]string{"jt"},
		[]string{"jt", "-r", "-a"},
		[]string{"jt", "--not-a-real-argument", "abc"},
		[]string{"jt", ""},
	}*/ //postponed until there is a way to suppress help message being printed

	test := cli.App("jt", "just a test app")
	test.Cmd.ErrorHandling = flag.ContinueOnError
	setupCLI(test)

	for _, v := range validCalls {
		// no action necessary as this is just a test
		test.Action = func() {}

		err := test.Run(v)
		if err != nil {
			t.Error(v)
		}
	}

	/*for _, v := range invalidCalls {
		test := cli.App("jt", "just a test app")
		setupCLI(test)

		err := test.Run(v)
		if err == nil {
			t.Fail()
		}
	}*/
}

func TestSetupPlan(t *T) {
	p := plan.NewPlan()

	*absolute = true
	overwrite = true
	*editor = "notepad++"
	*args = []string{"my", "list", "of", "args"}
	mkdir = true
	*check = true
	*delem = true

	setJobplanSettings(p)
	if p.AbsolutePaths != *absolute || p.Overwrite != overwrite ||
		p.Editor != *editor || !reflect.DeepEqual(p.EditorArgs, *args) ||
		p.CreateDirs != mkdir || p.StopToShow != *check ||
		p.DeleteEmpty != *delem {
		t.Fail()
	}

	*absolute = false
	overwrite = false
	*editor = "notepad++"
	*args = []string{"this", "is", "a", "different", "list"}
	mkdir = false
	*check = false
	*delem = false

	setJobplanSettings(p)
	if p.AbsolutePaths != *absolute || p.Overwrite != overwrite ||
		p.Editor != *editor || !reflect.DeepEqual(p.EditorArgs, *args) ||
		p.CreateDirs != mkdir || p.StopToShow != *check ||
		p.DeleteEmpty != *delem {
		t.Fail()
	}
}
