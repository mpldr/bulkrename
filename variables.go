package main

import "os"

var (
	// TempPath contains the temporary filepath excluding the ID which is provided by the plan
	TempPath = os.TempDir() + string(os.PathSeparator) + "br_"
)
