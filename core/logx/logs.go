package logx

import "sync"

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
