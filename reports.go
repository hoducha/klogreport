package main

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

func generateProjectReport(projects []Project) {
	printSectionHeader("Project Time Spent")
	
	// Collect project data with tag breakdown
	type ProjectInfo struct {
		Name     string
		Total    int
		Segments []TagSegment
	}
	
	var projectInfos []ProjectInfo
	maxMins := 0

	tagColorMap := getTagColorMap(projects)

	for _, project := range projects {
		tagTotals := make(map[string]int)
		totalMins := 0
		
		for _, record := range project.Data.Records {
			totalMins += record.TotalMins
			for _, entry := range record.Entries {
				for _, tag := range entry.Tags {
					tagTotals[tag] += entry.TotalMins
				}
			}
		}
		
		if totalMins > 0 {
			// Create segments sorted by time (largest first)
			var segments []TagSegment
			for tag, mins := range tagTotals {
				segments = append(segments, TagSegment{
					Tag:   tag,
					Mins:  mins,
					Color: tagColorMap[tag],
				})
			}
			
			sort.Slice(segments, func(i, j int) bool {
				return segments[i].Mins > segments[j].Mins
			})
			
			projectInfos = append(projectInfos, ProjectInfo{
				Name:     project.Name,
				Total:    totalMins,
				Segments: segments,
			})
			
			if totalMins > maxMins {
				maxMins = totalMins
			}
		}
	}

	if len(projectInfos) == 0 {
		fmt.Printf("  %s\n", color.HiBlackString("No time logged yet"))
		return
	}

	// Sort projects by total time
	sort.Slice(projectInfos, func(i, j int) bool {
		return projectInfos[i].Total > projectInfos[j].Total
	})

	// Calculate max label width
	var labels []string
	for _, pi := range projectInfos {
		labels = append(labels, pi.Name)
	}
	maxLabelWidth := calculateMaxLabelWidth(labels)

	// Print segmented bars
	for _, pi := range projectInfos {
		printSegmentedBar(pi.Name, pi.Segments, pi.Total, maxMins, maxLabelWidth)
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
	
	// Aggregate by day with tag breakdown
	type DayInfo struct {
		Date     string
		Total    int
		Segments []TagSegment
	}
	
	dayTagTotals := make(map[string]map[string]int) // date -> tag -> mins
	dayTotals := make(map[string]int)               // date -> total mins
	maxMins := 0

	// Get tag color mapping
	tagColorMap := getTagColorMap(projects)

	for _, project := range projects {
		for _, record := range project.Data.Records {
			if dayTagTotals[record.Date] == nil {
				dayTagTotals[record.Date] = make(map[string]int)
			}
			
			dayTotals[record.Date] += record.TotalMins
			
			for _, entry := range record.Entries {
				for _, tag := range entry.Tags {
					dayTagTotals[record.Date][tag] += entry.TotalMins
				}
			}
			
			if dayTotals[record.Date] > maxMins {
				maxMins = dayTotals[record.Date]
			}
		}
	}

	if len(dayTotals) == 0 {
		fmt.Printf("  %s\n", color.HiBlackString("No daily data found"))
		return
	}

	var dayInfos []DayInfo
	for date, total := range dayTotals {
		var segments []TagSegment
		for tag, mins := range dayTagTotals[date] {
			segments = append(segments, TagSegment{
				Tag:   tag,
				Mins:  mins,
				Color: tagColorMap[tag],
			})
		}
		
		// Sort segments by time (largest first)
		sort.Slice(segments, func(i, j int) bool {
			return segments[i].Mins > segments[j].Mins
		})
		
		dayInfos = append(dayInfos, DayInfo{
			Date:     date,
			Total:    total,
			Segments: segments,
		})
	}

	// Sort by date
	sort.Slice(dayInfos, func(i, j int) bool {
		return dayInfos[i].Date < dayInfos[j].Date
	})

	// Calculate max label width
	var labels []string
	for _, di := range dayInfos {
		labels = append(labels, di.Date)
	}
	maxLabelWidth := calculateMaxLabelWidth(labels)

	// Print segmented bars
	for _, di := range dayInfos {
		printSegmentedBar(di.Date, di.Segments, di.Total, maxMins, maxLabelWidth)
	}
}