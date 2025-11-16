package logx

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
)

func TestAddGlobalFields(t *testing.T) {
	globalFields = atomic.Value{}

	fmt.Println("=== 调试 AddGlobalFields ===")

	// 检查初始状态
	initial := globalFields.Load()
	fmt.Printf("初始状态: %v (类型: %T)\n", initial, initial)

	// 添加字段
	fields := []LogField{
		{Key: "service", Value: "api"},
		{Key: "version", Value: "1.0"},
	}

	fmt.Printf("准备添加 %d 个字段\n", len(fields))
	AddGlobalFields(fields...)

	// 检查结果
	result := globalFields.Load()
	fmt.Printf("添加后结果: %v (类型: %T)\n", result, result)

	if result != nil {
		if storedFields, ok := result.([]LogField); ok {
			fmt.Printf("成功获取 %d 个字段:\n", len(storedFields))
			for i, field := range storedFields {
				fmt.Printf("  [%d] Key: %s, Value: %v\n", i, field.Key, field.Value)
			}
		} else {
			fmt.Printf("类型断言失败! 实际类型: %T\n", result)
		}
	}
}

func TestContextWithFields(t *testing.T) {
	t.Run("空上下文添加字段", func(t *testing.T) {
		fmt.Println("=== 开始测试：空上下文添加字段 ===")

		ctx := context.Background()
		fmt.Printf("初始上下文: %v\n", ctx)

		fields := []LogField{
			{Key: "request_id", Value: "abc123"},
			{Key: "user_id", Value: 123},
		}
		fmt.Printf("要添加的字段: %+v\n", fields)

		newCtx := ContextWithFields(ctx, fields...)
		fmt.Printf("添加字段后的新上下文: %v\n", newCtx)

		// 验证字段已添加
		val := newCtx.Value(fieldsKey{})
		fmt.Printf("从上下文获取的值: %v (类型: %T)\n", val, val)

		if val == nil {
			t.Fatal("上下文应该包含字段")
		}

		storedFields, ok := val.([]LogField)
		if !ok {
			t.Fatalf("字段类型应该是 []LogField，实际是 %T", val)
		}

		fmt.Printf("转换后的字段切片: %+v\n", storedFields)
		fmt.Printf("字段数量: %d\n", len(storedFields))

		if len(storedFields) != 2 {
			t.Errorf("期望2个字段，实际得到%d个", len(storedFields))
		}

		// 验证字段内容
		fmt.Println("=== 验证字段内容 ===")
		for i, field := range storedFields {
			fmt.Printf("字段[%d]: Key=%s, Value=%v\n", i, field.Key, field.Value)
		}

		if storedFields[0].Key != "request_id" || storedFields[0].Value != "abc123" {
			t.Errorf("第一个字段不正确: 期望 Key=request_id, Value=abc123, 实际 Key=%s, Value=%v",
				storedFields[0].Key, storedFields[0].Value)
		}

		if storedFields[1].Key != "user_id" || storedFields[1].Value != 123 {
			t.Errorf("第二个字段不正确: 期望 Key=user_id, Value=123, 实际 Key=%s, Value=%v",
				storedFields[1].Key, storedFields[1].Value)
		}
	})

	t.Run("已有字段的上下文添加新字段", func(t *testing.T) {
		// 创建已有字段的上下文
		ctx := context.Background()
		ctx = ContextWithFields(ctx,
			LogField{Key: "service", Value: "api"},
			LogField{Key: "version", Value: "1.0"},
		)

		// 添加新字段
		newCtx := ContextWithFields(ctx,
			LogField{Key: "request_id", Value: "req123"},
			LogField{Key: "user_id", Value: 456},
		)

		// 验证字段合并
		val := newCtx.Value(fieldsKey{})
		storedFields := val.([]LogField)

		if len(storedFields) != 4 {
			t.Errorf("期望4个字段，实际得到%d个", len(storedFields))
		}

		// 验证字段顺序和内容
		expected := []struct {
			key   string
			value interface{}
		}{
			{"service", "api"},
			{"version", "1.0"},
			{"request_id", "req123"},
			{"user_id", 456},
		}

		for i, exp := range expected {
			if storedFields[i].Key != exp.key {
				t.Errorf("字段[%d] key 错误: 期望 %s, 实际 %s", i, exp.key, storedFields[i].Key)
			}
			if storedFields[i].Value != exp.value {
				t.Errorf("字段[%d] value 错误: 期望 %v, 实际 %v", i, exp.value, storedFields[i].Value)
			}
		}
	})

	t.Run("上下文隔离性", func(t *testing.T) {
		baseCtx := context.Background()

		// 从同一个上下文创建两个分支
		ctx1 := ContextWithFields(baseCtx,
			LogField{Key: "branch", Value: "A"},
		)
		ctx2 := ContextWithFields(baseCtx,
			LogField{Key: "branch", Value: "B"},
		)

		// 验证两个上下文互不影响
		fields1 := ctx1.Value(fieldsKey{}).([]LogField)
		fields2 := ctx2.Value(fieldsKey{}).([]LogField)

		if len(fields1) != 1 || len(fields2) != 1 {
			t.Error("每个上下文应该只有1个字段")
		}

		if fields1[0].Value != "A" {
			t.Error("ctx1 应该包含 branch=A")
		}
		if fields2[0].Value != "B" {
			t.Error("ctx2 应该包含 branch=B")
		}

		// 验证基础上下文未被修改
		if baseCtx.Value(fieldsKey{}) != nil {
			t.Error("基础上下文不应该被修改")
		}
	})
}
