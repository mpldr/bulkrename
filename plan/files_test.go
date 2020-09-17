package plan

import (
	"os"
	"testing"

	"github.com/mborders/logmatic"
)

func TestRecursiveFileList(t *testing.T) {
	var p Plan
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)

	if err := os.MkdirAll("test/not_allowed", 0700); err != nil {
		t.Error(err)
	}

	if err := os.MkdirAll("test/allowed_but_empty", 0700); err != nil {
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

	wg.Add(1)
	_ = p.listAllFiles("test")
	wg.Wait()

	result := p.InFiles

	if len(result) != 2 {
		if len(result) > 2 && result[1] == "test/not_allowed/permdenied" {
			t.Skipf("seems like Chmod failed. Skipping test.")
			return
		}
		t.Log(result)
		t.Error("list length does not match")
		return
	}
}
