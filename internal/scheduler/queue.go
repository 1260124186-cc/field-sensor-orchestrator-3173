package scheduler

import (
	"errors"
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/store"
	"sort"
	"time"
)

var ErrLeaseExpired = errors.New("cycle lease expired")

func NextRunnable(st *store.Store, site string, now time.Time) (domain.Cycle, error) {
	cycles := st.CyclesBySite(site)
	sort.SliceStable(cycles, func(i, j int) bool { return cycles[i].CreatedAt.Before(cycles[j].CreatedAt) })
	for _, cycle := range cycles {
		if cycle.State != domain.CycleLeased {
			continue
		}
		lease, err := st.Lease(cycle.LeaseID)
		if err != nil {
			return domain.Cycle{}, err
		}
		if !lease.ReleasedAt.IsZero() {
			continue
		}
		if !lease.ExpiresAt.After(now) {
			return domain.Cycle{}, ErrLeaseExpired
		}
		return cycle, nil
	}
	return domain.Cycle{}, store.ErrNotFound
}
