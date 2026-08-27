package lifecycle

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CycleLifecycle struct {
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

type CycleLifecycleDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type CycleLifecycleSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewCycles(name string, limit int) CycleLifecycle {
	if limit < 1 {
		limit = 1
	}
	return CycleLifecycle{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "cycle lifecycle"}, Sequence: []string{"queue", "lease", "run", "complete"}, Stage1: "queue-1", Stage2: "lease-2", Stage3: "run-3", Stage4: "complete-4", Stage5: "queue-5", Stage6: "lease-6", Stage7: "run-7", Stage8: "complete-8", Stage9: "queue-9", Stage10: "lease-10", Stage11: "run-11", Stage12: "complete-12", Stage13: "queue-13", Stage14: "lease-14", Stage15: "run-15", Stage16: "complete-16", Stage17: "queue-17", Stage18: "lease-18"}
}

func (v CycleLifecycle) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("cycle lifecycle stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("cycle lifecycle name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("cycle lifecycle bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("cycle lifecycle sequence is incomplete")
	}
	return nil
}

func (v CycleLifecycle) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+8)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8}...)
	sort.Strings(out)
	return out
}

func (v CycleLifecycle) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v CycleLifecycle) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v CycleLifecycle) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v CycleLifecycle) Decide(step string, used int, at time.Time) CycleLifecycleDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8}...)
	return CycleLifecycleDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v CycleLifecycle) Snapshot(at time.Time) CycleLifecycleSnapshot {
	return CycleLifecycleSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 8}
}
