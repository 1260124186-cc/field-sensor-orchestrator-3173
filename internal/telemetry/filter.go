package telemetry

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Filter struct {
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

type FilterDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type FilterSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewFilter(name string, limit int) Filter {
	if limit < 1 {
		limit = 1
	}
	return Filter{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "filter"}, Sequence: []string{"receive", "screen", "route", "emit"}, Stage1: "receive-1", Stage2: "screen-2", Stage3: "route-3", Stage4: "emit-4", Stage5: "receive-5", Stage6: "screen-6", Stage7: "route-7", Stage8: "emit-8", Stage9: "receive-9", Stage10: "screen-10", Stage11: "route-11", Stage12: "emit-12", Stage13: "receive-13", Stage14: "screen-14", Stage15: "route-15", Stage16: "emit-16", Stage17: "receive-17", Stage18: "screen-18"}
}

func (v Filter) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("filter stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("filter name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("filter bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("filter sequence is incomplete")
	}
	return nil
}

func (v Filter) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+6)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6}...)
	sort.Strings(out)
	return out
}

func (v Filter) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Filter) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Filter) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Filter) Decide(step string, used int, at time.Time) FilterDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6}...)
	return FilterDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Filter) Snapshot(at time.Time) FilterSnapshot {
	return FilterSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 6}
}
