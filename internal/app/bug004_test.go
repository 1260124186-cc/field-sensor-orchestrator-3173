package app

import (
 "context"
 "strings"
 "testing"
 "time"
 "example.com/field-sensor-orchestrator/internal/domain"
)

func TestSiteSummaryDoesNotLeakOtherSiteState(t *testing.T) {
 application := NewDemo()
 other := "valley-south"
 application.Store.PutSite(domain.Site{ID: other, Name: "South Valley", Region: "wetland", Enabled: true})
 application.Store.PutSensor(domain.Sensor{ID: "sensor-other", SiteID: other, Kind: "water", Label: "valley water", State: domain.SensorDegraded})
 cycle := domain.Cycle{ID: "cycle-other", SiteID: other, SensorIDs: []string{"sensor-other"}, State: domain.CycleCompleted, CreatedAt: time.Now(), FinishedAt: time.Now()}
 application.Store.PutCycle(cycle)
 application.Store.PutReading(domain.Reading{ID: "reading-other", CycleID: cycle.ID, SensorID: "sensor-other", Metric: "water", Value: 1, Unit: "m", ObservedAt: time.Now(), Quality: "verified", Sequence: 900})
 application.Alerts.Open(other, "sensor-other", "flood", "water high", "valley-tech")
 summary, err := application.Summary(); if err != nil { t.Fatal(err) }
 for _, forbidden := range []string{"degraded=1", "completed=1", "alerts=1", "readings=1"} { if strings.Contains(summary, forbidden) { t.Fatalf("ridge summary leaked other-site state: %s", summary) } }
 if !strings.Contains(summary, "site=ridge-north") || !strings.Contains(summary, "active=3") { t.Fatalf("unexpected ridge summary: %s", summary) }
 _ = context.Background()
}

func TestSummaryRepeatedCallIsStable(t *testing.T) {
 application := NewDemo(); first, err := application.Summary(); if err != nil { t.Fatal(err) }; second, err := application.Summary(); if err != nil { t.Fatal(err) }; if first != second { t.Fatalf("summary changed across reads: %s / %s", first, second) }
}
