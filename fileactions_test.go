package main

import (
	"os"
	"testing"

	"git.sr.ht/~poldi1405/glog"
)

func TestRemoveInvalidEntries(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	filelist := []string{
		"test/ok",
		"test/noexist&/&%",
		"test/not_allowed/permdenied",
	}

	if err := os.MkdirAll("test/not_allowed", 0o700); err != nil {
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

	err = os.Chmod("test/not_allowed", 0o000)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.Chmod("test/not_allowed", 0o700)
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
