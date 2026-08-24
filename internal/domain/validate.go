package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrInvalidTransition = errors.New("invalid state transition")

func ValidateSite(site Site) error {
	if strings.TrimSpace(site.ID) == "" || strings.TrimSpace(site.Name) == "" {
		return errors.New("site id and name are required")
	}
	if !site.Enabled {
		return errors.New("site is disabled")
	}
	return nil
}
func ValidateSensor(sensor Sensor) error {
	if sensor.ID == "" || sensor.SiteID == "" || sensor.Kind == "" {
		return errors.New("sensor identity is incomplete")
	}
	if sensor.State == "" {
		return errors.New("sensor state is required")
	}
	return nil
}
func ValidateCycle(c Cycle) error {
	if c.ID == "" || c.SiteID == "" {
		return errors.New("cycle identity is incomplete")
	}
	if len(c.SensorIDs) == 0 {
		return errors.New("cycle requires sensors")
	}
	return nil
}
func ValidateLease(l Lease, now time.Time) error {
	if l.ID == "" || l.CycleID == "" || l.Owner == "" {
		return errors.New("lease identity is incomplete")
	}
	if !l.ExpiresAt.After(now) {
		return errors.New("lease is already expired")
	}
	return nil
}
func ValidateReading(r Reading) error {
	if r.ID == "" || r.CycleID == "" || r.SensorID == "" || r.Metric == "" {
		return errors.New("reading identity is incomplete")
	}
	if r.Quality == "" {
		return errors.New("reading quality is required")
	}
	return nil
}
func TransitionSensor(current, next SensorState) error {
	allowed := map[SensorState]map[SensorState]bool{SensorActive: {SensorCalibrating: true, SensorDegraded: true, SensorRetired: true}, SensorCalibrating: {SensorActive: true, SensorDegraded: true}, SensorDegraded: {SensorActive: true, SensorRetired: true}, SensorRetired: {}}
	if current == next {
		return nil
	}
	if allowed[current][next] {
		return nil
	}
	return fmt.Errorf("%w: sensor %s -> %s", ErrInvalidTransition, current, next)
}
func TransitionCycle(current, next CycleState) error {
	allowed := map[CycleState]map[CycleState]bool{CycleQueued: {CycleLeased: true, CycleCancelled: true}, CycleLeased: {CycleRunning: true, CycleCancelled: true, CycleFailed: true}, CycleRunning: {CycleCompleted: true, CycleFailed: true, CycleCancelled: true}, CycleCompleted: {}, CycleFailed: {}, CycleCancelled: {}}
	if current == next {
		return nil
	}
	if allowed[current][next] {
		return nil
	}
	return fmt.Errorf("%w: cycle %s -> %s", ErrInvalidTransition, current, next)
}
