package pipe

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/snapp-incubator/stan-js-replicator/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/internal/streaming"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Pipe struct {
	CMQ       *cmq.CMQ
	Streaming *streaming.Streaming
	Tracer    trace.Tracer
	Logger    *zap.Logger
}

// New create pipe to pipe data from streaming to jetstream.
// you can use a single pipe instance for mutiple topics.
func New(c *cmq.CMQ, s *streaming.Streaming, logger *zap.Logger, tracer trace.Tracer) *Pipe {
	var pipe Pipe

	pipe.CMQ = c
	pipe.Tracer = tracer
	pipe.Logger = logger.Named("pipe")
	pipe.Streaming = s

	return &pipe
}

// Pipe start piping messages from streaming to jetstream based on given topic.
// its subscription on streaming isn't durable and it always start from 1 second behind.
// the reason here is to reduce load on the streaming server as much as possible.
func (p *Pipe) Pipe(topic string) {
	if _, err := p.Streaming.Conn.QueueSubscribe(topic, p.Streaming.Group, func(imsg *stan.Msg) {
		defer func() {
			_ = imsg.Ack()
		}()

		ctx := context.Background()

		ctx, span := p.Tracer.Start(ctx, fmt.Sprintf("pipe.%s.replicate", topic), trace.WithSpanKind(trace.SpanKindProducer))

		omsg := new(nats.Msg)

		omsg.Subject = topic
		omsg.Data = imsg.Data
		omsg.Header = make(nats.Header)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(omsg.Header))

		if _, err := p.CMQ.JConn.PublishMsg(omsg); err != nil {
			span.RecordError(err)
			p.Logger.Error("stan subscription failed", zap.Error(err))
		}

		span.End()
	}, stan.StartAtTimeDelta(time.Second), stan.SetManualAckMode()); err != nil {
		p.Logger.Fatal("stan subscription failed", zap.Error(err))
	}
}
