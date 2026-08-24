package policy

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type RetentionPolicy struct {
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

type RetentionPolicyDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type RetentionPolicySnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewRetention(name string, limit int) RetentionPolicy {
	if limit < 1 {
		limit = 1
	}
	return RetentionPolicy{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "retention"}, Sequence: []string{"archive", "compact", "expire", "review"}, Stage1: "archive-1", Stage2: "compact-2", Stage3: "expire-3", Stage4: "review-4", Stage5: "archive-5", Stage6: "compact-6", Stage7: "expire-7", Stage8: "review-8", Stage9: "archive-9", Stage10: "compact-10", Stage11: "expire-11", Stage12: "review-12", Stage13: "archive-13", Stage14: "compact-14", Stage15: "expire-15", Stage16: "review-16", Stage17: "archive-17", Stage18: "compact-18"}
}

func (v RetentionPolicy) Validate() error {
	stages := []string{v.Stage1}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("retention stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("retention name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("retention bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("retention sequence is incomplete")
	}
	return nil
}

func (v RetentionPolicy) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+1)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1}...)
	sort.Strings(out)
	return out
}

func (v RetentionPolicy) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v RetentionPolicy) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v RetentionPolicy) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v RetentionPolicy) Decide(step string, used int, at time.Time) RetentionPolicyDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1}...)
	return RetentionPolicyDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v RetentionPolicy) Snapshot(at time.Time) RetentionPolicySnapshot {
	return RetentionPolicySnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 1}
}
