package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

type Entry struct {
	Summary   string   `json:"summary"`
	Total     string   `json:"total"`
	TotalMins int      `json:"total_mins"`
	Tags      []string `json:"tags"`
}

type Record struct {
	Date      string  `json:"date"`
	Entries   []Entry `json:"entries"`
	TotalMins int     `json:"total_mins"`
}

type ProjectData struct {
	Records []Record `json:"records"`
}

type Project struct {
	Name string
	Data ProjectData
}

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

// printBar prints a colorful text-based bar chart with dynamic scaling
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

type TimeData struct {
	Label string
	Mins  int
}

func generateProjectReport(projects []Project) {
	printSectionHeader("Project Time Spent")
	var projectTimes []TimeData
	maxMins := 0

	for _, project := range projects {
		totalMins := 0
		for _, record := range project.Data.Records {
			totalMins += record.TotalMins
		}
		if totalMins > 0 { // Only include projects with time
			projectTimes = append(projectTimes, TimeData{project.Name, totalMins})
			if totalMins > maxMins {
				maxMins = totalMins
			}
		}
	}

	if len(projectTimes) == 0 {
		fmt.Printf("  %s\n", color.HiBlackString("No time logged yet"))
		return
	}

	sort.Slice(projectTimes, func(i, j int) bool {
		return projectTimes[i].Mins > projectTimes[j].Mins
	})

	var labels []string
	for _, pt := range projectTimes {
		labels = append(labels, pt.Label)
	}
	maxLabelWidth := calculateMaxLabelWidth(labels)

	projectColor := color.New(color.FgHiBlue)
	for _, pt := range projectTimes {
		printBar(pt.Label, pt.Mins, maxMins, projectColor, maxLabelWidth)
	}
}

func generateTagsReport(projects []Project) {
	printSectionHeader("Tags Time Spent")
	tagTotals := make(map[string]int)
	maxMins := 0

	for _, project := range projects {
		for _, record := range project.Data.Records {
			for _, entry := range record.Entries {
				for _, tag := range entry.Tags {
					tagTotals[tag] += entry.TotalMins
					if tagTotals[tag] > maxMins {
						maxMins = tagTotals[tag]
					}
				}
			}
		}
	}

	if len(tagTotals) == 0 {
		fmt.Printf("  %s\n", color.HiBlackString("No tags found"))
		return
	}

	var tagTimes []TimeData
	for tag, mins := range tagTotals {
		tagTimes = append(tagTimes, TimeData{tag, mins})
	}

	sort.Slice(tagTimes, func(i, j int) bool {
		return tagTimes[i].Mins > tagTimes[j].Mins
	})

	var labels []string
	for _, tt := range tagTimes {
		labels = append(labels, tt.Label)
	}
	maxLabelWidth := calculateMaxLabelWidth(labels)

	tagColor := color.New(color.FgHiGreen)
	for _, tt := range tagTimes {
		printBar(tt.Label, tt.Mins, maxMins, tagColor, maxLabelWidth)
	}
}

func generateTagsPerProjectReport(projects []Project) {
	printSectionHeader("Tags per Project")
	
	// Filter and sort projects by total time
	var activeProjects []Project
	for _, project := range projects {
		totalMins := 0
		tagTotals := make(map[string]int)
		
		for _, record := range project.Data.Records {
			totalMins += record.TotalMins
			for _, entry := range record.Entries {
				for _, tag := range entry.Tags {
					tagTotals[tag] += entry.TotalMins
				}
			}
		}
		
		if len(tagTotals) > 0 {
			activeProjects = append(activeProjects, project)
		}
	}
	
	if len(activeProjects) == 0 {
		fmt.Printf("  %s\n", color.HiBlackString("No tagged time found"))
		return
	}

	sort.Slice(activeProjects, func(i, j int) bool {
		totalI := 0
		totalJ := 0
		for _, record := range activeProjects[i].Data.Records {
			totalI += record.TotalMins
		}
		for _, record := range activeProjects[j].Data.Records {
			totalJ += record.TotalMins
		}
		return totalI > totalJ
	})

	tagColor := color.New(color.FgHiMagenta)
	
	for _, project := range activeProjects {
		tagTotals := make(map[string]int)
		maxMins := 0

		for _, record := range project.Data.Records {
			for _, entry := range record.Entries {
				for _, tag := range entry.Tags {
					tagTotals[tag] += entry.TotalMins
					if tagTotals[tag] > maxMins {
						maxMins = tagTotals[tag]
					}
				}
			}
		}

		printSubsectionHeader(project.Name)

		var tagTimes []TimeData
		for tag, mins := range tagTotals {
			tagTimes = append(tagTimes, TimeData{tag, mins})
		}

		sort.Slice(tagTimes, func(i, j int) bool {
			return tagTimes[i].Mins > tagTimes[j].Mins
		})

		var labels []string
		for _, tt := range tagTimes {
			labels = append(labels, "  "+tt.Label)
		}
		maxLabelWidth := calculateMaxLabelWidth(labels)

		for _, tt := range tagTimes {
			printBar("  "+tt.Label, tt.Mins, maxMins, tagColor, maxLabelWidth)
		}
		fmt.Println()
	}
}

func generateDailyReport(projects []Project) {
	printSectionHeader("Daily Working Time")
	dayTotals := make(map[string]int)
	maxMins := 0

	for _, project := range projects {
		for _, record := range project.Data.Records {
			dayTotals[record.Date] += record.TotalMins
			if dayTotals[record.Date] > maxMins {
				maxMins = dayTotals[record.Date]
			}
		}
	}

	if len(dayTotals) == 0 {
		fmt.Printf("  %s\n", color.HiBlackString("No daily data found"))
		return
	}

	var dayTimes []TimeData
	for date, mins := range dayTotals {
		dayTimes = append(dayTimes, TimeData{date, mins})
	}

	sort.Slice(dayTimes, func(i, j int) bool {
		return dayTimes[i].Label < dayTimes[j].Label // Sort by date
	})

	var labels []string
	for _, dt := range dayTimes {
		labels = append(labels, dt.Label)
	}
	maxLabelWidth := calculateMaxLabelWidth(labels)

	dayColor := color.New(color.FgHiRed)
	for _, dt := range dayTimes {
		printBar(dt.Label, dt.Mins, maxMins, dayColor, maxLabelWidth)
	}
}

func main() {
	// Get log directory from env or fallback
	logDir := os.Getenv("KLOG_DIR")
	if logDir == "" {
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, "klog")
	}

	args := os.Args[1:] // additional args for klog json

	var allProjects []Project

	// iterate over *.klg files
	filepath.WalkDir(logDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".klg") {
			return nil
		}

		projectName := strings.TrimSuffix(d.Name(), ".klg")

		// run klog json <file> <args>
		cmdArgs := append([]string{"json", path}, args...)
		out, err := exec.Command("klog", cmdArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running klog json on %s: %v\n", path, err)
			return nil
		}

		var data ProjectData
		if err := json.Unmarshal(out, &data); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON for %s: %v\n", path, err)
			return nil
		}

		allProjects = append(allProjects, Project{
			Name: projectName,
			Data: data,
		})

		return nil
	})

	if len(allProjects) == 0 {
		color.HiRed("❌ No .klg files found in %s", logDir)
		return
	}

	// Print header
	fmt.Println()
	color.HiCyan("═══════════════════════════════════════════════════════════")
	color.HiCyan("                     KLOG TIME REPORT                      ")  
	color.HiCyan("═══════════════════════════════════════════════════════════")

	// Generate all reports
	generateProjectReport(allProjects)
	generateTagsReport(allProjects)
	generateTagsPerProjectReport(allProjects)
	generateDailyReport(allProjects)

	// Print footer
	fmt.Println()
	color.HiBlack("─────────────────────────────────────────────────────────────")
	totalProjects := len(allProjects)
	color.HiBlack("%d project(s) processed", totalProjects)
	color.HiBlack("─────────────────────────────────────────────────────────────")
}
