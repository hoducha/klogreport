package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Bar chart rendering
func printBar(label string, mins int, maxMins int, barColor *color.Color, maxLabelWidth int) {
	maxBarLength := 30
	barLength := 0
	if maxMins > 0 {
		barLength = (mins * maxBarLength) / maxMins
	}
	if barLength == 0 && mins > 0 {
		barLength = 1
	}
	
	bar := strings.Repeat("█", barLength)
	duration := formatDuration(mins)
	
	fmt.Printf("  %-*s ", maxLabelWidth, label)
	barColor.Print(bar)
	fmt.Printf(" %s\n", color.HiWhiteString(duration))
}

func calculateMaxLabelWidth(labels []string) int {
	maxWidth := 0
	for _, label := range labels {
		if len(label) > maxWidth {
			maxWidth = len(label)
		}
	}
	return maxWidth
}

// Section headers
func printSectionHeader(title string) {
	fmt.Println()
	color.HiCyanString("┌" + strings.Repeat("─", len(title)+2) + "┐")
	color.HiCyan("│ %s │", title)
	color.HiCyanString("└" + strings.Repeat("─", len(title)+2) + "┘")
	fmt.Println()
}

func printSubsectionHeader(title string) {
	fmt.Printf("  %s\n", color.HiYellowString(title))
}

// Report headers and footers
func printHeader() {
	fmt.Println()
	color.HiCyan("═══════════════════════════════════════════════════════════")
	color.HiCyan("                     KLOG TIME REPORT                      ")  
	color.HiCyan("═══════════════════════════════════════════════════════════")
}

func printFooter(projectCount int) {
	fmt.Println()
	color.HiBlack("─────────────────────────────────────────────────────────────")
	color.HiBlack("%d project(s) processed", projectCount)
	color.HiBlack("─────────────────────────────────────────────────────────────")
}