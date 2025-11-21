package logx

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/YunFy26/mini-zero/core/lang"
)

const (
	hoursPerDay     = 24
	bufferSize      = 100
	defaultDirMode  = 0o755
	defaultFileMode = 0o600
	gzipExt         = ".gz"
	megaBytes       = 1 << 20
)

var (
	ErrorLogFileClosed = errors.New("error: log file closed")
	fileTimeFormat     = time.RFC3339
)

type (
	// RotateRule defines the interface for log rotation rules.
	RotateRule interface {
		BackupFileName() string
		MarkRotated()
		OutdatedFiles() []string
		ShallRotate(size int64) bool
	}

	// RotateLogger is a Logger that can rotate log files with given rules.
	RotateLogger struct {
		filename    string
		backup      string
		fp          *os.File
		channel     chan []byte
		done        chan lang.PlaceholderType
		rule        RotateRule
		compress    bool
		waitGroup   sync.WaitGroup
		closeOnce   sync.Once
		currentSize int64
	}

	// DailyRotateRule defines the daily rotation rule.
	DailyRotateRule struct {
		rotatedTime string
		filename    string
		delimiter   string
		days        int
		gzip        bool
	}

	// SizeLimitRotateRule defines the size limit rotation rule.
	SizeLimitRotateRule struct {
		DailyRotateRule
		maxSize    int64
		maxBackups int
	}
)

// ==================== DailyRotateRule Methods =========================

// BackupFileName returns the backup file name based on the current date.
func (r *DailyRotateRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, getNowDate())
}

// MarkRotated updates the rotated time to the current date.
func (r *DailyRotateRule) MarkRotated() {
	r.rotatedTime = getNowDate()
}

func (r *DailyRotateRule) OutdatedFiles() []string {
	if r.days < 0 {
		return nil
	}

	var pattern string
	if r.gzip {
		pattern = fmt.Sprintf("%s%s*%s", r.filename, r.delimiter, gzipExt)
	} else {
		pattern = fmt.Sprintf("%s%s*", r.filename, r.delimiter)
	}

}

func (r *DailyRotateRule) ShallRotate(size int64) bool {

}

func getNowDate() string {
	return time.Now().Format(time.DateOnly)
}

// NewSizeLimitRotateRule returns the rotation rule with size limit
func NewSizeLimitRotateRule(filename, delimiter string, days, maxSize, maxBackups int, gzip bool) RotateRule {
	return &SizeLimitRotateRule{
		DailyRotateRule: DailyRotateRule{
			rotatedTime: getNowDateInRFC3339Format(),
			filename:    filename,
			delimiter:   delimiter,
			days:        days,
			gzip:        gzip,
		},
		maxSize:    int64(maxSize) * megaBytes,
		maxBackups: maxBackups,
	}
}

func getNowDateInRFC3339Format() string {
	return time.Now().Format(fileTimeFormat)
}
