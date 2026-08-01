package cmd

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func blockWidth(s string) int {
	width := 0
	for _, line := range strings.Split(s, "\n") {
		if w := visibleWidth(line); w > width {
			width = w
		}
	}
	return width
}

func visibleWidth(s string) int {
	return utf8.RuneCountInString(ansiRegexp.ReplaceAllString(s, ""))
}

func padVisible(s string, width int) string {
	pad := width - visibleWidth(s)
	if pad <= 0 {
		return s
	}
	return s + strings.Repeat(" ", pad)
}

func sideBySide(left, right string, gap int) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
	leftWidth := blockWidth(left)

	n := len(leftLines)
	if len(rightLines) > n {
		n = len(rightLines)
	}

	var b strings.Builder
	for i := 0; i < n; i++ {
		var l, r string
		if i < len(leftLines) {
			l = leftLines[i]
		}
		if i < len(rightLines) {
			r = rightLines[i]
		}
		b.WriteString(padVisible(l, leftWidth))
		b.WriteString(strings.Repeat(" ", gap))
		b.WriteString(r)
		b.WriteString("\n")
	}
	return b.String()
}