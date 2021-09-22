package pipe_test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/internal/pipe"
	"github.com/snapp-incubator/stan-js-replicator/internal/streaming"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const Timeout = 10 * time.Second

type PipeSuite struct {
	suite.Suite

	stan *streaming.Streaming
	js   *cmq.CMQ

	pipe *pipe.Pipe
}

func (suite *PipeSuite) SetupSuite() {
	cfg := config.New()
	require := suite.Require()

	stan, err := streaming.New(cfg.Streaming, zap.NewNop())
	require.NoError(err)

	js, err := cmq.New(cfg.NATS, zap.NewNop())
	require.NoError(err)

	suite.pipe = pipe.New(js, stan, zap.NewNop(), trace.NewNoopTracerProvider().Tracer(""))
	suite.stan = stan
	suite.js = js
}

func (suite *PipeSuite) TestWithMessage() {
	require := suite.Require()

	// nolint: exhaustivestruct
	require.NoError(suite.js.Stream(&nats.StreamConfig{
		Name:     "hello",
		Storage:  nats.MemoryStorage,
		MaxAge:   time.Minute,
		Subjects: []string{"hello.world"},
	}))

	suite.pipe.Pipe("hello.world")

	sub, err := suite.js.JConn.SubscribeSync("hello.world", nats.AckExplicit(), nats.DeliverAll())
	require.NoError(err)

	defer func() {
		_ = sub.Unsubscribe()
	}()

	require.NoError(suite.stan.Conn.Publish("hello.world", []byte("Hello World")))

	msg, err := sub.NextMsg(Timeout)
	require.NoError(err)

	require.Equal(msg.Subject, "hello.world")
	require.Equal(msg.Data, []byte("Hello World"))

	meta, err := msg.Metadata()
	require.NoError(err)

	require.Equal(meta.Stream, "hello")
}

func TestPipe(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(PipeSuite))
}
