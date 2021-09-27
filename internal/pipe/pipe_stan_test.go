package pipe_test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/pipe"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/streaming"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const Timeout = 10 * time.Second

type PipeSTANSuite struct {
	suite.Suite

	stan *streaming.Streaming
	js   *cmq.CMQ

	pipe pipe.Pipe
}

func (suite *PipeSTANSuite) SetupSuite() {
	cfg := config.New()
	require := suite.Require()

	stan, err := streaming.New(cfg.Input.Streaming, zap.NewNop())
	require.NoError(err)

	js, err := cmq.New(cfg.Output, zap.NewNop())
	require.NoError(err)

	suite.pipe = pipe.NewSTAN(js, stan, zap.NewNop(), trace.NewNoopTracerProvider().Tracer(""))
	suite.stan = stan
	suite.js = js
}

// nolint: dupl
func (suite *PipeSTANSuite) TestWithMessage() {
	require := suite.Require()

	// nolint: exhaustivestruct
	require.NoError(suite.js.Stream(&nats.StreamConfig{
		Name:     "hello-stan",
		Storage:  nats.MemoryStorage,
		MaxAge:   time.Minute,
		Subjects: []string{"hello.world.stan"},
	}))

	suite.pipe.Pipe("hello.world.stan", "group")

	sub, err := suite.js.JConn.SubscribeSync("hello.world.stan", nats.AckExplicit(), nats.DeliverAll())
	require.NoError(err)

	defer func() {
		_ = sub.Unsubscribe()
	}()

	require.NoError(suite.stan.Conn.Publish("hello.world.stan", []byte("Hello World")))

	msg, err := sub.NextMsg(Timeout)
	require.NoError(err)

	require.Equal(msg.Subject, "hello.world.stan")
	require.Equal(msg.Data, []byte("Hello World"))

	meta, err := msg.Metadata()
	require.NoError(err)

	require.Equal(meta.Stream, "hello.stan")
}

func TestSTANPipe(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(PipeSTANSuite))
}
