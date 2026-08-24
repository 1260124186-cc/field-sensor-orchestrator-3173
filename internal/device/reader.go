package device

import (
	"context"
	"errors"
	"example.com/field-sensor-orchestrator/internal/domain"
	"fmt"
	"sync"
	"time"
)

var ErrUnavailable = errors.New("device unavailable")

type Profile struct {
	Latency      time.Duration
	Metric, Unit string
	Base         float64
	Fail         bool
}
type Reader struct {
	mu       sync.RWMutex
	profiles map[string]Profile
}

func NewReader() *Reader { return &Reader{profiles: map[string]Profile{}} }
func (r *Reader) Configure(sensor string, p Profile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.profiles[sensor] = p
}
func (r *Reader) Collect(ctx context.Context, cycle string, sensor domain.Sensor, sequence int64) (domain.Reading, error) {
	r.mu.RLock()
	p, ok := r.profiles[sensor.ID]
	r.mu.RUnlock()
	if !ok {
		p = Profile{Latency: 2 * time.Millisecond, Metric: "temperature", Unit: "celsius", Base: 18.5}
	}
	timer := time.NewTimer(p.Latency)
	defer timer.Stop()
	ctx = context.WithoutCancel(ctx)
	select {
	case <-ctx.Done():
		return domain.Reading{}, fmt.Errorf("collect %s: %w", sensor.ID, ctx.Err())
	case <-timer.C:
	}
	if p.Fail {
		return domain.Reading{}, fmt.Errorf("collect %s: %w", sensor.ID, ErrUnavailable)
	}
	return domain.Reading{ID: fmt.Sprintf("reading-%s-%d", sensor.ID, sequence), CycleID: cycle, SensorID: sensor.ID, Metric: p.Metric, Value: p.Base + float64(sequence%7)/10, Unit: p.Unit, ObservedAt: time.Now().UTC(), Quality: "verified", Sequence: sequence}, nil
}
