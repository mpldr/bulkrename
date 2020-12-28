package plan

import (
	"bufio"
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
	p.listAllFiles("test")
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

func TestWriteTempFileFails(t *testing.T) {
	var p Plan
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)

	if err := os.MkdirAll("test/not_allowed", 0700); err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll("test")
		if err != nil {
			t.Error(err)
		}
	}()

	err := os.Chmod("test/not_allowed", 0000)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.Chmod("test/not_allowed", 0700)
		if err != nil {
			t.Error(err)
		}
	}()

	tmpdir := os.Getenv("TMPDIR")
	os.Setenv("TMPDIR", "test/not_allowed")
	defer func() { os.Setenv("TMPDIR", tmpdir) }()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Create("test/not_allowed/ok")
	if err == nil {
		t.Skip()
	}

	err = p.writeTempFile()
	if err == nil {
		t.Fail()
	}
}

func TestWriteTempFile(t *testing.T) {
	var p Plan
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)

	if err := os.MkdirAll("test", 0700); err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll("test")
		if err != nil {
			t.Error(err)
		}
	}()

	tmpdir := os.Getenv("TMPDIR")
	os.Setenv("TMPDIR", "test")
	defer func() { os.Setenv("TMPDIR", tmpdir) }()

	p.InFiles = []string{
		"some filename",
		"ä w€iRd fiłenæm",
		"",
	}

	err := p.writeTempFile()
	if err != nil {
		t.Fail()
	}

	file, err := os.Open(p.TempFile())
	if err != nil {
		t.Fail()
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0

	for scanner.Scan() {
		if scanner.Text() != p.InFiles[i] {
			t.Fail()
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		t.Fail()
	}
}

func TestLoadFileList(t *testing.T) {
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

	_, err = os.Create("test/ok2")
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
	filelist := []string{"test/ok", "test/ok2", "test/allowed_but_empty"}
	p.LoadFileList(filelist, false)
	cwd, err := os.Getwd()
	if err != nil {
		t.Log(err)
		t.Skip()
	}

	m := make(map[string]bool)
	for i := 0; i < len(filelist); i++ {
		m[cwd+string(os.PathSeparator)+filelist[i]] = true
	}

	for _, f := range p.InFiles {
		if _, ok := m[f]; !ok {
			t.Log(m)
			t.Log("found unknown " + f)
			t.Fail()
		}
	}
}

func TestLoadFileListFails(t *testing.T) {
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

	_, err = os.Create("test/ok2")
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
	filelist := []string{"test/not_allowed/permdenied"}
	p.LoadFileList(filelist, false)
	cwd, err := os.Getwd()
	if err != nil {
		t.Log(err)
		t.Skip()
	}

	m := make(map[string]bool)
	for i := 0; i < len(filelist); i++ {
		m[cwd+string(os.PathSeparator)+filelist[i]] = true
	}

	if len(p.InFiles) != 0 {
		t.Log(p.InFiles)
		t.Fail()
	}
}
