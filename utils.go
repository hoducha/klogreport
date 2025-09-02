package main

import (
	"fmt"
	"strings"
)

// parseDuration converts "2h15m" -> minutes
func parseDuration(d string) int {
	h, m := 0, 0
	if strings.Contains(d, "h") {
		parts := strings.Split(d, "h")
		fmt.Sscanf(parts[0], "%d", &h)
		d = ""
		if len(parts) > 1 {
			d = parts[1]
		}
	}
	if strings.Contains(d, "m") {
		fmt.Sscanf(d, "%dm", &m)
	}
	return h*60 + m
}

// formatDuration converts minutes -> "2h15m"
func formatDuration(mins int) string {
	h := mins / 60
	m := mins % 60
	if h > 0 {
		if m > 0 {
			return fmt.Sprintf("%dh%dm", h, m)
		}
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", m)
}