// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"regexp"
	"strings"
)

type ArgumentsSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &ArgumentsSection{}

func (s *ArgumentsSection) Match(line string) bool {
	return regexp.MustCompile(`#+(\s)*argument.*`).MatchString(strings.ToLower(line))
}

func (s *ArgumentsSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *ArgumentsSection) GetHeading() Heading {
	return s.heading
}

func (s *ArgumentsSection) SetContent(content []string) {
	s.content = content
}

func (s *ArgumentsSection) GetContent() []string {
	return s.content
}

func (s *ArgumentsSection) Template() string {
	// TODO implement me
	panic("implement me")
}

// Normalize applies formatting normalization to the Arguments section content
// This should be called before parsing to clean up common formatting issues
func (s *ArgumentsSection) Normalize() (normalizedContents []string, hasChange bool) {
	if len(s.content) == 0 {
		return nil, false
	}
	// Simply call the existing normalizeArgumentsContent helper
	return normalizeArgumentsContent(s.content)
}

// normalizeArgumentsContent normalizes Arguments section content without section detection
func normalizeArgumentsContent(lines []string) (normalizedContents []string, hasChange bool) {
	normalized := make([]string, 0, len(lines))
	var skipThisLine int
	var inCodeBlock bool
	hasChange = false

	for idx, line := range lines {
		originalLine := line

		// Handle code block detection
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			normalized = append(normalized, line)
			continue
		}

		// Skip processing inside code blocks
		if inCodeBlock {
			normalized = append(normalized, line)
			continue
		}

		// Handle multi-line skip
		if skipThisLine > 0 {
			skipThisLine--
			continue
		}

		// Replace non-standard spaces (NBSP) with regular spaces
		line = replaceNBSP(line)

		// Fix separator: "--" -> "---"
		if line == "--" {
			line = "---"
		}

		// Fix list marker: convert "- `" to "* `"
		if strings.HasPrefix(line, "- `") {
			line = "*" + line[1:]
		}

		// Multi-line property merging: combine properties that span multiple lines
		if strings.HasPrefix(line, "*") && !strings.HasSuffix(line, ".") {
			idx2 := idx + 1
			for idx2 < len(lines) {
				l2 := lines[idx2]
				if l2 == "" {
					break
				}
				ch := l2[0]
				if ch == ' ' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
					if !strings.HasSuffix(line, " ") && !strings.HasPrefix(l2, " ") {
						line += " "
					}
					line += l2
					skipThisLine++
					idx2++
				} else {
					break
				}
			}
		}

		// Block head detection and formatting
		if tryBlockHeadDetect(line) {
			line = tryFixBlockHead(line)
		}

		// Add separator before block heads
		if blockHeadReg.MatchString(line) {
			isSep := func(l string) bool {
				return l == "---" || strings.HasPrefix(l, "#")
			}
			if idx > 1 && !isSep(lines[idx-1]) && !isSep(lines[idx-2]) {
				normalized = append(normalized, "---", "")
			}
		}

		// Fix Required/Optional position and case FIRST (before tryFixProp)
		line = regexp.MustCompile(`(?i)-\s*\((optional)\)\s*-`).ReplaceAllStringFunc(line, func(s string) string {
			return "- (Optional)"
		})
		line = regexp.MustCompile(`(?i)-\s*\((required)\)\s*-`).ReplaceAllStringFunc(line, func(s string) string {
			return "- (Required)"
		})

		line = regexp.MustCompile(`(?i)\((optional)\)\s*-`).ReplaceAllStringFunc(line, func(s string) string {
			return "- (Optional)"
		})
		line = regexp.MustCompile(`(?i)\((required)\)\s*-`).ReplaceAllStringFunc(line, func(s string) string {
			return "- (Required)"
		})

		// Property line processing (after position fixes to avoid double dash)
		if strings.HasPrefix(line, "*") {
			for k, v := range OrderFixMap {
				if strings.Contains(line, k) && !strings.Contains(line, v) {
					line = strings.Replace(line, k, v, 1)
				}
			}
			for k, v := range RequiredCaseFix {
				if strings.Contains(line, k) && !strings.Contains(line, v) {
					line = strings.Replace(line, k, v, 1)
				}
			}
			line = tryFixProp(line)
		}

		// Add missing marker prefix for properties
		if (strings.Contains(line, "(Optional)") || strings.Contains(line, "(Required)")) && strings.HasPrefix(line, "`") {
			line = "* " + line
		}

		// Remove redundant spaces
		line = removeRedundantSpace(line)

		normalized = append(normalized, line)

		if originalLine != line {
			hasChange = true
		}
	}

	return normalized, hasChange
}
