package httpx

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/rentieranalytics/rentier-core/httpx/endpoints"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	fx.Provide(
		AsRoute(endpoints.NewPingHandler),
		fx.Annotate(
			NewServerMux,
			fx.ParamTags(`group:"routes"`),
		),
	),
	fx.Invoke(NewHttpServer),
)

func NewHttpServer(
	lc fx.Lifecycle,
	mux *http.ServeMux,
	logger *slog.Logger,
	config HttpServerConfig,
) *http.Server {
	srv := &http.Server{Addr: config.GetHttpServerAddress(), Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Info("Starting HTTP", slog.String("addr", srv.Addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
