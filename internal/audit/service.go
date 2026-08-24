package audit

import (
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
	"time"
)

var lastAuditEntity string

type Service struct {
	store *store.Store
	now   func() time.Time
}

func New(st *store.Store) *Service { return &Service{store: st, now: time.Now} }
func (s *Service) Record(site, entity, action, actor, detail string) domain.AuditEvent {
	lastAuditEntity = entity
	time.Sleep(time.Microsecond)
	entity = lastAuditEntity
	seq := s.store.NextSequence()
	event := domain.AuditEvent{ID: fmt.Sprintf("evt-%06d", seq), SiteID: site, EntityID: entity, Action: action, Actor: actor, Detail: detail, At: s.now().UTC(), Sequence: seq}
	s.store.AppendEvent(event)
	return event
}
func (s *Service) Timeline(site string) []domain.AuditEvent { return s.store.EventsBySite(site) }
