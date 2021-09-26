package config

import (
	"log"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/logger"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/streaming"
	telemetry "github.com/snapp-incubator/stan-js-replicator/v2/internal/telemetry/config"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "sjr_"

	StreamingInput Type = "streaming"
	NATSInput      Type = "nats"
)

type (
	Type string

	// Config holds all configurations.
	Config struct {
		Logger    logger.Config    `koanf:"logger"`
		Telemetry telemetry.Config `koanf:"telemetry"`
		Output    cmq.Config       `koanf:"output"`
		Input     Input            `koanf:"input"`
		Stream    Stream           `koanf:"stream"`
		Channel   string           `koanf:"channel"`
		Topics    []string         `koanf:"topics"`
	}

	// Input specifies the source of the data which is pushed into jectstream.
	Input struct {
		Type      `koanf:"type"`
		NATS      cmq.Config       `koanf:"nats"`
		Streaming streaming.Config `koanf:"streaming"`
	}

	// Stream holds all the stream configuration, please check it with
	// https://pkg.go.dev/github.com/nats-io/nats.go#StreamConfig
	Stream struct {
		Replicas    int              `koanf:"replicas"`
		MaxAge      time.Duration    `koanf:"maxage"`
		StorageType nats.StorageType `koanf:"storagetype"`
	}
)

// New reads configuration with viper.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	log.Printf("following configuration is loaded:\n%+v", instance)

	return instance
}
