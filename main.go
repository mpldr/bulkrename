// Package main is for creating an executable
package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/mborders/logmatic"

	"mpldr.codes/br/plan"
)

var (
	buildVersion = "Version unknown! This was not built using the Makefile!"
	absolute     *bool
	recursive    *bool
	check        *bool
	mkdir        bool
	macro        *string
	editor       *string
	loglevel     *string
	args         *[]string
	overwrite    bool
	delem        *bool
	files        *[]string
	l            *logmatic.Logger
	//go:embed licensetext.txt
	licensetext []byte
	//go:embed LICENSE
	mpl []byte
)

func main() {
	br := cli.App("br", "Rename files in a bulk")
	l = logmatic.NewLogger()

	l.SetLevel(setupLogging())

	setupCLI(br)

	br.Action = func() {
		l.SetLevel(setupLogging())

		plan.L = l

		l.Info("setting up plan")
		jobplan := plan.NewPlan()
		setJobplanSettings(jobplan)

		l.Info("cleaning input")
		*files = RemoveInvalidEntries(*files)
		l.Info("loading filelist")
		jobplan.LoadFileList(*files, *recursive)
		l.Info("starting editor")
		err := jobplan.StartEditing()
		if err != nil {
			l.Info(err.Error())
			l.Fatal("error occurred when editing")
		}

		err = jobplan.PrepareExecution()
		if err != nil {
			l.Info(err.Error())
			l.Fatal("error occurred when preparing execution")
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
			l.Error("Errors occurred while executing the plan")

			for i, msg := range msgs {
				l.Info(msg)
				l.Info(errs[i].Error())
			}
			os.Exit(1)
		}

		// fmt.Println(jobplan.InFiles)
		// fmt.Printf("recursive: %v\nabsolute: %v\nstop to show: %v\ncreate directories: %v\nuse editor: %v\narguemnts: %v\noverwrite: %v\ndelete empty: %v\nfiles: %v", *recursive, *absolute, *check, mkdir, *editor, *args, overwrite, *delem, *files)
	}
	err := br.Run(os.Args)
	if err != nil {
		l.Info(err.Error())
		l.Fatal("unable to execute")
	}
}

func setupCLI(br *cli.Cli) {
	br.Version("v version", "bulkrename "+buildVersion)
	br.Spec = "[-r] [-a] [-d] [--macro | --editor --arg...] [--check] [--no-mkdir] [--no-overwrite] [--loglevel] [FILES...]"

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

	macro = br.String(cli.StringOpt{
		Name:   "macro",
		Desc:   "prepared macro to apply to the file",
		EnvVar: "MACRO",
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

	br.Command("licenses", "print license information", func(s *cli.Cmd) {
		s.Action = func() {
			fmt.Println(string(mpl))
			fmt.Println(string(licensetext))
		}
	})
}

func setupLogging() logmatic.LogLevel {
	trace := os.Getenv("BR_ENABLE_TRACE")
	if len(trace) > 0 {
		l.Debug("LogLevel set to TRACE")
		return logmatic.TRACE
	}
	if loglevel == nil || *loglevel == "" {
		return logmatic.WARN
	}

	switch *loglevel {
	case "trace":
		l.Debug("Set LogLevel to TRACE")
		return logmatic.TRACE

	case "debug":
		l.Debug("Set LogLevel to DEBUG")
		return logmatic.DEBUG

	case "info":
		l.Debug("Set LogLevel to INFO")
		return logmatic.INFO

	case "error":
		l.Debug("Set LogLevel to ERROR")
		return logmatic.ERROR

	case "fatal":
		l.Debug("Set LogLevel to FATAL")
		return logmatic.FATAL
	default:
		l.Debug("Set LogLevel to WARN")
		return logmatic.WARN
	}
}

func setJobplanSettings(jobplan *plan.Plan) {
	jobplan.AbsolutePaths = *absolute
	l.Debug("set AbsolutePaths to " + strconv.FormatBool(*absolute))
	jobplan.Overwrite = overwrite
	l.Debug("set Overwrite to " + strconv.FormatBool(overwrite))
	jobplan.Editor = *editor
	l.Debug("set Editor to " + *editor)
	jobplan.EditorArgs = *args
	l.Debug("set EditorArgs to " + strings.Join(*args, ", "))
	jobplan.CreateDirs = mkdir
	l.Debug("set CreateDirs to " + strconv.FormatBool(mkdir))
	jobplan.StopToShow = *check
	l.Debug("set StopToShow to " + strconv.FormatBool(*check))
	jobplan.DeleteEmpty = *delem
	l.Debug("set DeleteEmpty to " + strconv.FormatBool(*delem))
}
