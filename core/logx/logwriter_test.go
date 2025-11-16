package logx

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	timestamp := time.Now().Format("15:04:05.000")
	testMessage := fmt.Sprintf("[%s] test log data\n", timestamp)
	data := []byte(testMessage)

	logger := newLogWriter(log.New(os.Stdout, "", flags))

	n, err := logger.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Expected to write %d bytes, but wrote %d bytes", len(data), n)
	}
}
