package telemetry

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Statistics struct {
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

type StatisticsDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type StatisticsSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewStats(name string, limit int) Statistics {
	if limit < 1 {
		limit = 1
	}
	return Statistics{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "statistics"}, Sequence: []string{"collect", "summarize", "publish", "retain"}, Stage1: "collect-1", Stage2: "summarize-2", Stage3: "publish-3", Stage4: "retain-4", Stage5: "collect-5", Stage6: "summarize-6", Stage7: "publish-7", Stage8: "retain-8", Stage9: "collect-9", Stage10: "summarize-10", Stage11: "publish-11", Stage12: "retain-12", Stage13: "collect-13", Stage14: "summarize-14", Stage15: "publish-15", Stage16: "retain-16", Stage17: "collect-17", Stage18: "summarize-18"}
}

func (v Statistics) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("statistics stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("statistics name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("statistics bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("statistics sequence is incomplete")
	}
	return nil
}

func (v Statistics) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+5)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5}...)
	sort.Strings(out)
	return out
}

func (v Statistics) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Statistics) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Statistics) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Statistics) Decide(step string, used int, at time.Time) StatisticsDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5}...)
	return StatisticsDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Statistics) Snapshot(at time.Time) StatisticsSnapshot {
	return StatisticsSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 5}
}
