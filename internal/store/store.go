package store

import (
	"errors"
	"example.com/field-sensor-orchestrator/internal/domain"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("record not found")

type Store struct {
	mu           sync.RWMutex
	sites        map[string]domain.Site
	sensors      map[string]domain.Sensor
	cycles       map[string]domain.Cycle
	leases       map[string]domain.Lease
	readings     map[string]domain.Reading
	calibrations map[string]domain.Calibration
	alerts       map[string]domain.Alert
	events       []domain.AuditEvent
	sequence     int64
}

func New() *Store {
	return &Store{sites: map[string]domain.Site{}, sensors: map[string]domain.Sensor{}, cycles: map[string]domain.Cycle{}, leases: map[string]domain.Lease{}, readings: map[string]domain.Reading{}, calibrations: map[string]domain.Calibration{}, alerts: map[string]domain.Alert{}}
}
func (s *Store) NextSequence() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sequence++
	return s.sequence
}
func (s *Store) PutSite(v domain.Site) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sites[v.ID] = cloneSite(v)
}
func (s *Store) Site(id string) (domain.Site, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sites[id]
	if !ok {
		return domain.Site{}, ErrNotFound
	}
	return cloneSite(v), nil
}
func (s *Store) PutSensor(v domain.Sensor) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.Revision++
	s.sensors[v.ID] = cloneSensor(v)
}
func (s *Store) Sensor(id string) (domain.Sensor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.sensors[id]
	if !ok {
		return domain.Sensor{}, ErrNotFound
	}
	return cloneSensor(v), nil
}
func (s *Store) SensorsBySite(site string) []domain.Sensor {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Sensor{}
	for _, v := range s.sensors {
		if v.SiteID == site {
			out = append(out, cloneSensor(v))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (s *Store) PutCycle(v domain.Cycle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.Revision++
	s.cycles[v.ID] = cloneCycle(v)
}
func (s *Store) Cycle(id string) (domain.Cycle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.cycles[id]
	if !ok {
		return domain.Cycle{}, ErrNotFound
	}
	return cloneCycle(v), nil
}
func (s *Store) CyclesBySite(site string) []domain.Cycle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Cycle{}
	for _, v := range s.cycles {
		if v.SiteID == site {
			out = append(out, cloneCycle(v))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func (s *Store) PutLease(v domain.Lease) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.Revision++
	s.leases[v.ID] = v
}
func (s *Store) Lease(id string) (domain.Lease, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.leases[id]
	if !ok {
		return domain.Lease{}, ErrNotFound
	}
	return v, nil
}
func (s *Store) ReleaseLease(id string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.leases[id]
	if !ok {
		return ErrNotFound
	}
	v.ReleasedAt = at
	v.Revision++
	s.leases[id] = v
	return nil
}
func (s *Store) PutReading(v domain.Reading) { s.mu.Lock(); defer s.mu.Unlock(); s.readings[v.ID] = v }
func (s *Store) ReadingsByCycle(id string) []domain.Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Reading{}
	for _, v := range s.readings {
		if v.CycleID == id {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Sequence < out[j].Sequence })
	return out
}
func (s *Store) ReadingsBySite(site string) []domain.Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cycleIDs := map[string]bool{}
	for _, c := range s.cycles {
		if c.SiteID == site {
			cycleIDs[c.ID] = true
		}
	}
	out := []domain.Reading{}
	for _, v := range s.readings {
		if cycleIDs[v.CycleID] {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ObservedAt.Before(out[j].ObservedAt) })
	return out
}
