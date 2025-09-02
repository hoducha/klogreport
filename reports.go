package main

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

func generateProjectReport(projects []Project) {
	printSectionHeader("Project Time Spent")
	var projectTimes []TimeData
	maxMins := 0

	for _, project := range projects {
		totalMins := 0
		for _, record := range project.Data.Records {
			totalMins += record.TotalMins
		}
		if totalMins > 0 {
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

func generateTagsReport(projects []Project, tagColorMap map[string]*color.Color) {
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

	for _, tt := range tagTimes {
		printBar(tt.Label, tt.Mins, maxMins, tagColorMap[tt.Label], maxLabelWidth)
	}
}

func generateTagsPerProjectReport(projects []Project, tagColorMap map[string]*color.Color) {
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
			printBar("  "+tt.Label, tt.Mins, maxMins, tagColorMap[tt.Label], maxLabelWidth)
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
		return dayTimes[i].Label < dayTimes[j].Label
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