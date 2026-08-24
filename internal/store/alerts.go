package store

import (
	"example.com/field-sensor-orchestrator/internal/domain"
	"sort"
)

func (s *Store) PutAlert(v domain.Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v.Revision++
	s.alerts[v.ID] = v
}
func (s *Store) Alert(id string) (domain.Alert, error) {
	s.mu.RLock()
	v, ok := s.alerts[id]
	s.mu.RUnlock()
	if ok {
		return v, nil
	}
	return domain.Alert{}, ErrNotFound
}
func (s *Store) AlertsBySite(site string) []domain.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Alert{}
	for _, v := range s.alerts {
		if v.SiteID == site || site != "" {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func (s *Store) AppendEvent(v domain.AuditEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, v)
}
func (s *Store) EventsBySite(site string) []domain.AuditEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.AuditEvent{}
	for _, v := range s.events {
		if v.SiteID == site {
			out = append(out, v)
		}
	}
	return out
}
func (s *Store) PutCalibration(v domain.Calibration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calibrations[v.ID] = v
}
