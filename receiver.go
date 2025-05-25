package tracesimulationreceiver

import (
	"context"
	simulator "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/blueprint"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"time"
)

var _ receiver.Traces = (*traceSimReceiver)(nil)

type traceSimReceiver struct {
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	simulator    *simulator.Simulator[[]ptrace.Traces]
	interval     time.Duration
	blueprint    blueprint.Blueprint
}

func (r *traceSimReceiver) Start(ctx context.Context, _ component.Host) error {
	ctx, r.cancel = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()

		if err := r.emitTracesOnce(ctx); err != nil {
			return
		}

		for {
			select {
			case <-ticker.C:
				_ = r.emitTracesOnce(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *traceSimReceiver) emitTracesOnce(ctx context.Context) error {
	traces, err := r.simulator.Run(r.blueprint, time.Now())
	if err != nil {
		r.logger.Error("Error generating traces", zap.Error(err))
		return err
	}
	for _, trace := range traces {
		err = r.nextConsumer.ConsumeTraces(ctx, trace)
		if err != nil {
			r.logger.Error("Error sending traces", zap.Error(err))
		}
	}
	return nil
}

func (r *traceSimReceiver) Shutdown(_ context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}
