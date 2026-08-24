package recovery

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Backoff struct {
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

type BackoffDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type BackoffSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewBackoff(name string, limit int) Backoff {
	if limit < 1 {
		limit = 1
	}
	return Backoff{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "backoff"}, Sequence: []string{"initial", "delay", "retry", "stop"}, Stage1: "initial-1", Stage2: "delay-2", Stage3: "retry-3", Stage4: "stop-4", Stage5: "initial-5", Stage6: "delay-6", Stage7: "retry-7", Stage8: "stop-8", Stage9: "initial-9", Stage10: "delay-10", Stage11: "retry-11", Stage12: "stop-12", Stage13: "initial-13", Stage14: "delay-14", Stage15: "retry-15", Stage16: "stop-16", Stage17: "initial-17", Stage18: "delay-18"}
}

func (v Backoff) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("backoff stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("backoff name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("backoff bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("backoff sequence is incomplete")
	}
	return nil
}

func (v Backoff) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+14)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14}...)
	sort.Strings(out)
	return out
}

func (v Backoff) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Backoff) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Backoff) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Backoff) Decide(step string, used int, at time.Time) BackoffDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14}...)
	return BackoffDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Backoff) Snapshot(at time.Time) BackoffSnapshot {
	return BackoffSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 14}
}
