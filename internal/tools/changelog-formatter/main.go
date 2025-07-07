package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

type sectionType string

const (
	sectionTypeFeatures     sectionType = "FEATURES:"
	sectionTypeEnhancements sectionType = "ENHANCEMENTS:"
	sectionTypeBugs         sectionType = "BUG FIXES:"
	sectionTypeUnknown      sectionType = "UNKNOWN"
)

type changelog struct {
	pre  []string
	f    features
	e    enhancements
	b    bugs
	post []string
}

type features struct {
	general []string
}

type enhancements struct {
	dependencies []string
	dataSources  []string
	general      []string
}

type bugs struct {
	dataSources []string
	general     []string
}

var (
	sectionHeadingRegex = regexp.MustCompile(`(FEATURES.*|ENHANCEMENTS.*|BUG FIXES.*)`)
	versionHeadingRegex = regexp.MustCompile(`## \d*\.\d*\.\d* \(\w* \d{1,2}, \d{4}\)`)
)

func formatChangelog(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := strings.Split(string(bytes), "\n")

	cl := &changelog{}
	st := sectionTypeUnknown
	parsingNewEntries := false
	for idx, line := range content {
		if versionHeadingRegex.MatchString(line) {
			cl.post = content[idx:]
			break
		}

		if sectionHeadingRegex.MatchString(line) {
			parsingNewEntries = true

			st = determineSectionType(line)
			continue
		}

		if line != "" && parsingNewEntries {
			addChangelogEntry(cl, st, line)
		}

		if !parsingNewEntries {
			cl.pre = append(cl.pre, line)
		}
	}
	newContent := rebuildChangelog(cl)

	if err := os.WriteFile(path, []byte(strings.Join(newContent, "\n")), 0o644); err != nil {
		return fmt.Errorf("writing to `%s`", path)
	}

	return nil
}

func determineSectionType(line string) sectionType {
	switch {
	case strings.HasPrefix(line, "FEATURES"):
		return sectionTypeFeatures
	case strings.HasPrefix(line, "ENHANCEMENTS"):
		return sectionTypeEnhancements
	case strings.HasPrefix(line, "BUG FIXES"):
		return sectionTypeBugs
	default:
		return sectionTypeUnknown
	}
}

func addChangelogEntry(cl *changelog, st sectionType, line string) {
	if line == "" {
		return
	}

	switch st {
	case sectionTypeFeatures:
		cl.f.general = append(cl.f.general, line)
	case sectionTypeEnhancements:
		switch {
		case strings.Contains(line, "dependencies"):
			cl.e.dependencies = append(cl.e.dependencies, line)
		case strings.Contains(line, "Data Source"):
			cl.e.dataSources = append(cl.e.dataSources, line)
		default:
			cl.e.general = append(cl.e.general, line)
		}
	case sectionTypeBugs:
		if strings.Contains(line, "Data Source") {
			cl.b.dataSources = append(cl.b.dataSources, line)
		} else {
			cl.b.general = append(cl.b.general, line)
		}
	default:
		return
	}
}

func formatSection(entries []string, st sectionType) []string {
	result := []string{
		string(st),
		"",
		"",
	}
	return slices.Insert(result, 2, entries...)
}

func rebuildChangelog(cl *changelog) []string {
	newContent := make([]string, 0)
	newContent = append(newContent, cl.pre...)

	sort := func(s []string) {
		slices.SortFunc(s, func(a, b string) int {
			resourceName := func(s string) string {
				start := strings.Index(s, "`")
				if start == -1 {
					return s
				}

				end := strings.Index(s[start+1:], "`")
				if end == -1 {
					return s[start+1:]
				}

				return s[start+1 : start+1+end]
			}

			return strings.Compare(resourceName(a), resourceName(b))
		})
	}

	if len(cl.f.general) > 0 {
		sort(cl.f.general)
		newContent = append(newContent, formatSection(cl.f.general, sectionTypeFeatures)...)
	}

	tmpContent := make([]string, 0)
	if len(cl.e.dependencies) > 0 {
		sort(cl.e.dependencies)
		tmpContent = append(tmpContent, cl.e.dependencies...)
	}

	if len(cl.e.dataSources) > 0 {
		sort(cl.e.dataSources)
		tmpContent = append(tmpContent, cl.e.dataSources...)
	}

	if len(cl.e.general) > 0 {
		sort(cl.e.general)
		tmpContent = append(tmpContent, cl.e.general...)
	}

	if len(tmpContent) > 0 {
		newContent = append(newContent, formatSection(tmpContent, sectionTypeEnhancements)...)
	}

	tmpContent = make([]string, 0)
	if len(cl.b.dataSources) > 0 {
		sort(cl.b.dataSources)
		tmpContent = append(tmpContent, cl.b.dataSources...)
	}

	if len(cl.b.general) > 0 {
		sort(cl.b.general)
		tmpContent = append(tmpContent, cl.b.general...)
	}

	if len(tmpContent) > 0 {
		newContent = append(newContent, formatSection(tmpContent, sectionTypeBugs)...)
	}

	return append(newContent, cl.post...)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <path to changelog>")
		return
	}
	filePath := os.Args[1]

	if err := formatChangelog(filePath); err != nil {
		fmt.Println(fmt.Errorf("formatting changelog: %+v", err))
		return
	}
}
