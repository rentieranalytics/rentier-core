package middleware

import (
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func SentryMiddleware(options sentryhttp.Options) func(http.Handler) http.Handler {
	sentryHandler := sentryhttp.New(options)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := &loggingResponseWrite{ResponseWriter: w, statusCode: http.StatusOK}
			sentryHandler.Handle(next).ServeHTTP(lrw, r)
			if lrw.statusCode == http.StatusInternalServerError {
				hub := sentry.GetHubFromContext(r.Context())
				if hub != nil {
					hub.CaptureMessage("500 Internal Server Error: " + r.URL.Path)
					hub.Flush(2 * time.Second)
				}
			}
		})
	}
}
