package main

import (
	"flag"
	. "testing"

	"github.com/mborders/logmatic"

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

func TestSetupLoggingDefault(t *T) {
	l = logmatic.NewLogger()

	if setupLogging() != logmatic.WARN {
		t.Fail()
	}
}

func TestSetupLoggingValues(t *T) {
	l = logmatic.NewLogger()
	*loglevel = "trace"

	if setupLogging() != logmatic.TRACE {
		t.Fail()
	}
	*loglevel = "debug"

	if setupLogging() != logmatic.DEBUG {
		t.Fail()
	}
	*loglevel = "info"

	if setupLogging() != logmatic.INFO {
		t.Fail()
	}
	*loglevel = "error"

	if setupLogging() != logmatic.ERROR {
		t.Fail()
	}
	*loglevel = "fatal"

	if setupLogging() != logmatic.FATAL {
		t.Fail()
	}
	*loglevel = "other"

	if setupLogging() != logmatic.WARN {
		t.Fail()
	}
}
