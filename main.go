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

// printBar prints a simple text-based bar chart with dynamic scaling
func printBar(label string, mins int, maxMins int) {
	maxBarLength := 20
	barLength := 0
	if maxMins > 0 {
		barLength = (mins * maxBarLength) / maxMins
	}
	if barLength == 0 && mins > 0 {
		barLength = 1
	}
	fmt.Printf("%-12s %s %s\n", label, strings.Repeat("█", barLength), formatDuration(mins))
}

type TimeData struct {
	Label string
	Mins  int
}

// generateProjectReport generates "Total time by project" report
func generateProjectReport(projects []Project) {
	fmt.Println("=== Total Time by Project ===")
	var projectTimes []TimeData
	maxMins := 0

	for _, project := range projects {
		totalMins := 0
		for _, record := range project.Data.Records {
			totalMins += record.TotalMins
		}
		projectTimes = append(projectTimes, TimeData{project.Name, totalMins})
		if totalMins > maxMins {
			maxMins = totalMins
		}
	}

	sort.Slice(projectTimes, func(i, j int) bool {
		return projectTimes[i].Mins > projectTimes[j].Mins
	})

	for _, pt := range projectTimes {
		printBar(pt.Label, pt.Mins, maxMins)
	}
	fmt.Println()
}

// generateTagsReport generates "Total time by tags" report
func generateTagsReport(projects []Project) {
	fmt.Println("=== Total Time by Tags ===")
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

	var tagTimes []TimeData
	for tag, mins := range tagTotals {
		tagTimes = append(tagTimes, TimeData{tag, mins})
	}

	sort.Slice(tagTimes, func(i, j int) bool {
		return tagTimes[i].Mins > tagTimes[j].Mins
	})

	for _, tt := range tagTimes {
		printBar(tt.Label, tt.Mins, maxMins)
	}
	fmt.Println()
}

// generateTagsPerProjectReport generates "Total time by tags per project" report
func generateTagsPerProjectReport(projects []Project) {
	fmt.Println("=== Total Time by Tags per Project ===")
	
	sort.Slice(projects, func(i, j int) bool {
		totalI := 0
		totalJ := 0
		for _, record := range projects[i].Data.Records {
			totalI += record.TotalMins
		}
		for _, record := range projects[j].Data.Records {
			totalJ += record.TotalMins
		}
		return totalI > totalJ
	})

	for _, project := range projects {
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

		if len(tagTotals) == 0 {
			continue
		}

		fmt.Println(project.Name)

		var tagTimes []TimeData
		for tag, mins := range tagTotals {
			tagTimes = append(tagTimes, TimeData{tag, mins})
		}

		sort.Slice(tagTimes, func(i, j int) bool {
			return tagTimes[i].Mins > tagTimes[j].Mins
		})

		for _, tt := range tagTimes {
			printBar("  "+tt.Label, tt.Mins, maxMins)
		}
		fmt.Println()
	}
}

// generateDailyReport generates "Total time by days" report
func generateDailyReport(projects []Project) {
	fmt.Println("=== Total Time by Days ===")
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

	var dayTimes []TimeData
	for date, mins := range dayTotals {
		dayTimes = append(dayTimes, TimeData{date, mins})
	}

	sort.Slice(dayTimes, func(i, j int) bool {
		return dayTimes[i].Label < dayTimes[j].Label // Sort by date
	})

	for _, dt := range dayTimes {
		printBar(dt.Label, dt.Mins, maxMins)
	}
	fmt.Println()
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
		fmt.Println("No .klg files found in", logDir)
		return
	}

	// Generate all reports
	generateProjectReport(allProjects)
	generateTagsReport(allProjects)
	generateTagsPerProjectReport(allProjects)
	generateDailyReport(allProjects)
}
