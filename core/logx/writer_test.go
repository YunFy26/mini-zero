package logx

import (
	"bytes"
	"strings"
	"testing"
)

func TestWritePlainText(t *testing.T) {
	t.Run("测试基础日志写入", func(t *testing.T) {
		var buf bytes.Buffer
		writePlainText(&buf, "INFO", "测试消息")

		output := buf.String()
		if !strings.Contains(output, "INFO") {
			t.Errorf("日志应包含级别信息")
		}
		if !strings.Contains(output, "测试消息") {
			t.Errorf("日志应包含消息内容")
		}
	})

	t.Run("测试带字段的日志", func(t *testing.T) {
		var buf bytes.Buffer
		writePlainText(&buf, "ERROR", "错误发生", "field1", "field2")

		output := buf.String()
		if !strings.Contains(output, "ERROR") {
			t.Errorf("日志应包含错误级别")
		}
		if !strings.Contains(output, "field1") {
			t.Errorf("日志应包含字段1")
		}
		if !strings.Contains(output, "field2") {
			t.Errorf("日志应包含字段2")
		}
	})

	t.Run("测试writer为nil的情况", func(t *testing.T) {
		// 这个测试主要验证不会panic
		writePlainText(nil, "DEBUG", "nil writer测试")
	})
}

func TestWritePlainValue(t *testing.T) {
	t.Run("测试结构体日志", func(t *testing.T) {
		var buf bytes.Buffer

		type LogData struct {
			User    string `json:"user"`
			Action  string `json:"action"`
			Success bool   `json:"success"`
		}

		data := LogData{User: "testuser", Action: "login", Success: true}
		writePlainValue(&buf, "INFO", data)

		output := buf.String()
		if !strings.Contains(output, "INFO") {
			t.Errorf("日志应包含级别信息")
		}
		if !strings.Contains(output, "testuser") {
			t.Errorf("日志应包含用户信息")
		}
	})

	t.Run("测试map日志", func(t *testing.T) {
		var buf bytes.Buffer

		data := map[string]interface{}{
			"error_code": 500,
			"message":    "内部错误",
		}

		writePlainValue(&buf, "ERROR", data, "extra_field")

		output := buf.String()
		if !strings.Contains(output, "ERROR") {
			t.Errorf("日志应包含错误级别")
		}
		if !strings.Contains(output, "500") {
			t.Errorf("日志应包含错误码")
		}
		if !strings.Contains(output, "extra_field") {
			t.Errorf("日志应包含额外字段")
		}
	})

	t.Run("测试简单值日志", func(t *testing.T) {
		var buf bytes.Buffer
		writePlainValue(&buf, "DEBUG", "简单字符串值")

		output := buf.String()
		if !strings.Contains(output, "DEBUG") {
			t.Errorf("日志应包含调试级别")
		}
	})
}
