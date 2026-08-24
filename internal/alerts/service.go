package alerts

import (
	"example.com/field-sensor-orchestrator/internal/audit"
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
	"time"
)

type Service struct {
	store *store.Store
	audit *audit.Service
	now   func() time.Time
}

func New(st *store.Store, a *audit.Service) *Service {
	return &Service{store: st, audit: a, now: time.Now}
}
func (s *Service) Open(site, sensor, code, message, actor string) domain.Alert {
	seq := s.store.NextSequence()
	alert := domain.Alert{ID: fmt.Sprintf("alert-%06d", seq), SiteID: site, SensorID: sensor, Code: code, Message: message, State: domain.AlertOpen, CreatedAt: s.now().UTC()}
	s.store.PutAlert(alert)
	s.audit.Record(site, alert.ID, "alert.opened", actor, code)
	return alert
}
func (s *Service) Acknowledge(id, actor string) (domain.Alert, error) {
	alert, err := s.store.Alert(id)
	if err != nil {
		return domain.Alert{}, err
	}
	if alert.State == domain.AlertAcknowledged {
		return alert, nil
	}
	alert.State = domain.AlertAcknowledged
	alert.AcknowledgedAt = s.now().UTC()
	alert.AcknowledgedBy = actor
	s.store.PutAlert(alert)
	s.audit.Record(alert.SiteID, alert.ID, "alert.acknowledged", actor, alert.Code)
	return alert, nil
}
func (s *Service) OpenBySite(site string) []domain.Alert {
	out := []domain.Alert{}
	for _, a := range s.store.AlertsBySite(site) {
		if a.State == domain.AlertOpen {
			out = append(out, a)
		}
	}
	return out
}
