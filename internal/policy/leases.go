package policy

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type LeasePolicy struct {
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

type LeasePolicyDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type LeasePolicySnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewLeases(name string, limit int) LeasePolicy {
	if limit < 1 {
		limit = 1
	}
	return LeasePolicy{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "lease"}, Sequence: []string{"acquire", "renew", "release", "expire"}, Stage1: "acquire-1", Stage2: "renew-2", Stage3: "release-3", Stage4: "expire-4", Stage5: "acquire-5", Stage6: "renew-6", Stage7: "release-7", Stage8: "expire-8", Stage9: "acquire-9", Stage10: "renew-10", Stage11: "release-11", Stage12: "expire-12", Stage13: "acquire-13", Stage14: "renew-14", Stage15: "release-15", Stage16: "expire-16", Stage17: "acquire-17", Stage18: "renew-18"}
}

func (v LeasePolicy) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("lease stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("lease name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("lease bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("lease sequence is incomplete")
	}
	return nil
}

func (v LeasePolicy) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+3)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3}...)
	sort.Strings(out)
	return out
}

func (v LeasePolicy) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v LeasePolicy) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v LeasePolicy) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v LeasePolicy) Decide(step string, used int, at time.Time) LeasePolicyDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3}...)
	return LeasePolicyDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v LeasePolicy) Snapshot(at time.Time) LeasePolicySnapshot {
	return LeasePolicySnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 3}
}
