package tracesimulationreceiver

import (
	"context"
	simulator "github.com/k4ji/tracesimulator/pkg"
	"github.com/k4ji/tracesimulator/pkg/blueprint"
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

		for {
			select {
			case <-ticker.C:
				traces, err := r.simulator.Run(r.blueprint, time.Now())
				if err != nil {
					r.logger.Error("Error generating traces %v, shutting down receiver", zap.Error(err))
					return
				}
				for _, trace := range traces {
					err = r.nextConsumer.ConsumeTraces(ctx, trace)
					if err != nil {
						r.logger.Error("Error sending traces", zap.Error(err))
						continue
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *traceSimReceiver) Shutdown(_ context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}
