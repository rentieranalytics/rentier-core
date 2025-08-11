package httpx

import (
	"net/http"

	"github.com/rentier-io/rentieranalytics/httpx/middleware"
)

type HttpServerConfig interface {
	GetHttpServerAddress() string
}

type Route interface {
	http.Handler
	Pattern() string
	Middlewares() []middleware.Middleware
}

func NewServerMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		f := http.HandlerFunc(route.ServeHTTP)
		h := chainMiddlewares(f, route.Middlewares()...)
		mux.Handle(route.Pattern(), h)
	}
	return mux
}

func chainMiddlewares(
	handler http.Handler,
	middlewares ...middleware.Middleware,
) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
