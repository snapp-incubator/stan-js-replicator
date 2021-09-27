package pipe_test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/cmq"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/config"
	"github.com/snapp-incubator/stan-js-replicator/v2/internal/pipe"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PipeNATSSuite struct {
	suite.Suite

	nats *cmq.CMQ
	js   *cmq.CMQ

	pipe pipe.Pipe
}

func (suite *PipeNATSSuite) SetupSuite() {
	cfg := config.New()
	require := suite.Require()

	nats, err := cmq.New(cfg.Input.NATS, zap.NewNop())
	require.NoError(err)

	js, err := cmq.New(cfg.Output, zap.NewNop())
	require.NoError(err)

	suite.pipe = pipe.NewNATS(js, nats, zap.NewNop(), trace.NewNoopTracerProvider().Tracer(""))
	suite.nats = nats
	suite.js = js
}

func (suite *PipeNATSSuite) TestWithMessage() {
	require := suite.Require()

	// nolint: exhaustivestruct
	require.NoError(suite.js.Stream(&nats.StreamConfig{
		Name:     "hello.nats",
		Storage:  nats.MemoryStorage,
		MaxAge:   time.Minute,
		Subjects: []string{"hello.world.nats"},
	}))

	suite.pipe.Pipe("hello.world.nats", "group")

	sub, err := suite.js.JConn.SubscribeSync("hello.world.nats", nats.AckExplicit(), nats.DeliverAll())
	require.NoError(err)

	defer func() {
		_ = sub.Unsubscribe()
	}()

	require.NoError(suite.nats.Conn.Publish("hello.world.nats", []byte("Hello World")))

	msg, err := sub.NextMsg(Timeout)
	require.NoError(err)

	require.Equal(msg.Subject, "hello.world.nats")
	require.Equal(msg.Data, []byte("Hello World"))

	meta, err := msg.Metadata()
	require.NoError(err)

	require.Equal(meta.Stream, "hello.nats")
}

func TestPipe(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(PipeNATSSuite))
}
