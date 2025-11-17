package color

import (
	"fmt"
	"testing"
)

func TestWithColor(t *testing.T) {
	fmt.Println("\n🎨 颜色演示测试（应该在终端中看到彩色输出）")
	fmt.Println("=========================================")

	// 前景色演示
	fmt.Println("\n📝 前景色演示:")
	colors := []struct {
		name  string
		color Color
	}{
		{"黑色前景", FgBlack},
		{"红色前景", FgRed},
		{"绿色前景", FgGreen},
		{"黄色前景", FgYellow},
		{"蓝色前景", FgBlue},
		{"洋红前景", FgMagenta},
		{"青色前景", FgCyan},
		{"白色前景", FgWhite},
	}

	for _, c := range colors {
		coloredText := WithColor(c.name, c.color)
		fmt.Printf("  %s\n", coloredText)
	}

	// 背景色演示
	fmt.Println("\n🎨 背景色演示:")
	bgColors := []struct {
		name  string
		color Color
	}{
		{"黑色背景", BgBlack},
		{"红色背景", BgRed},
		{"绿色背景", BgGreen},
		{"黄色背景", BgYellow},
		{"蓝色背景", BgBlue},
		{"洋红背景", BgMagenta},
		{"青色背景", BgCyan},
		{"白色背景", BgWhite},
	}

	for _, bg := range bgColors {
		coloredText := WithColor(bg.name, bg.color)
		fmt.Printf("  %s\n", coloredText)
	}
}

func TestWithColorPadding(t *testing.T) {
	fmt.Println("\n📦 带内边距的颜色演示")
	fmt.Println("====================")

	examples := []struct {
		text  string
		color Color
	}{
		{"错误", BgRed},
		{"警告", BgYellow},
		{"成功", BgGreen},
		{"信息", BgBlue},
		{"调试", BgCyan},
	}

	for _, ex := range examples {
		paddedText := WithColorPadding(ex.text, ex.color)
		fmt.Printf("  %s\n", paddedText)
	}
}
