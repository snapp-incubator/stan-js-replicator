package config

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/logger"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/streaming"
	telemetry "github.com/snapp-incubator/stan-js-replicator/v2/internal/telemetry/config"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Logger: logger.Config{
			Level: "debug",
			Syslog: logger.Syslog{
				Enabled: false,
				Network: "",
				Address: "",
				Tag:     "",
			},
		},
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: false,
				Ratio:   1.0,
				Agent: telemetry.Agent{
					Host: "127.0.0.1",
					Port: "6831",
				},
			},
			Profiler: telemetry.Profiler{
				Enabled: false,
				Address: "http://127.0.0.1:4040",
			},
			Metric: telemetry.Metric{
				Address: ":8080",
				Enabled: true,
			},
		},
		Output: cmq.Config{
			URL: "nats://127.0.0.1:4222",
		},
		Input: Input{
			Type:  StreamingInput,
			Group: "sjr",
			Streaming: streaming.Config{
				URL:       "nats://127.0.0.1:4223",
				ClusterID: "snapp",
			},
			NATS: cmq.Config{
				URL: "nats://127.0.0.1:4224",
			},
		},
		Channel: "koochooloo",
		Topics:  []string{"k1", "k2"},
		Stream: Stream{
			MaxAge:      1 * time.Hour,
			StorageType: nats.MemoryStorage,
			Replicas:    1,
		},
	}
}
