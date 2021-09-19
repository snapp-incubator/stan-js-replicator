package cmq

import (
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type CMQ struct {
	Conn   *nats.Conn
	JConn  nats.JetStreamContext
	Logger *zap.Logger
}

// New creates a new connection to nats cluster with jetstream support.
func New(cfg Config, logger *zap.Logger) (*CMQ, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed %w", err)
	}

	logger.Info("nats connection successful",
		zap.String("connected-addr", nc.ConnectedAddr()),
		zap.Strings("discovered-servers", nc.DiscoveredServers()))

	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logger.Fatal("nats disconnected", zap.Error(err))
	})

	nc.SetReconnectHandler(func(c *nats.Conn) {
		logger.Warn("nats reconnected")
	})

	jsm, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("jetstream context extraction failed %w", err)
	}

	return &CMQ{
		Conn:   nc,
		JConn:  jsm,
		Logger: logger,
	}, nil
}

func (c *CMQ) Stream(sc *nats.StreamConfig) error {
	info, err := c.JConn.StreamInfo(sc.Name)

	switch {
	case errors.Is(err, nats.ErrStreamNotFound):
		stream, err := c.JConn.AddStream(sc)
		if err != nil {
			return fmt.Errorf("cannot create stream %w", err)
		}

		info = stream
	case err != nil:
		return fmt.Errorf("cannot read stream information %w", err)
	}

	c.Logger.Info("events stream", zap.Any("stream", info))

	return nil
}
