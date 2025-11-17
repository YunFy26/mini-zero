package color

import "github.com/fatih/color"

const (
	NoColor Color = iota
	FgBlack
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	BgBlack
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

var colors = map[Color][]color.Attribute{
	FgBlack:   {color.FgBlack, color.Bold},
	FgRed:     {color.FgRed, color.Bold},
	FgGreen:   {color.FgGreen, color.Bold},
	FgYellow:  {color.FgYellow, color.Bold},
	FgBlue:    {color.FgBlue, color.Bold},
	FgMagenta: {color.FgMagenta, color.Bold},
	FgCyan:    {color.FgCyan, color.Bold},
	FgWhite:   {color.FgWhite, color.Bold},
	BgBlack:   {color.BgBlack, color.FgHiWhite, color.Bold},
	BgRed:     {color.BgRed, color.FgHiWhite, color.Bold},
	BgGreen:   {color.BgGreen, color.FgHiWhite, color.Bold},
	BgYellow:  {color.BgYellow, color.FgHiWhite, color.Bold},
	BgBlue:    {color.BgBlue, color.FgHiWhite, color.Bold},
	BgMagenta: {color.BgMagenta, color.FgHiWhite, color.Bold},
	BgCyan:    {color.BgCyan, color.FgHiWhite, color.Bold},
	BgWhite:   {color.BgWhite, color.FgBlack, color.Bold},
}

type Color uint32

func WithColor(text string, colour Color) string {
	c := color.New(colors[colour]...)
	return c.Sprint(text)
}

func WithColorPadding(text string, colour Color) string {
	return WithColor(" "+text+" ", colour)
}
