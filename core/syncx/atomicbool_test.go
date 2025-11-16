package syncx

import (
	"fmt"
	"testing"
)

// 调试用：打印实际的地址值
func TestDebug_PrintAddresses(t *testing.T) {
	fmt.Println("=== 调试：查看 new() 返回的地址 ===")

	b1 := NewAtomicBool()
	fmt.Printf("b1 地址: %p\n", b1)

	b2 := NewAtomicBool()
	fmt.Printf("b2 地址: %p\n", b2)

	b3 := ForAtomicBool(true)
	fmt.Printf("b3 地址: %p\n", b3)

	b4 := ForAtomicBool(false)
	fmt.Printf("b4 地址: %p\n", b4)
}
