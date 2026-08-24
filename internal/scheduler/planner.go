package scheduler

import (
	"errors"
	"example.com/field-sensor-orchestrator/internal/audit"
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
	"time"
)

var ErrNoSensors = errors.New("no active sensors")

type Planner struct {
	store *store.Store
	audit *audit.Service
	now   func() time.Time
}

func New(st *store.Store, a *audit.Service) *Planner {
	return &Planner{store: st, audit: a, now: time.Now}
}
func (p *Planner) Schedule(site, actor string, ttl time.Duration) (domain.Cycle, error) {
	if _, err := p.store.Site(site); err != nil {
		return domain.Cycle{}, fmt.Errorf("load site: %w", err)
	}
	ids := []string{}
	for _, sensor := range p.store.SensorsBySite(site) {
		if sensor.State == domain.SensorActive {
			ids = append(ids, sensor.ID)
		}
	}
	if len(ids) == 0 {
		return domain.Cycle{}, ErrNoSensors
	}
	seq := p.store.NextSequence()
	now := p.now().UTC()
	cycle := domain.Cycle{ID: fmt.Sprintf("cycle-%06d", seq), SiteID: site, SensorIDs: ids, State: domain.CycleQueued, CreatedAt: now}
	if err := domain.ValidateCycle(cycle); err != nil {
		return domain.Cycle{}, err
	}
	lease := domain.Lease{ID: fmt.Sprintf("lease-%06d", seq), CycleID: cycle.ID, Owner: actor, AcquiredAt: now, ExpiresAt: now.Add(ttl)}
	if err := domain.ValidateLease(lease, now); err != nil {
		return domain.Cycle{}, err
	}
	cycle.State = domain.CycleLeased
	cycle.LeaseID = lease.ID
	p.store.PutCycle(cycle)
	p.store.PutLease(lease)
	p.audit.Record(site, cycle.ID, "cycle.scheduled", actor, fmt.Sprintf("sensors=%d ttl=%s", len(ids), ttl))
	return cycle, nil
}
