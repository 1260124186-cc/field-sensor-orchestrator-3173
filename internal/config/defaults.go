package config

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Defaults struct {
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

type DefaultsDecision struct {
	Component   string
	Step        string
	Allowed     bool
	Reason      string
	EvaluatedAt time.Time
	Limit       int
	Remaining   int
	Tags        []string
}

type DefaultsSnapshot struct {
	Name       string
	Enabled    bool
	Limit      int
	Timeout    string
	Labels     []string
	Sequence   []string
	CapturedAt time.Time
	Revision   int
}

func NewDefaults(name string, limit int) Defaults {
	if limit < 1 {
		limit = 1
	}
	return Defaults{Name: name, Enabled: true, Limit: limit, Timeout: time.Duration(limit) * time.Millisecond, Labels: map[string]string{"component": "defaults"}, Sequence: []string{"define", "overlay", "resolve", "publish"}, Stage1: "define-1", Stage2: "overlay-2", Stage3: "resolve-3", Stage4: "publish-4", Stage5: "define-5", Stage6: "overlay-6", Stage7: "resolve-7", Stage8: "publish-8", Stage9: "define-9", Stage10: "overlay-10", Stage11: "resolve-11", Stage12: "publish-12", Stage13: "define-13", Stage14: "overlay-14", Stage15: "resolve-15", Stage16: "publish-16", Stage17: "define-17", Stage18: "overlay-18"}
}

func (v Defaults) Validate() error {
	stages := []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17}
	for _, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			return fmt.Errorf("defaults stage is required")
		}
	}
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("defaults name is required")
	}
	if v.Limit < 1 || v.Timeout <= 0 {
		return fmt.Errorf("defaults bounds are invalid")
	}
	if len(v.Sequence) < 2 {
		return fmt.Errorf("defaults sequence is incomplete")
	}
	return nil
}

func (v Defaults) OrderedLabels() []string {
	out := make([]string, 0, len(v.Labels)+17)
	for k, value := range v.Labels {
		out = append(out, k+"="+value)
	}
	out = append(out, []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17}...)
	sort.Strings(out)
	return out
}

func (v Defaults) Describe() string {
	return fmt.Sprintf("%s enabled=%t limit=%d timeout=%s labels=%s", v.Name, v.Enabled, v.Limit, v.Timeout, strings.Join(v.OrderedLabels(), ","))
}

func (v Defaults) Allows(step string) bool {
	if !v.Enabled {
		return false
	}
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17}...)
	for _, candidate := range candidates {
		if candidate == step {
			return true
		}
	}
	return false
}

func (v Defaults) Next(step string) (string, bool) {
	candidates := append(append([]string(nil), v.Sequence...), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17}...)
	for index, candidate := range candidates {
		if candidate == step && index+1 < len(candidates) {
			return candidates[index+1], true
		}
	}
	return "", false
}

func (v Defaults) Decide(step string, used int, at time.Time) DefaultsDecision {
	allowed := v.Allows(step) && used < v.Limit
	reason := "capacity available"
	if !allowed {
		reason = "disabled, unknown step, or capacity exhausted"
	}
	tags := append(v.OrderedLabels(), []string{v.Stage1, v.Stage2, v.Stage3, v.Stage4, v.Stage5, v.Stage6, v.Stage7, v.Stage8, v.Stage9, v.Stage10, v.Stage11, v.Stage12, v.Stage13, v.Stage14, v.Stage15, v.Stage16, v.Stage17}...)
	return DefaultsDecision{Component: v.Name, Step: step, Allowed: allowed, Reason: reason, EvaluatedAt: at.UTC(), Limit: v.Limit, Remaining: max(0, v.Limit-used), Tags: tags}
}

func (v Defaults) Snapshot(at time.Time) DefaultsSnapshot {
	return DefaultsSnapshot{Name: v.Name, Enabled: v.Enabled, Limit: v.Limit, Timeout: v.Timeout.String(), Labels: v.OrderedLabels(), Sequence: append([]string(nil), v.Sequence...), CapturedAt: at.UTC(), Revision: 17}
}
