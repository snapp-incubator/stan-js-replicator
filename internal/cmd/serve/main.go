package serve

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/pipe"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/streaming"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/telemetry/metric"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/telemetry/profiler"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	profiler.Start(cfg.Telemetry.Profiler)

	metric.NewServer(cfg.Telemetry.Metric).Start(logger.Named("metrics"))

	c, err := cmq.New(cfg.Output, logger.Named("cmq"))
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

	var p pipe.Pipe

	switch cfg.Input.Type {
	case config.StreamingInput:
		str, err := streaming.New(cfg.Input.Streaming, logger.Named("streaming"))
		if err != nil {
			logger.Fatal("nats streaming initiation failed", zap.Error(err))
		}

		p = pipe.NewSTAN(c, str, logger.Named("pipe"), tracer)
	case config.NATSInput:
		nats, err := cmq.New(cfg.Input.NATS, logger.Named("nats"))
		if err != nil {
			logger.Fatal("nats streaming initiation failed", zap.Error(err))
		}

		p = pipe.NewNATS(c, nats, logger.Named("pipe"), tracer)
	default:
		logger.Fatal("invalid input type", zap.String("input-type", string(cfg.Input.Type)))
	}

	for _, topic := range cfg.Topics {
		go p.Pipe(topic, cfg.Input.Group)
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
