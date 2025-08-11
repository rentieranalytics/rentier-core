package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rentier-io/rentier-core/httpx/middleware"
)

type Ping struct{}

func NewPingHandler(logger *slog.Logger) *Ping {
	return &Ping{}
}

func (t *Ping) Pattern() string {
	return "GET /ping"
}

func (h *Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}

func (h *Ping) Middlewares() []middleware.Middleware {
	return []middleware.Middleware{}
}
