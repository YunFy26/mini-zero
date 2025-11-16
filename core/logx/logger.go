package logx

import (
	"context"
	"time"
)

type Logger interface {
	Debug(...any)
	Debugf(string, ...any)
	Debugfn(func() any)
	Debugv(any)
	Debugw(string, ...LogField)

	Info(...any)
	Infof(string, ...any)
	Infofn(func() any)
	Infov(any)
	Infow(string, ...LogField)

	Slow(...any)
	Slowf(string, ...any)
	Slowfn(func() any)
	Slowv(any)
	Sloww(string, ...LogField)

	WithCallerSkip(skip int) Logger
	WithContext(ctx context.Context) Logger
	WithDuration(d time.Duration) Logger
	WithFields(fields ...LogField) Logger
}
