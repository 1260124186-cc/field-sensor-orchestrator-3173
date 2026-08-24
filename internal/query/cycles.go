package query

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CycleQuery struct {
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

type CycleQueryDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type CycleQuerySnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewCycles(name string, limit int) CycleQuery {
	if limit < 1 {
		limit = 1
	}
	return CycleQuery{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "cycle query"}, Sequence: []string{"scope", "state", "time", "return"}, Stage1: "scope-1", Stage2: "state-2", Stage3: "time-3", Stage4: "return-4", Stage5: "scope-5", Stage6: "state-6", Stage7: "time-7", Stage8: "return-8", Stage9: "scope-9", Stage10: "state-10", Stage11: "time-11", Stage12: "return-12", Stage13: "scope-13", Stage14: "state-14", Stage15: "time-15", Stage16: "return-16", Stage17: "scope-17", Stage18: "state-18"}
}

func (v CycleQuery) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("cycle query stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("cycle query name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("cycle query bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("cycle query sequence is incomplete")
	}
	return nil
}

func (v CycleQuery) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+11)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11}...)
	sort.Strings(out)
	return out
}

func (v CycleQuery) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v CycleQuery) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v CycleQuery) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v CycleQuery) Decide(step string, used int, at time.Time) CycleQueryDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11}...)
	return CycleQueryDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v CycleQuery) Snapshot(at time.Time) CycleQuerySnapshot {
	return CycleQuerySnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 11}
}
