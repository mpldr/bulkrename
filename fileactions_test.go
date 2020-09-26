package main

import (
	"os"
	"testing"

	"github.com/mborders/logmatic"
)

func TestRemoveInvalidEntries(t *testing.T) {
	l = logmatic.NewLogger()
	l.SetLevel(logmatic.LogLevel(42))
	l.SetLevel(logmatic.FATAL)
	filelist := []string{
		"test/ok",
		"test/noexist&/&%",
		"test/not_allowed/permdenied",
	}

	if err := os.MkdirAll("test/not_allowed", 0700); err != nil {
		t.Error(err)
	}

	defer func() {
		err := os.RemoveAll("test")
		if err != nil {
			t.Error(err)
		}
	}()

	_, err := os.Create("test/ok")
	if err != nil {
		t.Error(err)
	}

	_, err = os.Create("test/not_allowed/permdenied")
	if err != nil {
		t.Error(err)
	}

	err = os.Chmod("test/not_allowed", 0000)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.Chmod("test/not_allowed", 0700)
		if err != nil {
			t.Error(err)
		}
	}()

	result := RemoveInvalidEntries(filelist)

	if len(result) != 1 {
		if len(result) > 1 && result[1] == "test/not_allowed/permdenied" {
			t.Skipf("seems like Chmod failed. Skipping test.")
			return
		}
		t.Log(result)
		t.Error("list too long")
	}

	if result[0] != "test/ok" {
		t.Log(result)
		t.Error("wrong file kept")
	}
}
