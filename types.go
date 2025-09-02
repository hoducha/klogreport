package main

// Entry represents a single time entry in a klog record
type Entry struct {
	Summary   string   `json:"summary"`
	Total     string   `json:"total"`
	TotalMins int      `json:"total_mins"`
	Tags      []string `json:"tags"`
}

// Record represents a daily record containing multiple entries
type Record struct {
	Date      string  `json:"date"`
	Entries   []Entry `json:"entries"`
	TotalMins int     `json:"total_mins"`
}

// ProjectData holds all records for a project
type ProjectData struct {
	Records []Record `json:"records"`
}

// Project represents a complete project with its data
type Project struct {
	Name string
	Data ProjectData
}

// TimeData is used for sorting and displaying time information
type TimeData struct {
	Label string
	Mins  int
}