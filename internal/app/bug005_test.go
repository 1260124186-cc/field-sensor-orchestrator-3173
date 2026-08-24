package app

import (
 "context"
 "errors"
 "testing"
 "time"
 "example.com/field-sensor-orchestrator/internal/device"
 "example.com/field-sensor-orchestrator/internal/domain"
)

func TestFailedCycleReleasesLeaseForNextSchedule(t *testing.T) {
 application := NewDemo()
 application.Reader.Configure("sensor-wind", device.Profile{Latency: time.Millisecond, Metric: "wind_speed", Unit: "mps", Fail: true})
 first, err := application.Schedule("coordinator"); if err != nil { t.Fatal(err) }
 if _, err = application.Run(context.Background(), "coordinator"); err == nil || !errors.Is(err, device.ErrUnavailable) { t.Fatalf("run error=%v", err) }
 failed, err := application.Store.Cycle(first.ID); if err != nil { t.Fatal(err) }
 if failed.State != domain.CycleFailed { t.Fatalf("state=%s", failed.State) }
 lease, err := application.Store.Lease(failed.LeaseID); if err != nil { t.Fatal(err) }
 if lease.ReleasedAt.IsZero() { t.Fatal("failed cycle retained its lease") }
 application.Reader.Configure("sensor-wind", device.Profile{Latency: time.Millisecond, Metric: "wind_speed", Unit: "mps", Base: 4.2})
 second, err := application.Schedule("coordinator"); if err != nil { t.Fatalf("next schedule blocked: %v", err) }
 completed, err := application.Run(context.Background(), "coordinator"); if err != nil { t.Fatal(err) }
 if completed.ID != second.ID || completed.State != domain.CycleCompleted { t.Fatalf("next cycle=%+v", completed) }
}

func TestSuccessfulCycleReleasesLease(t *testing.T) {
 application := NewDemo(); cycle, err := application.Schedule("coordinator"); if err != nil { t.Fatal(err) }; if _, err := application.Run(context.Background(), "coordinator"); err != nil { t.Fatal(err) }; lease, err := application.Store.Lease(cycle.LeaseID); if err != nil { t.Fatal(err) }; if lease.ReleasedAt.IsZero() { t.Fatal("successful cycle retained lease") }
}
