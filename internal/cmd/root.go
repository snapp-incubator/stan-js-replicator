package cmd

import (
	"os"

	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmd/serve"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/logger"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/telemetry/trace"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	tracer := trace.New(cfg.Telemetry.Trace)

	// nolint: exhaustivestruct
	root := &cobra.Command{
		Use:   "sjr",
		Short: "replicate streaming messages on jetstream",
	}

	serve.Register(root, cfg, logger, tracer)

	if err := root.Execute(); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
