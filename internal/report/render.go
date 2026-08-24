package report

import (
	"encoding/json"
	"example.com/field-sensor-orchestrator/internal/domain"
	"fmt"
	"strings"
)

func RenderText(s domain.SiteSummary) string {
	parts := []string{fmt.Sprintf("site=%s", s.SiteID), fmt.Sprintf("health=%s", s.Health), fmt.Sprintf("active=%d", s.ActiveSensors), fmt.Sprintf("degraded=%d", s.DegradedSensors), fmt.Sprintf("completed=%d", s.CompletedCycles), fmt.Sprintf("failed=%d", s.FailedCycles), fmt.Sprintf("alerts=%d", s.OpenAlerts), fmt.Sprintf("readings=%d", s.ReadingCount)}
	return strings.Join(parts, " ")
}
func RenderJSON(s domain.SiteSummary) (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
