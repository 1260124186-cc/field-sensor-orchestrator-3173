package store

import "example.com/field-sensor-orchestrator/internal/domain"

func cloneSite(v domain.Site) domain.Site       { v.Labels = cloneMap(v.Labels); return v }
func cloneSensor(v domain.Sensor) domain.Sensor { v.Metadata = cloneMap(v.Metadata); return v }
func cloneCycle(v domain.Cycle) domain.Cycle {
	v.SensorIDs = append([]string(nil), v.SensorIDs...)
	return v
}
func cloneMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
func (s *Store) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := Snapshot{Sites: map[string]domain.Site{}, Sensors: map[string]domain.Sensor{}, Cycles: map[string]domain.Cycle{}, Leases: map[string]domain.Lease{}, Readings: map[string]domain.Reading{}, Alerts: map[string]domain.Alert{}, Events: append([]domain.AuditEvent(nil), s.events...)}
	for k, v := range s.sites {
		out.Sites[k] = cloneSite(v)
	}
	for k, v := range s.sensors {
		out.Sensors[k] = cloneSensor(v)
	}
	for k, v := range s.cycles {
		out.Cycles[k] = cloneCycle(v)
	}
	for k, v := range s.leases {
		out.Leases[k] = v
	}
	for k, v := range s.readings {
		out.Readings[k] = v
	}
	for k, v := range s.alerts {
		out.Alerts[k] = v
	}
	return out
}

type Snapshot struct {
	Sites    map[string]domain.Site
	Sensors  map[string]domain.Sensor
	Cycles   map[string]domain.Cycle
	Leases   map[string]domain.Lease
	Readings map[string]domain.Reading
	Alerts   map[string]domain.Alert
	Events   []domain.AuditEvent
}
