package natsx

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"

	wjetstream "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
	wmessage "github.com/ThreeDotsLabs/watermill/message"
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
		nats.Name("api-estimator"),
		nats.ReconnectWait(2*time.Second), // czas oczekiwania między próbami reconnectu
		nats.MaxReconnects(-1),            // -1 oznacza nieskończone próby reconnectu
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

func NewNatsPublisher(
	logger watermill.LoggerAdapter,
	con *nats.Conn,
	sConfig jetstream.StreamConfig,
) wmessage.Publisher {
	publisher, err := wjetstream.NewPublisher(
		wjetstream.PublisherConfig{
			Conn:   con,
			Logger: logger,
		},
	)
	if err != nil {
		panic(err)
	}
	return publisher

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
