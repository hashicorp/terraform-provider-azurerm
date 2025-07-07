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
	sectionTypeBreakingChanges sectionType = "BREAKING CHANGES:"
	sectionTypeFeatures        sectionType = "FEATURES:"
	sectionTypeEnhancements    sectionType = "ENHANCEMENTS:"
	sectionTypeBugs            sectionType = "BUG FIXES:"
	sectionTypeUnknown         sectionType = "UNKNOWN"
)

type changelog struct {
	pre          []string
	breaking     []string
	features     features
	enhancements enhancements
	bugs         bugs
	post         []string
}

type features struct {
	dataSources []string
	general     []string
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
	sectionHeadingRegex = regexp.MustCompile(`(BREAKING CHANGES.*|FEATURES.*|ENHANCEMENTS.*|BUG FIXES.*)`)
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

	changeLog := &changelog{}
	st := sectionTypeUnknown
	parsingNewEntries := false
	for idx, line := range content {
		// Once we hit a released version's changelog heading, add the remaining content to `changelog.post`.
		// We won't format any of the lines added here.
		if versionHeadingRegex.MatchString(line) {
			changeLog.post = content[idx:]
			break
		}

		if sectionHeadingRegex.MatchString(line) {
			parsingNewEntries = true

			st = determineSectionType(line)
			continue
		}

		if line != "" && parsingNewEntries {
			addChangelogEntry(changeLog, st, line)
		}

		if !parsingNewEntries {
			// If we're not currently parsing new changelog entries, i.e we haven't encountered a `BREAKING CHANGES`, `FEATURES`, `ENHANCEMENTS`, or `BUG FIXES` heading
			// track it in `changelog.pre`. We won't format any of the lines added here.
			changeLog.pre = append(changeLog.pre, line)
		}
	}
	newContent := rebuildChangelog(changeLog)

	if err := os.WriteFile(path, []byte(strings.Join(newContent, "\n")), 0o644); err != nil {
		return fmt.Errorf("writing to `%s`", path)
	}

	return nil
}

func determineSectionType(line string) sectionType {
	switch {
	case strings.HasPrefix(line, "BREAKING CHANGES"):
		return sectionTypeBreakingChanges
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

func addChangelogEntry(changeLog *changelog, st sectionType, line string) {
	if line == "" {
		return
	}

	switch st {
	case sectionTypeBreakingChanges:
		changeLog.breaking = append(changeLog.breaking, line)
	case sectionTypeFeatures:
		if strings.Contains(line, "Data Source") {
			changeLog.features.dataSources = append(changeLog.features.dataSources, line)
		} else {
			changeLog.features.general = append(changeLog.features.general, line)
		}
	case sectionTypeEnhancements:
		switch {
		case strings.Contains(line, "dependencies"):
			changeLog.enhancements.dependencies = append(changeLog.enhancements.dependencies, line)
		case strings.Contains(line, "Data Source"):
			changeLog.enhancements.dataSources = append(changeLog.enhancements.dataSources, line)
		default:
			changeLog.enhancements.general = append(changeLog.enhancements.general, line)
		}
	case sectionTypeBugs:
		if strings.Contains(line, "Data Source") {
			changeLog.bugs.dataSources = append(changeLog.bugs.dataSources, line)
		} else {
			changeLog.bugs.general = append(changeLog.bugs.general, line)
		}
	default:
		return
	}
}

// formatSection ensures that a changelog section is formatted as expected.
// The empty strings ensure there is always a newline after the heading and the last changelog entry, the entries are inserted between them.
func formatSection(entries []string, st sectionType) []string {
	result := []string{
		string(st),
		"",
		"",
	}
	return slices.Insert(result, 2, entries...)
}

func rebuildChangelog(changeLog *changelog) []string {
	newContent := make([]string, 0)
	newContent = append(newContent, changeLog.pre...)

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

	if len(changeLog.breaking) > 0 {
		fmt.Println("changelog contains breaking changes, please sort the entries manually")
		newContent = append(newContent, formatSection(changeLog.breaking, sectionTypeBreakingChanges)...)
	}

	tmpContent := make([]string, 0)
	if len(changeLog.features.dataSources) > 0 {
		sort(changeLog.features.dataSources)
		tmpContent = append(tmpContent, changeLog.features.dataSources...)
	}

	if len(changeLog.features.general) > 0 {
		sort(changeLog.features.general)
		tmpContent = append(tmpContent, changeLog.features.general...)
	}

	if len(tmpContent) > 0 {
		newContent = append(newContent, formatSection(tmpContent, sectionTypeFeatures)...)
	}

	tmpContent = make([]string, 0)
	if len(changeLog.enhancements.dependencies) > 0 {
		sort(changeLog.enhancements.dependencies)
		tmpContent = append(tmpContent, changeLog.enhancements.dependencies...)
	}

	if len(changeLog.enhancements.dataSources) > 0 {
		sort(changeLog.enhancements.dataSources)
		tmpContent = append(tmpContent, changeLog.enhancements.dataSources...)
	}

	if len(changeLog.enhancements.general) > 0 {
		sort(changeLog.enhancements.general)
		tmpContent = append(tmpContent, changeLog.enhancements.general...)
	}

	if len(tmpContent) > 0 {
		newContent = append(newContent, formatSection(tmpContent, sectionTypeEnhancements)...)
	}

	tmpContent = make([]string, 0)
	if len(changeLog.bugs.dataSources) > 0 {
		sort(changeLog.bugs.dataSources)
		tmpContent = append(tmpContent, changeLog.bugs.dataSources...)
	}

	if len(changeLog.bugs.general) > 0 {
		sort(changeLog.bugs.general)
		tmpContent = append(tmpContent, changeLog.bugs.general...)
	}

	if len(tmpContent) > 0 {
		newContent = append(newContent, formatSection(tmpContent, sectionTypeBugs)...)
	}

	return append(newContent, changeLog.post...)
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
