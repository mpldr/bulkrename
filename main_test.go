package main

import (
	"flag"
	. "testing"

	cli "github.com/jawher/mow.cli"
)

func TestCLISetup(t *T) {
	validCalls := [][]string{
		[]string{"jt", "abc", "def", "ghi"},
		[]string{"jt", "-r", "-a", "Makefile"},
		[]string{"jt", "-r", "Makefile"},
		[]string{"jt", "--editor", "word.exe", "Makefile"},
		[]string{"jt", "--editor", "word.exe", "-a", "Makefile"},
		[]string{"jt", "--check", "--no-mkdir", "--no-overwrite", "Makefile"},
	}
	/*invalidCalls := [][]string{
		[]string{"jt"},
		[]string{"jt", "-r", "-a"},
		[]string{"jt", "--not-a-real-argument", "abc"},
		[]string{"jt", ""},
	}*/ //postponed until there is a way to supress help message being printed

	test := cli.App("jt", "just a test app")
	test.Cmd.ErrorHandling = flag.ContinueOnError
	setupCLI(test)

	for _, v := range validCalls {
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
