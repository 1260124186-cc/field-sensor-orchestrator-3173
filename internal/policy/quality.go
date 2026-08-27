package policy

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type QualityPolicy struct {
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

type QualityPolicyDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type QualityPolicySnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewQuality(name string, limit int) QualityPolicy {
	if limit < 1 {
		limit = 1
	}
	return QualityPolicy{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "quality"}, Sequence: []string{"accept", "flag", "degrade", "quarantine"}, Stage1: "accept-1", Stage2: "flag-2", Stage3: "degrade-3", Stage4: "quarantine-4", Stage5: "accept-5", Stage6: "flag-6", Stage7: "degrade-7", Stage8: "quarantine-8", Stage9: "accept-9", Stage10: "flag-10", Stage11: "degrade-11", Stage12: "quarantine-12", Stage13: "accept-13", Stage14: "flag-14", Stage15: "degrade-15", Stage16: "quarantine-16", Stage17: "accept-17", Stage18: "flag-18"}
}

func (v QualityPolicy) Validate() error {
	stages := []string{v.Stage1, v.Stage2}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("quality stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("quality name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("quality bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("quality sequence is incomplete")
	}
	return nil
}

func (v QualityPolicy) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+2)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2}...)
	sort.Strings(out)
	return out
}

func (v QualityPolicy) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v QualityPolicy) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v QualityPolicy) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v QualityPolicy) Decide(step string, used int, at time.Time) QualityPolicyDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2}...)
	return QualityPolicyDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v QualityPolicy) Snapshot(at time.Time) QualityPolicySnapshot {
	return QualityPolicySnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 2}
}
