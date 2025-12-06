package natsx

import (
	"go.uber.org/fx"

	"context"
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wjetstream "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type CommonSubscriberDeps struct {
	fx.In

	Logger  *slog.Logger
	MakeSub NatsSubscriberFactory
}

type NatsSubscriberConfig struct {
	ConsumerName   string
	FilterSubject  string
	FilterSubjects []string
	Durable        string
	AckWaitTimeout time.Duration
	AckPolicy      jetstream.AckPolicy
	NakDelay       time.Duration
	DeliverPolicy  jetstream.DeliverPolicy
}

type NatsSubscriberFactory func(cfg NatsSubscriberConfig) *wjetstream.Subscriber

func NewNatsSubscriberFactory(
	logger *slog.Logger,
	con *nats.Conn,
) NatsSubscriberFactory {
	return func(cfg NatsSubscriberConfig) *wjetstream.Subscriber {
		if cfg.Durable == "" {
			cfg.Durable = cfg.ConsumerName
		}
		if cfg.AckWaitTimeout == 0 {
			cfg.AckWaitTimeout = 300 * time.Second
		}
		if cfg.NakDelay == 0 {
			cfg.NakDelay = 1 * time.Minute
		}
		if cfg.DeliverPolicy == 0 {
			cfg.DeliverPolicy = jetstream.DeliverNewPolicy
		}

		l := watermill.NewSlogLogger(logger)

		sub, err := wjetstream.NewSubscriber(
			wjetstream.SubscriberConfig{
				Conn:           con,
				Logger:         l,
				AckWaitTimeout: cfg.AckWaitTimeout,
				NakDelay:       wjetstream.NewStaticDelay(cfg.NakDelay),
				ResourceInitializer: func(
					ctx context.Context,
					js jetstream.JetStream,
					topic string,
				) (
					jetstream.Consumer,
					func(context.Context, watermill.LoggerAdapter),
					error,
				) {
					consumerCfg := jetstream.ConsumerConfig{
						Durable:        cfg.Durable,
						AckPolicy:      cfg.AckPolicy,
						FilterSubject:  cfg.FilterSubject,
						FilterSubjects: cfg.FilterSubjects,
						DeliverPolicy:  cfg.DeliverPolicy,
					}

					if cfg.FilterSubject != "" {
						consumerCfg.FilterSubject = cfg.FilterSubject
					}
					if len(cfg.FilterSubjects) > 0 {
						consumerCfg.FilterSubjects = cfg.FilterSubjects
					}
					c, ce := js.CreateOrUpdateConsumer(
						ctx,
						topic,
						consumerCfg,
					)
					if ce != nil {
						return nil, nil, ce
					}
					return c, func(context.Context, watermill.LoggerAdapter) {}, nil
				},
			})
		if err != nil {
			panic(err)
		}
		return sub
	}
}
