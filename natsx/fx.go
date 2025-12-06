package natsx

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"natsx",
	fx.Provide(
		fx.Annotate(
			NewWatermillLogger,
			fx.As(new(watermill.LoggerAdapter)),
		),
		fx.Annotate(
			NewRouter,
			fx.ParamTags(`group:"subHandlers"`),
		),
		NewNatsConnection,
		NewNatsJetStream,
	),
	fx.Invoke(
		RunRouter,
	),
)

type SubscriberHandler interface {
	AddHandler(*message.Router)
}

func NewWatermillLogger(logger *slog.Logger) watermill.LoggerAdapter {
	return watermill.NewSlogLogger(logger)
}

func NewNatsConnection(config NatsConfig, logger *slog.Logger) *nats.Conn {
	conn, natsConnectErr := nats.Connect(
		config.GetNatsURL(),
		nats.UserCredentials(config.GetNatsJWTUserFilePath()),
		nats.Name(config.GetClientName()),
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error(
				"Disconnected due to error, will attempt reconnects",
				slog.Any("error", err),
			)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info(
				"Reconnected to nats",
				slog.Any("url", nc.ConnectedUrl()),
			)
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Info(
				"Connection closed",
				slog.Any("reason", nc.LastError()),
			)
		}),
	)
	if natsConnectErr != nil {
		panic(natsConnectErr)
	}
	return conn
}

func NewNatsJetStream(
	con *nats.Conn,
) (jetstream.JetStream, error) {
	js, err := jetstream.New(con)
	if err != nil {
		return nil, err
	}
	return js, nil
}
