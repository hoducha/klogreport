package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func loadProjects(logDir string, args []string) []Project {
	var allProjects []Project

	filepath.WalkDir(logDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".klg") {
			return nil
		}

		projectName := strings.TrimSuffix(d.Name(), ".klg")

		// Run `klog json` command to get the JSON data
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

	return allProjects
}


func main() {
	logDir := os.Getenv("KLOG_DIR")
	if logDir == "" {
		home, _ := os.UserHomeDir()
		logDir = filepath.Join(home, "klog")
	}

	args := os.Args[1:] // arguments to be passed to `klog json` command
	allProjects := loadProjects(logDir, args)

	if len(allProjects) == 0 {
		color.HiRed("❌ No .klg files found in %s", logDir)
		return
	}

	printHeader()
	generateProjectReport(allProjects)
	generateTagsReport(allProjects)
	generateTagsPerProjectReport(allProjects)
	generateDailyReport(allProjects)
	printFooter(len(allProjects))
}
