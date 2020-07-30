package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

var (
	buildVersion = "Version unknown! This was not built using the Makefile!"
	absolute     *bool
	recursive    *bool
	check        *bool
	mkdir        bool
	editor       *string
	overwrite    bool
	delem        *bool
	files        *[]string
)

func main() {
	br := cli.App("br", "Rename files in a bulk")

	setupCLI(br)

	br.Action = func() {
		fmt.Printf("recursive: %v\nabsolute: %v\nstop to show: %v\ncreate directories: %v\nuse editor: %v\noverwrite: %v\ndelete empty: %v\nfiles: %v", *recursive, *absolute, *check, mkdir, *editor, overwrite, *delem, *files)
	}
	br.Run(os.Args)
}

func setupCLI(br *cli.Cli) {
	br.Version("v version", "bulkrename "+buildVersion)
	br.Spec = "[-r] [-a] [--editor] [--check] [--no-mkdir] [--no-overwrite] FILES..."

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

	editor = br.String(cli.StringOpt{
		Name:   "editor",
		Desc:   "executable of the editor",
		Value:  DefaultEditor,
		EnvVar: "EDITOR",
	})

	files = br.Strings(cli.StringsArg{
		Name: "FILES",
		Desc: "the source files that will be added to the editor",
	})

}
