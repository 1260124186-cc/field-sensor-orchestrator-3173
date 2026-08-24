package main

import (
	"context"
	"example.com/field-sensor-orchestrator/internal/app"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	action := "demo"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}
	if err := app.Execute(ctx, action); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
