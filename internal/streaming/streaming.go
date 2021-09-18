package streaming

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

const ClientIDLen = 15

type Streaming struct {
	Conn   stan.Conn
	Logger *zap.Logger
	Group  string
}

// New creates a new connection to nats cluster with streaming support.
func New(cfg Config, logger *zap.Logger) (*Streaming, error) {
	id := make([]byte, ClientIDLen)

	if _, err := rand.Read(id); err != nil {
		return nil, fmt.Errorf("cannot create random stan id %w", err)
	}

	nc, err := stan.Connect(
		cfg.ClusterID,
		fmt.Sprintf("sjr-%s", base64.URLEncoding.EncodeToString(id)),
		stan.NatsURL(cfg.URL),
	)
	if err != nil {
		return nil, fmt.Errorf("stan connection failed %w", err)
	}

	logger.Info("nats connection successful",
		zap.String("connected-addr", nc.NatsConn().ConnectedAddr()),
		zap.Strings("discovered-servers", nc.NatsConn().DiscoveredServers()))

	return &Streaming{
		Conn:   nc,
		Logger: logger,
		Group:  cfg.Group,
	}, nil
}
