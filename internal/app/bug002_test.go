package app

import (
 "context"
 "errors"
 "testing"
 "time"
 "example.com/field-sensor-orchestrator/internal/domain"
)

func TestCancelledRunCycleStopsCollection(t *testing.T) {
 application := NewDemo()
 cycle, err := application.Schedule("coordinator")
 if err != nil { t.Fatal(err) }
 ctx, cancel := context.WithCancel(context.Background())
 cancel()
 _, err = application.Run(ctx, "coordinator")
 if err == nil || !errors.Is(err, context.Canceled) { t.Fatalf("run error=%v want context canceled", err) }
 stored, err := application.Store.Cycle(cycle.ID); if err != nil { t.Fatal(err) }
 if stored.State != domain.CycleFailed { t.Fatalf("state=%s want failed", stored.State) }
 if got := len(application.Store.ReadingsByCycle(cycle.ID)); got != 0 { t.Fatalf("readings=%d want=0", got) }
 lease, err := application.Store.Lease(stored.LeaseID); if err != nil { t.Fatal(err) }
 if lease.ReleasedAt.IsZero() { t.Fatal("cancelled cycle lease was not released") }
}

func TestNormalRunCycleStillCompletes(t *testing.T) {
 application := NewDemo(); if _, err := application.Schedule("coordinator"); err != nil { t.Fatal(err) }
 cycle, err := application.Run(context.Background(), "coordinator"); if err != nil { t.Fatal(err) }
 if cycle.State != domain.CycleCompleted { t.Fatalf("state=%s", cycle.State) }
 if len(application.Store.ReadingsByCycle(cycle.ID)) == 0 { t.Fatal("normal run produced no readings") }
 _ = time.Second
}
