package pipe

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NATSPipe struct {
	CMQ    *cmq.CMQ
	NATS   *cmq.CMQ
	Tracer trace.Tracer
	Logger *zap.Logger
}

// New create pipe to pipe data from nats to jetstream.
// you can use a single pipe instance for mutiple topics.
func NewNATS(c *cmq.CMQ, n *cmq.CMQ, logger *zap.Logger, tracer trace.Tracer) Pipe {
	var pipe NATSPipe

	pipe.CMQ = c
	pipe.Tracer = tracer
	pipe.Logger = logger.Named("pipe")
	pipe.NATS = n

	return &pipe
}

// Pipe start piping messages from nats to jetstream based on given topic.
func (p *NATSPipe) Pipe(topic, group string) {
	piped := NewPiped(topic)

	if _, err := p.NATS.Conn.QueueSubscribe(topic, group, func(imsg *nats.Msg) {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(imsg.Header))

		ctx, span := p.Tracer.Start(ctx, fmt.Sprintf("pipe.%s.replicate", topic), trace.WithSpanKind(trace.SpanKindProducer))

		omsg := new(nats.Msg)

		omsg.Subject = topic
		omsg.Data = imsg.Data
		omsg.Header = make(nats.Header)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(omsg.Header))

		if _, err := p.CMQ.JConn.PublishMsg(omsg); err != nil {
			piped.FailedMessages.Inc()
			span.RecordError(err)
			p.Logger.Error("jetstream publish failed", zap.Error(err))
		}

		span.End()
		piped.PipedMessages.Inc()
	}); err != nil {
		p.Logger.Fatal("nats subscription failed", zap.Error(err))
	}
}
