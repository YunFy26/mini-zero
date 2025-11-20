package logx

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/YunFy26/mini-zero/core/lang"
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

func (r *DailyRotateRule) BackupFileName() string {
	return fmt.Sprintf("%s%s%s", r.filename, r.delimiter, getNowDate())
}

func getNowDate() string {
	return time.Now().Format(time.DateOnly)
}
