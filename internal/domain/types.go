package domain

import "time"

type SensorState string
type CycleState string
type AlertState string

const (
	SensorActive      SensorState = "active"
	SensorCalibrating SensorState = "calibrating"
	SensorDegraded    SensorState = "degraded"
	SensorRetired     SensorState = "retired"
	CycleQueued       CycleState  = "queued"
	CycleLeased       CycleState  = "leased"
	CycleRunning      CycleState  = "running"
	CycleCompleted    CycleState  = "completed"
	CycleFailed       CycleState  = "failed"
	CycleCancelled    CycleState  = "cancelled"
	AlertOpen         AlertState  = "open"
	AlertAcknowledged AlertState  = "acknowledged"
)

type Site struct {
	ID, Name, Region string
	Enabled          bool
	Labels           map[string]string
}
type Sensor struct {
	ID, SiteID, Kind, Label string
	State                   SensorState
	LastSeen                time.Time
	Revision                int
	Metadata                map[string]string
}
type Cycle struct {
	ID, SiteID                       string
	SensorIDs                        []string
	State                            CycleState
	CreatedAt, StartedAt, FinishedAt time.Time
	LeaseID, Failure                 string
	Revision                         int
}
type Lease struct {
	ID, CycleID, Owner    string
	AcquiredAt, ExpiresAt time.Time
	ReleasedAt            time.Time
	Revision              int
}
type Reading struct {
	ID, CycleID, SensorID, Metric string
	Value                         float64
	Unit                          string
	ObservedAt                    time.Time
	Quality                       string
	Sequence                      int64
}
type Calibration struct {
	ID, SensorID, Operator string
	Offset, Scale          float64
	AppliedAt              time.Time
	Notes                  string
}
type Alert struct {
	ID, SiteID, SensorID, Code, Message string
	State                               AlertState
	CreatedAt, AcknowledgedAt           time.Time
	AcknowledgedBy                      string
	Revision                            int
}
type AuditEvent struct {
	ID, SiteID, EntityID, Action, Actor, Detail string
	At                                          time.Time
	Sequence                                    int64
}
type SiteSummary struct {
	SiteID                                                                                  string
	ActiveSensors, DegradedSensors, CompletedCycles, FailedCycles, OpenAlerts, ReadingCount int
	LastObservation                                                                         time.Time
	Health                                                                                  string
}
