package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func getTagColorMap(projects []Project) map[string]*color.Color {
	tagTotals := make(map[string]int)
	
	for _, project := range projects {
		for _, record := range project.Data.Records {
			for _, entry := range record.Entries {
				for _, tag := range entry.Tags {
					tagTotals[tag] += entry.TotalMins
				}
			}
		}
	}
	
	var tags []string
	for tag := range tagTotals {
		tags = append(tags, tag)
	}
	
	sort.Strings(tags)
	
	const selectedPalette = "tableau10"
	palette := colorPalettes[selectedPalette]
	paletteSize := len(palette.Colors)
	colorMap := make(map[string]*color.Color)
	for i, tag := range tags {
		if i < paletteSize {
			colorMap[tag] = palette.Colors[i]
		} else {
			colorMap[tag] = generateOverflowColor(i, paletteSize)
		}
	}
	
	return colorMap
}

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