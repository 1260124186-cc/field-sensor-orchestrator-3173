package report

import (
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
)

type Builder struct{ store *store.Store }

func New(st *store.Store) *Builder { return &Builder{store: st} }
func (b *Builder) Build(site string) (domain.SiteSummary, error) {
	if _, err := b.store.Site(site); err != nil {
		return domain.SiteSummary{}, fmt.Errorf("load site: %w", err)
	}
	summary := domain.SiteSummary{SiteID: site}
	for _, sensor := range b.store.SensorsBySite(site) {
		switch sensor.State {
		case domain.SensorActive:
			summary.ActiveSensors++
		case domain.SensorDegraded:
			summary.DegradedSensors++
		}
	}
	for _, cycle := range b.store.CyclesBySite(site) {
		switch cycle.State {
		case domain.CycleCompleted:
			summary.CompletedCycles++
		case domain.CycleFailed:
			summary.FailedCycles++
		}
	}
	for _, alert := range b.store.AlertsBySite(site) {
		if alert.State == domain.AlertOpen {
			summary.OpenAlerts++
		}
	}
	readings := b.store.ReadingsBySite(site)
	summary.ReadingCount = len(readings)
	for _, reading := range readings {
		if reading.ObservedAt.After(summary.LastObservation) {
			summary.LastObservation = reading.ObservedAt
		}
	}
	summary.Health = health(summary)
	return summary, nil
}
func health(s domain.SiteSummary) string {
	if s.DegradedSensors > 0 || s.OpenAlerts > 0 || s.FailedCycles > 0 {
		return "attention"
	}
	if s.ActiveSensors == 0 {
		return "inactive"
	}
	return "healthy"
}
