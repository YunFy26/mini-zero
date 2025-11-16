package logx

import (
	"context"
	"sync"
	"sync/atomic"
)

var (
	globalFields     atomic.Value
	globalFieldsLock sync.Mutex
)

type fieldsKey struct{}

func AddGlobalFields(fields ...LogField) {
	globalFieldsLock.Lock()
	defer globalFieldsLock.Unlock()
	old := globalFields.Load()
	if old == nil {
		globalFields.Store(append([]LogField(nil), fields...))
	} else {
		globalFields.Store(append(old.([]LogField), fields...))
	}
}

func ContextWithFields(ctx context.Context, fields ...LogField) context.Context {
	if val := ctx.Value(fieldsKey{}); val != nil {
		if arr, ok := val.([]LogField); ok {
			allFields := make([]LogField, 0, len(arr)+len(fields))
			allFields = append(allFields, arr...)
			allFields = append(allFields, fields...)
			return context.WithValue(ctx, fieldsKey{}, allFields)
		}
	}
	return context.WithValue(ctx, fieldsKey{}, fields)
}

func WithFields(ctx context.Context, fields ...LogField) context.Context {
	return ContextWithFields(ctx, fields...)
}
