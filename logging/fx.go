package logging

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var Module = fx.Module(
	"logging",
	fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
		return &fxevent.SlogLogger{
			Logger: log,
		}
	}),
	fx.Provide(
		NewLogger,
	),
)
