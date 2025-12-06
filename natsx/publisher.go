package natsx

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	wjetstream "github.com/ThreeDotsLabs/watermill-nats/v2/pkg/jetstream"
	wmessage "github.com/ThreeDotsLabs/watermill/message"
)

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
