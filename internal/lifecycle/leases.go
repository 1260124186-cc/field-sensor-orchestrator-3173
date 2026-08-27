package lifecycle

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type LeaseLifecycle struct {
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

type LeaseLifecycleDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type LeaseLifecycleSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewLeases(name string, limit int) LeaseLifecycle {
	if limit < 1 {
		limit = 1
	}
	return LeaseLifecycle{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "lease lifecycle"}, Sequence: []string{"create", "hold", "release", "close"}, Stage1: "create-1", Stage2: "hold-2", Stage3: "release-3", Stage4: "close-4", Stage5: "create-5", Stage6: "hold-6", Stage7: "release-7", Stage8: "close-8", Stage9: "create-9", Stage10: "hold-10", Stage11: "release-11", Stage12: "close-12", Stage13: "create-13", Stage14: "hold-14", Stage15: "release-15", Stage16: "close-16", Stage17: "create-17", Stage18: "hold-18"}
}

func (v LeaseLifecycle) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("lease lifecycle stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("lease lifecycle name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("lease lifecycle bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("lease lifecycle sequence is incomplete")
	}
	return nil
}

func (v LeaseLifecycle) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+9)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9}...)
	sort.Strings(out)
	return out
}

func (v LeaseLifecycle) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v LeaseLifecycle) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v LeaseLifecycle) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v LeaseLifecycle) Decide(step string, used int, at time.Time) LeaseLifecycleDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9}...)
	return LeaseLifecycleDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v LeaseLifecycle) Snapshot(at time.Time) LeaseLifecycleSnapshot {
	return LeaseLifecycleSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 9}
}
