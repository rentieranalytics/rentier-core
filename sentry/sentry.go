package sentry

type SentryConfig interface {
	GetSentryDSN() string
	GetSentryEnv() string
	GetSentrySampleRate() float64
	GetVersion() string
	GetServerName() string
}
