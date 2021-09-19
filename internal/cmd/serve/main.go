package serve

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/internal/pipe"
	"github.com/snapp-incubator/stan-js-replicator/internal/streaming"
	"github.com/snapp-incubator/stan-js-replicator/internal/telemetry/metric"
	"github.com/snapp-incubator/stan-js-replicator/internal/telemetry/profiler"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	profiler.Start(cfg.Telemetry.Profiler)

	metric.NewServer(cfg.Telemetry.Metric).Start(logger.Named("metrics"))

	c, err := cmq.New(cfg.NATS, logger.Named("cmq"))
	if err != nil {
		logger.Fatal("nats initiation failed", zap.Error(err))
	}

	// nolint: exhaustivestruct
	sc := &nats.StreamConfig{
		Name:     cfg.Channel,
		Subjects: cfg.Topics,
		Replicas: cfg.Stream.Replicas,
		Storage:  cfg.Stream.StorageType,
		MaxAge:   cfg.Stream.MaxAge,
	}

	if err := c.Stream(sc); err != nil {
		logger.Fatal("nats stream creation failed", zap.Error(err))
	}

	str, err := streaming.New(cfg.Streaming, logger.Named("streaming"))
	if err != nil {
		logger.Fatal("nats streaming initiation failed", zap.Error(err))
	}

	p := pipe.New(c, str, logger.Named("pipe"), tracer)

	for _, topic := range cfg.Topics {
		go p.Pipe(topic)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
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
