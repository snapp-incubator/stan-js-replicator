package pipe

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

// Piped contains metrics to meter the number of piped messages.
type Piped struct {
	PipedMessages  prometheus.Counter
	FailedMessages prometheus.Counter
	TimeLag        prometheus.Histogram
}

// NewPiped creates piped metrics based on topic name with const labels.
// nolint: funlen
func NewPiped(name string) Piped {
	piped := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "sjr",
		Name:      "piped_total",
		Help:      "total number of piped messages",
		Subsystem: "pipe",
		ConstLabels: prometheus.Labels{
			"topic": name,
		},
	})

	if err := prometheus.Register(piped); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &are); ok {
			piped, ok = are.ExistingCollector.(prometheus.Counter)
			if !ok {
				panic("piped must be a counter")
			}
		} else {
			panic(err)
		}
	}

	failed := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "sjr",
		Name:      "failed_total",
		Help:      "total number of failed messages",
		Subsystem: "pipe",
		ConstLabels: prometheus.Labels{
			"topic": name,
		},
	})

	if err := prometheus.Register(failed); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &are); ok {
			failed, ok = are.ExistingCollector.(prometheus.Counter)
			if !ok {
				panic("failed must be a counter")
			}
		} else {
			panic(err)
		}
	}

	lag := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "sjr",
		Name:      "time_lag_seconds",
		Help:      "message time lag",
		Subsystem: "pipe",
		Buckets:   prometheus.DefBuckets,
		ConstLabels: prometheus.Labels{
			"topic": name,
		},
	})

	if err := prometheus.Register(lag); err != nil {
		var are prometheus.AlreadyRegisteredError
		if ok := errors.As(err, &are); ok {
			lag, ok = are.ExistingCollector.(prometheus.Histogram)
			if !ok {
				panic("lag must be a histogram")
			}
		} else {
			panic(err)
		}
	}

	return Piped{
		PipedMessages:  piped,
		FailedMessages: failed,
		TimeLag:        lag,
	}
}
