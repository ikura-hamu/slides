package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"slices"
	"strings"
)

type slideInfo struct {
	Name        string
	Title       string
	Description string
	Date        string
}

func main() {
	entries, err := os.ReadDir("../src")
	if err != nil {
		panic(fmt.Sprintf("failed to read directory: %v", err))
	}

	slides := make([]slideInfo, 0, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		slides = append(slides, loadFile(e.Name()))
	}

	slices.SortFunc[[]slideInfo](slides, func(a, b slideInfo) int {
		if a.Date > b.Date {
			return -1
		} else if a.Date == b.Date {
			return 0
		} else {
			return 1
		}
	})

	t := template.Must(template.ParseFiles("template.html"))

	f, err := os.Create("../docs/index.html")
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %v", err))
	}
	defer f.Close()

	err = t.Execute(f, slides)
	if err != nil {
		panic(fmt.Sprintf("failed to execute template: %v", err))
	}
}

func loadFile(fileName string) slideInfo {
	f, err := os.Open("../src/" + fileName)
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %v", err))
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("failed to read file content: %v", err))
	}

	var title string
	var description string
	var date string
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "title: ") {
			title = strings.TrimSpace(strings.ReplaceAll(l, "title: ", ""))
		}
		if strings.HasPrefix(l, "description: ") {
			description = strings.TrimSpace(strings.ReplaceAll(l, "description: ", ""))
		}
		if strings.HasPrefix(l, "date: ") {
			date = strings.TrimSpace(strings.ReplaceAll(l, "date: ", ""))
		}

		if title != "" && description != "" && date != "" {
			break
		}
	}

	return slideInfo{
		Name:        strings.TrimRight(fileName, ".md"),
		Title:       title,
		Description: description,
		Date:        date,
	}
}
