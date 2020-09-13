package main

import (
	"bufio"
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/mborders/logmatic"

	"gitlab.com/poldi1405/bulkrename/plan"
	"gitlab.com/poldi1405/go-ansi"
)

var (
	buildVersion = "Version unknown! This was not built using the Makefile!"
	absolute     *bool
	recursive    *bool
	check        *bool
	mkdir        bool
	editor       *string
	loglevel     *string
	args         *[]string
	overwrite    bool
	delem        *bool
	files        *[]string
	l            *logmatic.Logger
)

func main() {
	br := cli.App("br", "Rename files in a bulk")
	l = logmatic.NewLogger()
	trace := os.Getenv("BR_ENABLE_TRACE")
	if len(trace) > 0 {
		l.SetLevel(logmatic.TRACE)
		l.Debug("LogLevel set to TRACE")
	} else {
		l.SetLevel(logmatic.WARN)
	}

	setupCLI(br)

	br.Action = func() {
		switch *loglevel {
		case "trace":
			l.Debug("Set LogLevel to TRACE")
			l.SetLevel(logmatic.TRACE)

		case "debug":
			l.Debug("Set LogLevel to DEBUG")
			l.SetLevel(logmatic.DEBUG)

		case "info":
			l.Debug("Set LogLevel to INFO")
			l.SetLevel(logmatic.INFO)

		case "error":
			l.Debug("Set LogLevel to ERROR")
			l.SetLevel(logmatic.ERROR)

		case "fatal":
			l.Debug("Set LogLevel to FATAL")
			l.SetLevel(logmatic.FATAL)
		}
		if len(trace) > 0 {
			l.SetLevel(logmatic.TRACE)
			l.Debug("Reset LogLevel to TRACE because BR_ENABLE_TRACE is set")
		}

		plan.L = l

		l.Info("setting up plan")
		jobplan := plan.NewPlan()

		jobplan.AbsolutePaths = *absolute
		l.Debug("set AbsolutePaths to", *absolute)
		jobplan.Overwrite = overwrite
		l.Debug("set Overwrite to", overwrite)
		jobplan.Editor = *editor
		l.Debug("set Editor to", *editor)
		jobplan.EditorArgs = *args
		l.Debug("set EditorArgs to", *args)
		jobplan.CreateDirs = mkdir
		l.Debug("set CreateDirs to", mkdir)
		jobplan.StopToShow = *check
		l.Debug("set StopToShow to", *check)
		jobplan.DeleteEmpty = *delem
		l.Debug("set DeleteEmpty to", *delem)

		l.Info("cleaning input")
		*files = RemoveInvalidEntries(*files)
		l.Info("loading filelist")
		jobplan.LoadFileList(*files, *recursive)
		l.Info("starting editor")
		err := jobplan.StartEditing()
		if err != nil {
			l.Fatal("error occured when editing", err)
		}

		err = jobplan.PrepareExecution()
		if err != nil {
			os.Exit(1)
		}

		if jobplan.StopToShow {
			jobplan.PreviewPlan()
			fmt.Print("\nDo you wish to continue? [Y/n] ")
			reader := bufio.NewReader(os.Stdin)
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			switch char {
			case 'n', 'N':
				os.Exit(0)
			}
		}

		errOcc, msgs, errs := jobplan.Execute()
		if errOcc {
			fmt.Print(ansi.Bold(ansi.Red("ERROR!")), "\nThe following errors occures while executing the plan:\n\n")

			for i, msg := range msgs {
				fmt.Println(msg, errs[i])
			}
			os.Exit(1)
		}

		//fmt.Println(jobplan.InFiles)
		//fmt.Printf("recursive: %v\nabsolute: %v\nstop to show: %v\ncreate directories: %v\nuse editor: %v\narguemnts: %v\noverwrite: %v\ndelete empty: %v\nfiles: %v", *recursive, *absolute, *check, mkdir, *editor, *args, overwrite, *delem, *files)
	}
	err := br.Run(os.Args)
	if err != nil {
		l.Fatal("unable to execute", err)
	}
}

func setupCLI(br *cli.Cli) {
	br.Version("v version", "bulkrename "+buildVersion)
	br.Spec = "[-r] [-a] [-d] [--editor] [--arg...] [--check] [--no-mkdir] [--no-overwrite] [--loglevel] FILES..."

	recursive = br.Bool(cli.BoolOpt{
		Name:   "r recursive",
		Value:  false,
		Desc:   "recursively list files",
		EnvVar: "BR_RECURSIVE",
	})

	absolute = br.Bool(cli.BoolOpt{
		Name:  "a absolute",
		Value: false,
		Desc:  "list files with their absolute path",
	})

	check = br.Bool(cli.BoolOpt{
		Name:  "check",
		Value: false,
		Desc:  "show actions that will be performed",
	})

	nomkdir := br.Bool(cli.BoolOpt{
		Name:  "no-mkdir",
		Value: false,
		Desc:  "do not create directories that do not exist",
	})
	mkdir = !*nomkdir

	nooverwrite := br.Bool(cli.BoolOpt{
		Name:  "no-overwrite",
		Value: false,
		Desc:  "do not overwrite files",
	})
	overwrite = !*nooverwrite

	delem = br.Bool(cli.BoolOpt{
		Name:  "d delete-empty",
		Desc:  "delete files that were renamed to empty strings",
		Value: false,
	})

	loglevel = br.String(cli.StringOpt{
		Name:  "loglevel",
		Desc:  "set the loglevel",
		Value: "warn",
	})

	editor = br.String(cli.StringOpt{
		Name:   "editor",
		Desc:   "executable of the editor",
		Value:  DefaultEditor,
		EnvVar: "EDITOR",
	})

	args = br.Strings(cli.StringsOpt{
		Name:  "arg",
		Desc:  "arguments for the editor",
		Value: []string{"{}"},
	})

	files = br.Strings(cli.StringsArg{
		Name: "FILES",
		Desc: "the source files that will be added to the editor",
	})

}
