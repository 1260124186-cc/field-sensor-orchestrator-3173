package recovery

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Classification struct {
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

type ClassificationDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type ClassificationSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewClassify(name string, limit int) Classification {
	if limit < 1 {
		limit = 1
	}
	return Classification{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "classification"}, Sequence: []string{"inspect", "match", "label", "route"}, Stage1: "inspect-1", Stage2: "match-2", Stage3: "label-3", Stage4: "route-4", Stage5: "inspect-5", Stage6: "match-6", Stage7: "label-7", Stage8: "route-8", Stage9: "inspect-9", Stage10: "match-10", Stage11: "label-11", Stage12: "route-12", Stage13: "inspect-13", Stage14: "match-14", Stage15: "label-15", Stage16: "route-16", Stage17: "inspect-17", Stage18: "match-18"}
}

func (v Classification) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("classification stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("classification name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("classification bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("classification sequence is incomplete")
	}
	return nil
}

func (v Classification) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+15)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15}...)
	sort.Strings(out)
	return out
}

func (v Classification) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Classification) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Classification) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Classification) Decide(step string, used int, at time.Time) ClassificationDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15}...)
	return ClassificationDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Classification) Snapshot(at time.Time) ClassificationSnapshot {
	return ClassificationSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 15}
}
