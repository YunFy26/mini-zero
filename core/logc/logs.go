package logc

import (
	"context"
	"fmt"

	"github.com/YunFy26/mini-zero/core/logx"
)

type (
	LogConf  = logx.LogConf
	LogField = logx.LogField
)

func AddGlobalFields(fields ...LogField) {
	logx.AddGlobalFields(fields...)
}

func Alert(_ context.Context, v string) {
	logx.Alert(v)
}

func Close() error {
	return logx.Close()
}

func Debug(ctx context.Context, v ...interface{}) {
	getLogger(ctx).Debug(v...)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	getLogger(ctx).Debugf(format, v...)
}

func Debugfn(ctx context.Context, fn func() any) {
	getLogger(ctx).Debugfn(fn)
}

func Debugv(ctx context.Context, v interface{}) {
	getLogger(ctx).Debugv(v)
}

func Debugw(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Debugw(msg, fields...)
}

func Error(ctx context.Context, v ...any) {
	getLogger(ctx).Error(v...)
}

func Errorf(ctx context.Context, format string, v ...any) {
	getLogger(ctx).Errorf(fmt.Errorf(format, v...).Error())
}

func Errorfn(ctx context.Context, fn func() any) {
	getLogger(ctx).Errorfn(fn)
}

func Errorv(ctx context.Context, v any) {
	getLogger(ctx).Errorv(v)
}

func Errorw(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Errorw(msg, fields...)
}

func Field(key string, value any) LogField {
	return logx.Field(key, value)
}

func Info(ctx context.Context, v ...any) {
	getLogger(ctx).Info(v...)
}

func Infof(ctx context.Context, format string, v ...any) {
	getLogger(ctx).Infof(format, v...)
}

func Infofn(ctx context.Context, fn func() any) {
	getLogger(ctx).Infofn(fn)
}

func Infov(ctx context.Context, v any) {
	getLogger(ctx).Infov(v)
}

func Infow(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Infow(msg, fields...)
}

func Must(err error) {
	logx.Must(err)
}

func MustSetup(c logx.LogConf) {
	logx.MustSetup(c)
}

func SetLevel(level uint32) {
	logx.SetLevel(level)
}

func SetUp(c LogConf) error {
	return logx.SetUp(c)
}

func Slow(ctx context.Context, v ...any) {
	getLogger(ctx).Slow(v...)
}

func Slowf(ctx context.Context, format string, v ...any) {
	getLogger(ctx).Slowf(format, v...)
}

func Slowfn(ctx context.Context, fn func() any) {
	getLogger(ctx).Slowfn(fn)
}

func Slowv(ctx context.Context, v any) {
	getLogger(ctx).Slowv(v)
}

func Sloww(ctx context.Context, msg string, fields ...LogField) {
	getLogger(ctx).Sloww(msg, fields...)
}

func getLogger(ctx context.Context) logx.Logger {
	return logx.WithContext(ctx).WithCallerSkip(1)
}
