package tracesimulationreceiver

import (
	"context"
	"fmt"
	"github.com/k4ji/tracesimulationreceiver/internal/config"
	"github.com/k4ji/tracesimulationreceiver/internal/config/blueprint"
	"github.com/k4ji/tracesimulationreceiver/internal/config/global"
	"github.com/k4ji/tracesimulationreceiver/internal/metadata"
	simulator "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/adapter/opentelemetry"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver"
)

// createDefaultConfig creates the default configuration for the Trace Simulation receiver.
func createDefaultConfig() component.Config {
	return &config.Config{
		Global:    global.Default(),
		Blueprint: blueprint.Default(),
	}
}

func createTracesReceiver(_ context.Context, params receiver.Settings, baseCfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	logger := params.Logger
	cfg := baseCfg.(*config.Config)
	bp, err := cfg.Blueprint.To()
	if err != nil {
		return nil, fmt.Errorf("failed to convert blueprint: %w", err)
	}

	rcvr := traceSimReceiver{
		logger:       logger,
		nextConsumer: consumer,
		simulator:    simulator.New[[]ptrace.Traces](opentelemetry.NewAdapter()),
		interval:     cfg.Global.Interval,
		blueprint:    bp,
	}

	return &rcvr, nil
}

// NewFactory creates a factory for tracesimulation receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, metadata.TracesStability))
}
