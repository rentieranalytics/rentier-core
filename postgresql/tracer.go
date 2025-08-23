package postgresql

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Options struct {
	MaxSQLLen       int
	CaptureDBErrors bool
}

type Tracer struct {
	opts Options
}

func NewTracer(opts Options) *Tracer { return &Tracer{opts: opts} }

func (t *Tracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	if sentry.TransactionFromContext(ctx) == nil && sentry.SpanFromContext(ctx) == nil {
		return ctx
	}
	sql := compactSQL(data.SQL)
	if t.opts.MaxSQLLen > 0 && utf8.RuneCountInString(sql) > t.opts.MaxSQLLen {
		sql = string([]rune(sql)[:t.opts.MaxSQLLen]) + "…"
	}

	span := sentry.StartSpan(ctx, "db.postgresql.query",
		sentry.WithOpName("db.postgresql.query"),
		sentry.WithDescription(sql),
	)
	span.SetData("db.system", "postgresql")
	if conn != nil && conn.Config() != nil {
		cfg := conn.Config()
		span.SetData("server.address", cfg.Host)
		span.SetData("server.port", cfg.Port)
		span.SetData("db.user", cfg.User)
		span.SetData("db.name", cfg.Database)
	}

	return span.Context()
}

func (t *Tracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := sentry.SpanFromContext(ctx)
	if span == nil {
		return
	}
	if data.Err != nil {
		span.Status = sentry.SpanStatusInternalError
		var pgErr *pgconn.PgError
		if ok := errorAs(data.Err, &pgErr); ok && pgErr != nil {
			span.SetData("db.pg.code", pgErr.Code)
			span.SetData("db.pg.severity", pgErr.Severity)
			span.SetData("db.pg.detail", pgErr.Detail)
		}
		if t.opts.CaptureDBErrors {
			sentry.CaptureException(data.Err)
		}
	} else {
		span.Status = sentry.SpanStatusOK
	}
	span.Finish()
}

func (t *Tracer) TraceConnectStart(
	ctx context.Context,
	data pgx.TraceConnectStartData,
) context.Context {
	span := sentry.StartSpan(
		ctx,
		"db.postgresql.connect",
		sentry.WithOpName("db.postgresql.connect"),
	)
	return span.Context()
}

func (t *Tracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	span := sentry.SpanFromContext(ctx)
	if span == nil {
		return
	}
	if data.Err != nil {
		span.Status = sentry.SpanStatusInternalError
		sentry.CaptureException(data.Err)
	} else {
		span.Status = sentry.SpanStatusOK
	}
	span.Finish()
}

func compactSQL(s string) string {
	out := strings.TrimSpace(s)
	out = strings.Join(strings.Fields(out), " ")
	return out
}

func errorAs(err error, target interface{}) bool {
	type causer interface{ As(any) bool }
	if err == nil {
		return false
	}
	if e, ok := err.(causer); ok {
		return e.As(target)
	}
	return false
}
