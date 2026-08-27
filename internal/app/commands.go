package app

import (
	"context"
	"fmt"
)

func Execute(ctx context.Context, action string) error {
	application := NewDemo()
	switch action {
	case "demo":
		if err := application.Demo(ctx); err != nil {
			return err
		}
		fmt.Println("demo completed")
	case "register-sensor":
		sensor, err := application.RegisterSensor("field-tech")
		if err != nil {
			return err
		}
		fmt.Printf("registered sensor %s\n", sensor.ID)
	case "schedule-cycle":
		cycle, err := application.Schedule("coordinator")
		if err != nil {
			return err
		}
		fmt.Printf("scheduled cycle %s\n", cycle.ID)
	case "run-cycle":
		if _, err := application.Schedule("coordinator"); err != nil {
			return err
		}
		cycle, err := application.Run(ctx, "coordinator")
		if err != nil {
			return err
		}
		fmt.Printf("completed cycle %s\n", cycle.ID)
	case "calibrate-sensor":
		sensor, err := application.Calibrate("field-tech")
		if err != nil {
			return err
		}
		fmt.Printf("calibrated sensor %s\n", sensor.ID)
	case "ack-alert":
		alert, err := application.Acknowledge("coordinator")
		if err != nil {
			return err
		}
		fmt.Printf("acknowledged alert %s\n", alert.ID)
	case "export-summary":
		summary, err := application.Summary()
		if err != nil {
			return err
		}
		fmt.Printf("summary exported %s\n", summary)
	default:
		return fmt.Errorf("unknown action %q", action)
	}
	return nil
}
