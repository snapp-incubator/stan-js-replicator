package profiler

import (
	"log"

	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
	"github.com/snapp-incubator/stan-js-replicator/internal/telemetry/config"
)

func Start(cfg config.Profiler) {
	if cfg.Enabled {
		// nolint: exhaustivestruct
		if _, err := profiler.Start(profiler.Config{
			ApplicationName: "snapp.stan-js-replicator",
			ServerAddress:   cfg.Address,
		}); err != nil {
			log.Printf("failed to start the profiler")
		}
	}
}
