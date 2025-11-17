package logx

import (
	"fmt"
	"reflect"
	"sync"
)

const callerDepth = 5

var (
	timeFormat              = "2006-01-02T15:04:05.000Z07:00"
	encoding         uint32 = jsonEncodingType
	maxContentLength uint32
	disableStat      uint32
	logLevel         uint32
	options          logOptions
	writer           = new(atomicWriter)
	setupOnce        sync.Once
)

type (
	// LogField is a key-value pair that will be added to the log entry.
	LogField struct {
		Key   string
		Value any
	}

	// LogOption defines the method to customize the logging.
	LogOption func(options *logOptions)

	logEntry map[string]any

	logOptions struct {
		gzipEnabled           bool
		logStackCooldownMills int
		keepDays              int
		maxBackups            int
		maxSize               int
		rotationRule          string
	}
)

func Field(key string, value any) LogField {
	return LogField{
		Key:   key,
		Value: value,
	}
}

// err -> string
func encodeError(err error) (ret string) {
	return encodeWithRecover(err, func() string {
		// 可能会NPE
		return err.Error()
	})
}

// stringer -> string
func encodeStringer(s fmt.Stringer) string {
	return encodeWithRecover(s, func() string {
		// 可能会NPE
		return s.String()
	})
}

// recover() panic
func encodeWithRecover(arg any, fn func() string) (ret string) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(arg); v.Kind() == reflect.Ptr && v.IsNil() {
				ret = nilAngleString
			} else {
				ret = fmt.Sprintf("panic: %v", err)
			}
		}
	}()
	return fn()
}
