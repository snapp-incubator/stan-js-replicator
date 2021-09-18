package pipe

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

// Piped contains metrics to meter the number of piped messages.
type Piped struct {
	PipedMessages  prometheus.Counter
	FailedMessages prometheus.Counter
}

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
			piped, ok = are.ExistingCollector.(prometheus.Counter)
			if !ok {
				panic("failed must be a counter")
			}
		} else {
			panic(err)
		}
	}

	return Piped{
		PipedMessages:  piped,
		FailedMessages: failed,
	}
}
