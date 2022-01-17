package plan

import (
	"os"
	"testing"

	j "git.sr.ht/~poldi1405/bulkrename/plan/jobdescriptor"
	"git.sr.ht/~poldi1405/glog"
)

func TestTempFileRemoved(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan

	err := p.CreatePlan("probably does not exist. If it does, delete it.")
	if err == nil {
		t.Fail()
	}

	if !os.IsNotExist(err) {
		t.Fail()
	}
}

func TestDetectCircles(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan

	p.InFiles = []string{"1", "2", "3"}
	p.OutFiles = []string{"2", "3", "1"}
	p.jobs = []j.JobDescriptor{
		{
			Action:     0,
			SourcePath: "1",
			DstPath:    "2",
		},
		{
			Action:     0,
			SourcePath: "2",
			DstPath:    "3",
		},
		{
			Action:     0,
			SourcePath: "3",
			DstPath:    "1",
		},
	}

	results := p.findCollisions()
	if len(results) != 3 {
		t.Error("got", len(results), "prerules generated instead of 3")
	}
}

func TestDetectLinearReplace(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan

	p.InFiles = []string{"1", "2", "3"}
	p.OutFiles = []string{"2", "3", "1"}
	p.jobs = []j.JobDescriptor{
		{
			Action:     0,
			SourcePath: "1",
			DstPath:    "2",
		},
		{
			Action:     0,
			SourcePath: "2",
			DstPath:    "3",
		},
		{
			Action:     0,
			SourcePath: "3",
			DstPath:    "4",
		},
	}

	results := p.findCollisions()
	if len(results) != 2 {
		t.Error("got", len(results), "prerules generated instead of 3")
	}
}

func TestDetectNoCircles(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan

	p.InFiles = []string{"1", "2", "3"}
	p.OutFiles = []string{"2", "3", "1"}
	p.jobs = []j.JobDescriptor{
		{
			Action:     0,
			SourcePath: "1",
			DstPath:    "4",
		},
		{
			Action:     0,
			SourcePath: "2",
			DstPath:    "5",
		},
		{
			Action:     0,
			SourcePath: "3",
			DstPath:    "6",
		},
	}

	results := p.findCollisions()
	if len(results) != 0 {
		t.Error("got", len(results), "prerules generated instead of 0")
	}
}

func TestFailGetAbsolutePath(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Skipf(err.Error())
	}

	err = os.Mkdir(pwd+"test_failabspath", 0o777)
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove(pwd + "test_failabspath")

	err = os.Chdir(pwd + "test_failabspath")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer func() {
		err := os.Chdir(pwd)
		if err != nil {
			t.Log(err)
		}
	}()

	f, err := os.Create(pwd + "test.txt")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove(pwd + "test.txt")
	defer f.Close()
	_, err = f.WriteString("Hello World")
	if err != nil {
		t.Skipf(err.Error())
	}

	err = os.Remove(pwd + "test_failabspath")
	if err != nil {
		t.Skipf(err.Error())
	}
	var p Plan
	err = p.CreatePlan(pwd + "test.txt")
	if err == nil {
		t.Fail()
	}
}

func TestDeleteEmptyLines(t *testing.T) {
	_, err := os.Create("test.txt")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("test.txt")

	var p Plan
	p.InFiles = []string{"hey there!"}
	p.DeleteEmpty = true

	err = p.CreatePlan("test.txt")
	if err != nil {
		t.Fail()
	}

	if len(p.jobs) != 1 {
		t.FailNow()
	}

	if (p.jobs[0]).Action != -1 {
		t.Fail()
	}
}

func TestNoUnnecessaryPrerules(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan

	p.InFiles = []string{"1"}
	p.OutFiles = []string{"2"}
	p.jobs = []j.JobDescriptor{
		{
			Action:     0,
			SourcePath: "1",
			DstPath:    "2",
		},
	}

	_, err := os.Create("1")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("1")

	err = p.PrepareExecution()
	if err != nil {
		t.Error(err)
	}

	if len(p.jobs) != 1 {
		t.Error("got", len(p.jobs)-1, "prejobs")
	}
}

func TestFailBecauseActionForbidden(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan
	reset := func() {
		p.InFiles = []string{"1"}
		p.OutFiles = []string{"3/2/1"}
		p.jobs = []j.JobDescriptor{
			{
				Action:     0,
				SourcePath: "1",
				DstPath:    "3/1",
			},
		}
	}
	reset()

	_, err := os.Create("1")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("1")

	_, err = os.Create("3")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("3")

	p.Overwrite = false
	p.CreateDirs = false
	err = p.PrepareExecution()
	if err != errMultipleChoiceNotAllowed {
		t.Error("did not fail when overwriting and creating directories is forbidden")
	}

	reset()

	p.Overwrite = false
	p.CreateDirs = true
	err = p.PrepareExecution()
	if err != errMultipleChoiceNotAllowed {
		t.Error("did not fail when only overwriting is forbidden")
		if err != nil {
			t.Log(err)
		}
	}

	reset()

	p.Overwrite = true
	p.CreateDirs = false
	err = p.PrepareExecution()
	if err == errDirCreationNotAllowed {
		t.Error("did not fail when only creating directories is forbidden")
	}
}

func TestFailBecauseMkdirForbidden(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan
	reset := func() {
		p.InFiles = []string{"1"}
		p.OutFiles = []string{"2/1"}
		p.jobs = []j.JobDescriptor{
			{
				Action:     0,
				SourcePath: "1",
				DstPath:    "2/1",
			},
		}
	}
	reset()

	_, err := os.Create("1")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("1")

	p.Overwrite = true
	p.CreateDirs = false
	err = p.PrepareExecution()
	if err != errDirCreationNotAllowed {
		t.Error("did not fail when directories is forbidden")
	}
}

func TestCreateMkdirPrerule(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan
	reset := func() {
		p.InFiles = []string{"1"}
		p.OutFiles = []string{"2/1"}
		p.jobs = []j.JobDescriptor{
			{
				Action:     0,
				SourcePath: "1",
				DstPath:    "2/1",
			},
		}
	}
	reset()

	_, err := os.Create("1")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("1")

	p.Overwrite = true
	p.CreateDirs = true
	err = p.PrepareExecution()
	if err != nil {
		t.Error(err)
	}

	if len(p.jobs) != 2 || p.jobs[0].Action != 2 {
		t.Fail()
	}
}

func TestReplaceFileWithDirectoryPrerules(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan
	reset := func() {
		p.InFiles = []string{"1"}
		p.OutFiles = []string{"3/1"}
		p.jobs = []j.JobDescriptor{
			{
				Action:     0,
				SourcePath: "1",
				DstPath:    "3/1",
			},
		}
	}
	reset()

	_, err := os.Create("1")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("1")

	_, err = os.Create("3")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("3")

	p.Overwrite = true
	p.CreateDirs = true
	err = p.PrepareExecution()
	if err != nil {
		t.Error(err)
	}

	if len(p.jobs) != 3 || p.jobs[0].Action != -1 || p.jobs[1].Action != 2 {
		t.Fail()
	}
}

func TestIgnoreRingDetectionRules(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan
	reset := func() {
		p.InFiles = []string{"1"}
		p.OutFiles = []string{"3/1"}
		p.jobs = []j.JobDescriptor{
			{
				Action:     3,
				SourcePath: "1",
				DstPath:    "3/1",
			},
		}
	}
	reset()

	_, err := os.Create("1")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("1")

	_, err = os.Create("3")
	if err != nil {
		t.Skipf(err.Error())
	}
	defer os.Remove("3")

	p.Overwrite = true
	p.CreateDirs = true
	err = p.PrepareExecution()
	if err != nil {
		t.Error(err)
	}

	if len(p.jobs) != 1 {
		t.Fail()
	}
}

func TestPlanningSourceFileNotExist(t *testing.T) {
	glog.SetLevel(glog.Level(42))
	var p Plan
	reset := func() {
		p.InFiles = []string{"1"}
		p.OutFiles = []string{"3/1"}
		p.jobs = []j.JobDescriptor{
			{
				Action:     0,
				SourcePath: "1",
				DstPath:    "3/1",
			},
		}
	}
	reset()

	p.Overwrite = true
	p.CreateDirs = true
	err := p.PrepareExecution()
	if !os.IsNotExist(err) {
		t.Fail()
	}
}
