package cli

import (
	"fmt"
	"strings"
)

const ansiEsc = "\x1b"

type TextStyle int

const (
	Black TextStyle = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	_
	Default
)

const (
	Reset TextStyle = iota
	Bold
	Dim
	Italic
	Underline
	Blinking
	_
	Inverse
	Hidden
	Strikethrough
)

func fg(c TextStyle) int {
	return int(c) + 30
}

func Style(s string, styles ...TextStyle) string {
	if len(styles) == 0 {
		return s
	}
	var styleStrs []string
	for _, style := range styles {
		styleStrs = append(styleStrs, fmt.Sprintf("%d", style))
	}
	ansiCode := strings.Join(styleStrs, ";")
	return fmt.Sprintf("%s[%sm%s%s[39m", ansiEsc, ansiCode, s, ansiEsc)
}
