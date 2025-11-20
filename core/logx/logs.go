package logx

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
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

// ============================= Log Writer ===============================
// type (
// 	Writer interface {
// 		Alert(v any)                     // 警告日志
// 		Close() error                    // 关闭写入器
// 		Debug(v any, fields ...LogField) // 调试日志
// 		Error(v any, fields ...LogField) // 错误日志
// 		Info(v any, fields ...LogField)  // 信息日志
// 		Severe(v any)                    // 严重日志
// 		Slow(v any, fields ...LogField)  // 慢查询日志
// 		Stack(v any)                     // 堆栈日志
// 		Stat(v any, fields ...LogField)  // 统计日志
// 	}

// 	// 原子写入
// 	atomicWriter struct {
// 		writer Writer
// 		lock   sync.RWMutex
// 	}
// 	// 组合写入
// 	comboWriter struct {
// 		writers []Writer
// 	}
// 	// 具体的日志写入器实现
// 	concreteWriter struct {
// 		infoLog   io.WriteCloser
// 		errorLog  io.WriteCloser
// 		severeLog io.WriteCloser
// 		slowLog   io.WriteCloser
// 		statLog   io.WriteCloser
// 		stackLog  io.Writer
// 	}

//	nopWriter struct{}
//
// )

// AddWriter adds a log writer, supporting multiple writers writing simultaneously
func AddWriter(w Writer) {
	ow := Reset()
	if ow == nil {
		SetWriter(w)
	} else {
		SetWriter(comboWriter{
			writers: []Writer{ow, w},
		})
	}
}

// getWriter gets the current log writer
func getWriter() Writer {
	w := writer.Load()
	if w == nil {
		w = writer.StoreIfNil(newConsoleWriter())
	}
	return w
}

// Reset the current log writer and return the old writer
func Reset() Writer {
	return writer.Swap(nil)
}

// SetWriter sets the log writer
func SetWriter(w Writer) {
	if atomic.LoadUint32(&logLevel) != disableLevel {
		writer.Store(w)
	}
}

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

// ============================= Log Options ===============================

//	logOptions struct {
//		gzipEnabled           bool
//		logStackCooldownMills int
//		keepDays              int
//		maxBackups            int
//		maxSize               int
//		rotationRule          string
//	}
//
// WithCoolDownMillis sets the cooldown milliseconds for logging stack traces
func WithCoolDownMillis(millis int) LogOption {
	return func(options *logOptions) {
		options.logStackCooldownMills = millis
	}
}

// WithKeepDays sets the number of days to keep log files
func WithKeepDays(days int) LogOption {
	return func(options *logOptions) {
		options.keepDays = days
	}
}

// WithGzip enables gzip compression for log files
func WithGzip() LogOption {
	return func(options *logOptions) {
		options.gzipEnabled = true
	}
}

// WithMaxBackups sets the maximum number of log file backups
func WithMaxBackups(count int) LogOption {
	return func(options *logOptions) {
		options.maxBackups = count
	}
}

// WithMaxSize sets the maximum size of a single log file in megabytes
func WithMaxSize(size int) LogOption {
	return func(options *logOptions) {
		options.maxSize = size
	}
}

// WithRotation sets the rotation rule for log files
func WithRotation(r string) LogOption {
	return func(options *logOptions) {
		options.rotationRule = r
	}
}

// handleOptions applies the given log options to the global options
func handleOptions(opts []LogOption) {
	for _, opt := range opts {
		opt(&options)
	}
}

func createOutput(path string) (io.WriteCloser, error) {
	if len(path) == 0 {
		return nil, ErrLogPathNotSet
	}

	var rule RotateRule
}
