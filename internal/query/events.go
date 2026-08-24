package query

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type EventQuery struct {
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

type EventQueryDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type EventQuerySnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewEvents(name string, limit int) EventQuery {
	if limit < 1 {
		limit = 1
	}
	return EventQuery{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "event query"}, Sequence: []string{"scope", "sequence", "limit", "return"}, Stage1: "scope-1", Stage2: "sequence-2", Stage3: "limit-3", Stage4: "return-4", Stage5: "scope-5", Stage6: "sequence-6", Stage7: "limit-7", Stage8: "return-8", Stage9: "scope-9", Stage10: "sequence-10", Stage11: "limit-11", Stage12: "return-12", Stage13: "scope-13", Stage14: "sequence-14", Stage15: "limit-15", Stage16: "return-16", Stage17: "scope-17", Stage18: "sequence-18"}
}

func (v EventQuery) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("event query stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("event query name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("event query bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("event query sequence is incomplete")
	}
	return nil
}

func (v EventQuery) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+12)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12}...)
	sort.Strings(out)
	return out
}

func (v EventQuery) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v EventQuery) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v EventQuery) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v EventQuery) Decide(step string, used int, at time.Time) EventQueryDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12}...)
	return EventQueryDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v EventQuery) Snapshot(at time.Time) EventQuerySnapshot {
	return EventQuerySnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 12}
}
