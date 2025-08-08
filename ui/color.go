package ui

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	ansiReset = "\u001b[0m"

	ansiBold   = "\u001b[1m"
	ansiDim    = "\u001b[2m"
	ansiRed    = "\u001b[31m"
	ansiGreen  = "\u001b[32m"
	ansiYellow = "\u001b[33m"
	ansiBlue   = "\u001b[34m"
	ansiMagenta= "\u001b[35m"
	ansiCyan   = "\u001b[36m"
)

var enableColor = detectColor()

func detectColor() bool {
	// Respect NO_COLOR convention
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	return true
}

func wrap(code, s string) string {
	if !enableColor {
		return s
	}
	return code + s + ansiReset
}

func Bold(s string) string   { return wrap(ansiBold, s) }
func Dim(s string) string    { return wrap(ansiDim, s) }
func Red(s string) string    { return wrap(ansiRed, s) }
func Green(s string) string  { return wrap(ansiGreen, s) }
func Yellow(s string) string { return wrap(ansiYellow, s) }
func Cyan(s string) string   { return wrap(ansiCyan, s) }

// Convenience helpers for formatting strings with placeholders
// Example: FmtGreen("Error: %s\n", err)
func FmtRed(format string, a ...any) string    { return wrap(ansiRed, sprintf(format, a...)) }
func FmtGreen(format string, a ...any) string  { return wrap(ansiGreen, sprintf(format, a...)) }
func FmtYellow(format string, a ...any) string { return wrap(ansiYellow, sprintf(format, a...)) }
func FmtCyan(format string, a ...any) string   { return wrap(ansiCyan, sprintf(format, a...)) }
func FmtBold(format string, a ...any) string   { return wrap(ansiBold, sprintf(format, a...)) }
func FmtDim(format string, a ...any) string    { return wrap(ansiDim, sprintf(format, a...)) }

// local sprintf to avoid importing fmt in user code when using Fmt* helpers here
func sprintf(format string, a ...any) string {
	return fmtSprintf(format, a...)
}

// Minimal indirection to keep fmt usage internal to this package
var fmtSprintf = func(format string, a ...any) string { return fmt.Sprintf(format, a...) }

