package cmd

import (
	"os"

	"golang.org/x/term"
)

const (
	defaultTermWidth  = 120
	defaultTermHeight = 30
)

func terminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 || h <= 0 {
		return defaultTermWidth, defaultTermHeight
	}
	return w, h
}