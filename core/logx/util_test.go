package logx

import (
	"fmt"
	"testing"
)

func TestGetTimestamp(t *testing.T) {
	time := getTimestamp()
	if time == "" {
		t.Errorf("getTimestamp() = %v, want non-empty string", time)
	}
	fmt.Printf("getTimestamp 时间格式：%s\n", time)
}
