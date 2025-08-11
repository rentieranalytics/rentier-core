package sentry

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"sentry",
	fx.Invoke(func(c SentryConfig) {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              c.GetSentryDSN(),
			EnableTracing:    true,
			TracesSampleRate: c.GetSentrySampleRate(),
			Environment:      c.GetSentryEnv(),
			Release:          c.GetVersion(),
		}); err != nil {
			panic(err)
		}
	}),
)
