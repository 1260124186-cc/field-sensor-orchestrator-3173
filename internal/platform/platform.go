package platform

import (
	"example.com/field-sensor-orchestrator/internal/config"
	"example.com/field-sensor-orchestrator/internal/lifecycle"
	"example.com/field-sensor-orchestrator/internal/policy"
	"example.com/field-sensor-orchestrator/internal/query"
	"example.com/field-sensor-orchestrator/internal/recovery"
	"example.com/field-sensor-orchestrator/internal/telemetry"
	"fmt"
)

func Describe() []string {
	items := []interface {
		Validate() error
		Describe() string
	}{config.NewConfig("runtime", 4), config.NewDefaults("defaults", 3), config.NewValidate("validation", 2), lifecycle.NewSensors("sensors", 8), lifecycle.NewCycles("cycles", 8), lifecycle.NewLeases("leases", 8), policy.NewRetention("retention", 7), policy.NewQuality("quality", 6), policy.NewLeases("lease-policy", 5), query.NewSensors("sensor-query", 5), query.NewCycles("cycle-query", 5), query.NewEvents("event-query", 5), recovery.NewPlan("plan", 4), recovery.NewBackoff("backoff", 4), recovery.NewClassify("classify", 4), telemetry.NewWindow("window", 10), telemetry.NewStats("stats", 10), telemetry.NewFilter("filter", 10)}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if err := item.Validate(); err != nil {
			out = append(out, fmt.Sprintf("invalid=%v", err))
			continue
		}
		out = append(out, item.Describe())
	}
	return out
}
