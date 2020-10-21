package plan

import (
	"os"
	"testing"

	"github.com/mborders/logmatic"
	j "gitlab.com/poldi1405/bulkrename/plan/jobdescriptor"
)

func TestTempFileRemoved(t *testing.T) {
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)
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
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)
	var p Plan

	p.InFiles = []string{"1", "2", "3"}
	p.OutFiles = []string{"2", "3", "1"}
	p.jobs = []j.JobDescriptor{
		j.JobDescriptor{
			Action:     0,
			SourcePath: "1",
			DstPath:    "2",
		},
		j.JobDescriptor{
			Action:     0,
			SourcePath: "2",
			DstPath:    "3",
		},
		j.JobDescriptor{
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
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)
	var p Plan

	p.InFiles = []string{"1", "2", "3"}
	p.OutFiles = []string{"2", "3", "1"}
	p.jobs = []j.JobDescriptor{
		j.JobDescriptor{
			Action:     0,
			SourcePath: "1",
			DstPath:    "2",
		},
		j.JobDescriptor{
			Action:     0,
			SourcePath: "2",
			DstPath:    "3",
		},
		j.JobDescriptor{
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
	L = logmatic.NewLogger()
	L.SetLevel(logmatic.LogLevel(42))
	L.SetLevel(logmatic.FATAL)
	var p Plan

	p.InFiles = []string{"1", "2", "3"}
	p.OutFiles = []string{"2", "3", "1"}
	p.jobs = []j.JobDescriptor{
		j.JobDescriptor{
			Action:     0,
			SourcePath: "1",
			DstPath:    "4",
		},
		j.JobDescriptor{
			Action:     0,
			SourcePath: "2",
			DstPath:    "5",
		},
		j.JobDescriptor{
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
