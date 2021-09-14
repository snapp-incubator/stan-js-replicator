package serve

import (
	"github.com/snapp-incubator/stan-js-replicator/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/internal/telemetry/metric"
	"github.com/snapp-incubator/stan-js-replicator/internal/telemetry/profiler"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, _ trace.Tracer) {
	profiler.Start(cfg.Telemetry.Profiler)

	metric.NewServer(cfg.Telemetry.Metric).Start(logger.Named("metrics"))

	c, err := cmq.New(cfg.NATS, logger)
	if err != nil {
		logger.Fatal("nats initiation failed", zap.Error(err))
	}

	if err := c.Stream(cfg.Channel, cfg.Topics); err != nil {
		logger.Fatal("nats stream creation failed", zap.Error(err))
	}
}

// Register server command.
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "serve",
			Short: "read events from streaming and publishes them on nats",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg, logger, tracer)
			},
		},
	)
}
