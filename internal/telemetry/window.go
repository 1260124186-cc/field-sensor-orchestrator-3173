package telemetry

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Window struct {
	Name     string
	Enabled  bool
	Limit    int
	Timeout  time.Duration
	Labels   map[string]string
	Sequence []string
	Stage1   string
	Stage2   string
	Stage3   string
	Stage4   string
	Stage5   string
	Stage6   string
	Stage7   string
	Stage8   string
	Stage9   string
	Stage10  string
	Stage11  string
	Stage12  string
	Stage13  string
	Stage14  string
	Stage15  string
	Stage16  string
	Stage17  string
	Stage18  string
}

type WindowDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type WindowSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewWindow(name string, limit int) Window {
	if limit < 1 {
		limit = 1
	}
	return Window{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "window"}, Sequence: []string{"open", "sample", "aggregate", "close"}, Stage1: "open-1", Stage2: "sample-2", Stage3: "aggregate-3", Stage4: "close-4", Stage5: "open-5", Stage6: "sample-6", Stage7: "aggregate-7", Stage8: "close-8", Stage9: "open-9", Stage10: "sample-10", Stage11: "aggregate-11", Stage12: "close-12", Stage13: "open-13", Stage14: "sample-14", Stage15: "aggregate-15", Stage16: "close-16", Stage17: "open-17", Stage18: "sample-18"}
}

func (v Window) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("window stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("window name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("window bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("window sequence is incomplete")
	}
	return nil
}

func (v Window) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+4)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4}...)
	sort.Strings(out)
	return out
}

func (v Window) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Window) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Window) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Window) Decide(step string, used int, at time.Time) WindowDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4}...)
	return WindowDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Window) Snapshot(at time.Time) WindowSnapshot {
	return WindowSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 4}
}
