package sentry

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/getsentry/sentry-go"
)

type SentryConfig interface {
	GetSentryDSN() string
	GetSentryEnv() string
	GetSentrySampleRate() float64
	GetVersion() string
	GetServerName() string
}

type Trace struct {
	Trace  string
	Bagged string
}

func NewTrace(trace, bagged string) *Trace {
	return &Trace{
		Trace:  trace,
		Bagged: bagged,
	}
}

func GetTXFromMsg(
	ctx context.Context,
	msg *message.Message,
	txName string,
) *sentry.Span {
	md := msg.Metadata
	traceFromReq := NewTrace(
		md.Get("sentry-trace"),
		md.Get("baggage"),
	)
	return StartTransaction(
		ctx,
		txName,
		traceFromReq,
	)
}

func StartTransaction(
	ctx context.Context,
	title string,
	trace *Trace,
) *sentry.Span {
	var tx *sentry.Span
	if trace != nil {
		tx = sentry.StartTransaction(
			ctx,
			title,
			sentry.ContinueFromHeaders(trace.Trace, trace.Bagged),
			sentry.WithTransactionSource(sentry.SourceCustom),
		)
	} else {
		tx = sentry.StartTransaction(
			ctx,
			title,
		)
	}
	return tx
}

func SetTxError(ctx context.Context, tx *sentry.Span, err error) {
	tx.Status = sentry.SpanStatusInternalError
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
	}
	hub.Scope().SetSpan(tx)
	hub.CaptureException(err)
}

type contextKeySpan struct{}

func ContextWithSpan(ctx context.Context, span *sentry.Span) context.Context {
	return context.WithValue(ctx, contextKeySpan{}, span)
}

func TxToCtx(
	ctx context.Context,
	msg *message.Message,
	subsriberName string,
) (context.Context, *sentry.Span) {
	tx := GetTXFromMsg(ctx, msg, subsriberName)
	return ContextWithSpan(ctx, tx), tx
}

func SpanFromContext(ctx context.Context) *sentry.Span {
	if s, ok := ctx.Value(contextKeySpan{}).(*sentry.Span); ok {
		return s
	}
	return nil
}

func SetTags(tx *sentry.Span, data map[string]string) {
	for key, value := range data {
		tx.SetTag(key, value)
	}
}
