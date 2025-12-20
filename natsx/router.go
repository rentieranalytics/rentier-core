package natsx

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	wmessage "github.com/ThreeDotsLabs/watermill/message"
	wmiddleware "github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
)

func NewRouter(
	handlers []SubscriberHandler,
	logger watermill.LoggerAdapter,
) *wmessage.Router {
	router, err := wmessage.NewRouter(wmessage.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		wmiddleware.Recoverer,
	)
	for _, h := range handlers {
		h.AddHandler(router)
	}
	return router
}

func RunRouter(router *wmessage.Router) {
	ctx := context.Background()
	go router.Run(ctx)
}
