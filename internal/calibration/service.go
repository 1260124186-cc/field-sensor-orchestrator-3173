package calibration

import (
	"errors"
	"example.com/field-sensor-orchestrator/internal/audit"
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
	"math"
	"time"
)

var ErrInvalidCalibration = errors.New("invalid calibration")

type Service struct {
	store *store.Store
	audit *audit.Service
	now   func() time.Time
}

func New(st *store.Store, a *audit.Service) *Service {
	return &Service{store: st, audit: a, now: time.Now}
}
func (s *Service) Apply(sensorID, operator string, offset, scale float64) (domain.Sensor, error) {
	sensor, err := s.store.Sensor(sensorID)
	if err != nil {
		return domain.Sensor{}, err
	}
	if math.IsNaN(offset) || math.IsNaN(scale) || scale <= 0 {
		return domain.Sensor{}, ErrInvalidCalibration
	}
	if err := domain.TransitionSensor(sensor.State, domain.SensorCalibrating); err != nil {
		return domain.Sensor{}, err
	}
	sensor.State = domain.SensorCalibrating
	s.store.PutSensor(sensor)
	record := domain.Calibration{ID: fmt.Sprintf("cal-%d", s.store.NextSequence()), SensorID: sensorID, Operator: operator, Offset: offset, Scale: scale, AppliedAt: s.now().UTC(), Notes: "field reference applied"}
	s.store.PutCalibration(record)
	if err := domain.TransitionSensor(sensor.State, domain.SensorActive); err != nil {
		return domain.Sensor{}, err
	}
	sensor.State = domain.SensorActive
	sensor.LastSeen = s.now().UTC()
	s.store.PutSensor(sensor)
	s.audit.Record(sensor.SiteID, sensor.ID, "sensor.calibrated", operator, fmt.Sprintf("offset=%.2f scale=%.2f", offset, scale))
	return sensor, nil
}
