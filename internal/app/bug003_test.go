package app

import (
 "context"
 "errors"
 "testing"
 "time"
 "example.com/field-sensor-orchestrator/internal/device"
 "example.com/field-sensor-orchestrator/internal/domain"
)

func TestDeviceFailurePropagatesAndRollsBackCycle(t *testing.T) {
 application := NewDemo()
 application.Reader.Configure("sensor-wind", device.Profile{Latency: time.Millisecond, Metric: "wind_speed", Unit: "mps", Base: 4.2, Fail: true})
 cycle, err := application.Schedule("coordinator"); if err != nil { t.Fatal(err) }
 _, err = application.Run(context.Background(), "coordinator")
 if err == nil || !errors.Is(err, device.ErrUnavailable) { t.Fatalf("run error=%v want device unavailable", err) }
 stored, err := application.Store.Cycle(cycle.ID); if err != nil { t.Fatal(err) }
 if stored.State != domain.CycleFailed { t.Fatalf("state=%s want failed", stored.State) }
 if got := len(application.Store.ReadingsByCycle(cycle.ID)); got != 2 { t.Fatalf("successful readings=%d want=2", got) }
 lease, err := application.Store.Lease(stored.LeaseID); if err != nil { t.Fatal(err) }
 if lease.ReleasedAt.IsZero() { t.Fatal("failed cycle lease not released") }
}

func TestHealthyDevicesCompleteCycle(t *testing.T) {
 application := NewDemo(); if _, err := application.Schedule("coordinator"); err != nil { t.Fatal(err) }
 cycle, err := application.Run(context.Background(), "coordinator"); if err != nil { t.Fatal(err) }
 if cycle.State != domain.CycleCompleted { t.Fatalf("state=%s", cycle.State) }
}
