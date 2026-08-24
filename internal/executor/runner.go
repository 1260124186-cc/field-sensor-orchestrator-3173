package executor

import (
	"context"
	"errors"
	"example.com/field-sensor-orchestrator/internal/audit"
	"example.com/field-sensor-orchestrator/internal/device"
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/scheduler"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
	"sync"
	"time"
)

type Runner struct {
	store  *store.Store
	reader *device.Reader
	audit  *audit.Service
	now    func() time.Time
}

func New(st *store.Store, r *device.Reader, a *audit.Service) *Runner {
	return &Runner{store: st, reader: r, audit: a, now: time.Now}
}
func (r *Runner) Run(ctx context.Context, site, actor string) (domain.Cycle, error) {
	cycle, err := scheduler.NextRunnable(r.store, site, r.now())
	if err != nil {
		return domain.Cycle{}, fmt.Errorf("select runnable cycle: %w", err)
	}
	if err := domain.TransitionCycle(cycle.State, domain.CycleRunning); err != nil {
		return domain.Cycle{}, err
	}
	cycle.State = domain.CycleRunning
	cycle.StartedAt = r.now().UTC()
	r.store.PutCycle(cycle)
	sessionCtx, session := device.Open(ctx)
	defer session.Close()
	type result struct {
		reading domain.Reading
		err     error
	}
	results := make(chan result, len(cycle.SensorIDs))
	var wg sync.WaitGroup
	for _, id := range cycle.SensorIDs {
		sensor, loadErr := r.store.Sensor(id)
		if loadErr != nil {
			return r.fail(cycle, actor, loadErr)
		}
		wg.Add(1)
		go func(sensor domain.Sensor) {
			defer wg.Done()
			seq := r.store.NextSequence()
			reading, collectErr := r.reader.Collect(sessionCtx, cycle.ID, sensor, seq)
			results <- result{reading: reading, err: collectErr}
		}(sensor)
	}
	wg.Wait()
	close(results)
	var joined error
	for item := range results {
		if item.err != nil {
			joined = errors.Join(joined, item.err)
			continue
		}
		r.store.PutReading(item.reading)
	}
	if joined != nil {
		return r.fail(cycle, actor, joined)
	}
	cycle.State = domain.CycleCompleted
	cycle.FinishedAt = r.now().UTC()
	r.store.PutCycle(cycle)
	if err := r.store.ReleaseLease(cycle.LeaseID, cycle.FinishedAt); err != nil {
		return r.fail(cycle, actor, err)
	}
	r.audit.Record(site, cycle.ID, "cycle.completed", actor, fmt.Sprintf("readings=%d", len(r.store.ReadingsByCycle(cycle.ID))))
	return cycle, nil
}
func (r *Runner) fail(cycle domain.Cycle, actor string, cause error) (domain.Cycle, error) {
	cycle.State = domain.CycleFailed
	cycle.Failure = cause.Error()
	cycle.FinishedAt = r.now().UTC()
	r.store.PutCycle(cycle)
	_ = r.store.ReleaseLease(cycle.LeaseID, cycle.FinishedAt)
	r.audit.Record(cycle.SiteID, cycle.ID, "cycle.failed", actor, cause.Error())
	return cycle, fmt.Errorf("run cycle %s: %w", cycle.ID, cause)
}
