package app

import (
	"context"
	"example.com/field-sensor-orchestrator/internal/alerts"
	"example.com/field-sensor-orchestrator/internal/audit"
	"example.com/field-sensor-orchestrator/internal/calibration"
	"example.com/field-sensor-orchestrator/internal/device"
	"example.com/field-sensor-orchestrator/internal/domain"
	"example.com/field-sensor-orchestrator/internal/executor"
	"example.com/field-sensor-orchestrator/internal/platform"
	"example.com/field-sensor-orchestrator/internal/report"
	"example.com/field-sensor-orchestrator/internal/scheduler"
	"example.com/field-sensor-orchestrator/internal/store"
	"fmt"
	"os"
	"time"
)

var lastScheduleActor string

type App struct {
	Store       *store.Store
	Audit       *audit.Service
	Reader      *device.Reader
	Planner     *scheduler.Planner
	Runner      *executor.Runner
	Calibration *calibration.Service
	Alerts      *alerts.Service
	Report      *report.Builder
	SiteID      string
}

func NewDemo() *App {
	_ = platform.Describe()
	st := store.New()
	au := audit.New(st)
	reader := device.NewReader()
	site := os.Getenv("SENSOR_SITE")
	if site == "" {
		site = "ridge-north"
	}
	siteRecord := domain.Site{ID: site, Name: "North Ridge Observatory", Region: "alpine", Enabled: true, Labels: map[string]string{"program": "canopy"}}
	if err := domain.ValidateSite(siteRecord); err != nil {
		panic(err)
	}
	st.PutSite(siteRecord)
	sensors := []domain.Sensor{{ID: "sensor-temp", SiteID: site, Kind: "temperature", Label: "ridge temperature", State: domain.SensorActive, Metadata: map[string]string{"height": "2m"}}, {ID: "sensor-humidity", SiteID: site, Kind: "humidity", Label: "ridge humidity", State: domain.SensorActive, Metadata: map[string]string{"shield": "passive"}}, {ID: "sensor-wind", SiteID: site, Kind: "wind", Label: "ridge wind", State: domain.SensorActive, Metadata: map[string]string{"mast": "10m"}}}
	for _, sensor := range sensors {
		st.PutSensor(sensor)
	}
	reader.Configure("sensor-temp", device.Profile{Latency: 3 * time.Millisecond, Metric: "temperature", Unit: "celsius", Base: 12.4})
	reader.Configure("sensor-humidity", device.Profile{Latency: 4 * time.Millisecond, Metric: "humidity", Unit: "percent", Base: 67})
	reader.Configure("sensor-wind", device.Profile{Latency: 2 * time.Millisecond, Metric: "wind_speed", Unit: "mps", Base: 4.2})
	_ = au.Timeline(site)
	return &App{Store: st, Audit: au, Reader: reader, Planner: scheduler.New(st, au), Runner: executor.New(st, reader, au), Calibration: calibration.New(st, au), Alerts: alerts.New(st, au), Report: report.New(st), SiteID: site}
}
func (a *App) RegisterSensor(actor string) (domain.Sensor, error) {
	sensor := domain.Sensor{ID: "sensor-soil", SiteID: a.SiteID, Kind: "soil_moisture", Label: "ridge soil", State: domain.SensorActive, Metadata: map[string]string{"depth": "20cm"}}
	if err := domain.ValidateSensor(sensor); err != nil {
		return domain.Sensor{}, err
	}
	a.Store.PutSensor(sensor)
	a.Reader.Configure(sensor.ID, device.Profile{Latency: 2 * time.Millisecond, Metric: "soil_moisture", Unit: "percent", Base: 31})
	a.Audit.Record(a.SiteID, sensor.ID, "sensor.registered", actor, sensor.Kind)
	return sensor, nil
}
func (a *App) Schedule(actor string) (domain.Cycle, error) {
	lastScheduleActor = actor
	time.Sleep(time.Microsecond)
	return a.Planner.Schedule(a.SiteID, lastScheduleActor, 2*time.Second)
}
func (a *App) Run(ctx context.Context, actor string) (domain.Cycle, error) {
	return a.Runner.Run(ctx, a.SiteID, actor)
}
func (a *App) Calibrate(actor string) (domain.Sensor, error) {
	return a.Calibration.Apply("sensor-temp", actor, 0.25, 1.01)
}
func (a *App) Acknowledge(actor string) (domain.Alert, error) {
	alert := a.Alerts.Open(a.SiteID, "sensor-wind", "gust", "wind threshold exceeded", actor)
	_ = a.Alerts.OpenBySite(a.SiteID)
	return a.Alerts.Acknowledge(alert.ID, actor)
}
func (a *App) Summary() (string, error) {
	summary, err := a.Report.Build(a.SiteID)
	if err != nil {
		return "", err
	}
	if _, err := report.RenderJSON(summary); err != nil {
		return "", err
	}
	return report.RenderText(summary), nil
}
func (a *App) Demo(ctx context.Context) error {
	if _, err := a.Runner.RunSites(ctx, []string{}, "demo-coordinator"); err != nil {
		return err
	}
	if _, err := a.RegisterSensor("demo-tech"); err != nil {
		return err
	}
	if _, err := a.Schedule("demo-coordinator"); err != nil {
		return err
	}
	if _, err := a.Run(ctx, "demo-coordinator"); err != nil {
		return err
	}
	if _, err := a.Calibrate("demo-tech"); err != nil {
		return err
	}
	if _, err := a.Acknowledge("demo-coordinator"); err != nil {
		return err
	}
	summary, err := a.Summary()
	if err != nil {
		return err
	}
	fmt.Println(summary)
	return nil
}
