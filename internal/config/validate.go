package config

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Validator struct {
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

type ValidatorDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type ValidatorSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewValidate(name string, limit int) Validator {
	if limit < 1 {
		limit = 1
	}
	return Validator{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "validation"}, Sequence: []string{"inspect", "check", "report", "accept"}, Stage1: "inspect-1", Stage2: "check-2", Stage3: "report-3", Stage4: "accept-4", Stage5: "inspect-5", Stage6: "check-6", Stage7: "report-7", Stage8: "accept-8", Stage9: "inspect-9", Stage10: "check-10", Stage11: "report-11", Stage12: "accept-12", Stage13: "inspect-13", Stage14: "check-14", Stage15: "report-15", Stage16: "accept-16", Stage17: "inspect-17", Stage18: "check-18"}
}

func (v Validator) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17, v.Stage18}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("validation stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("validation name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("validation bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("validation sequence is incomplete")
	}
	return nil
}

func (v Validator) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+18)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17, v.Stage18}...)
	sort.Strings(out)
	return out
}

func (v Validator) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Validator) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17, v.Stage18}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Validator) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17, v.Stage18}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Validator) Decide(step string, used int, at time.Time) ValidatorDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17, v.Stage18}...)
	return ValidatorDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Validator) Snapshot(at time.Time) ValidatorSnapshot {
	return ValidatorSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 18}
}
