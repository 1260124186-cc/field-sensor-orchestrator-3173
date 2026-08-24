package app

import (
	"context"
	"strings"
	"testing"
)

func TestDemoWorkflows(t *testing.T) {
	a := NewDemo()
	if _, err := a.RegisterSensor("tester"); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Schedule("tester"); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Run(context.Background(), "tester"); err != nil {
		t.Fatal(err)
	}
	summary, err := a.Summary()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(summary, "completed=1") {
		t.Fatalf("summary=%s", summary)
	}
}
