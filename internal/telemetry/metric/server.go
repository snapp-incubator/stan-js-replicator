package metric

import (
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/snapp-incubator/stan-js-replicator/internal/telemetry/config"
	"go.uber.org/zap"
)

// Server contains information about metrics server.
type Server struct {
	srv     *http.ServeMux
	address string
}

// NewServer creates a new monitoring server.
func NewServer(cfg config.Metric) Server {
	var srv *http.ServeMux

	if cfg.Enabled {
		srv = http.NewServeMux()
		srv.Handle("/metrics", promhttp.Handler())
	}

	return Server{
		address: cfg.Address,
		srv:     srv,
	}
}

// Start creates and run a metric server for prometheus in new go routine.
func (s Server) Start(logger *zap.Logger) {
	go func() {
		if err := http.ListenAndServe(s.address, s.srv); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("metric server initiation failed", zap.Error(err))
		}
	}()
}
