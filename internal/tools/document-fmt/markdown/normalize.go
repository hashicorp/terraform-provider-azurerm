// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"regexp"
	"strings"
	"unicode"
)

// Regular expressions for normalization
var (
	tryBlockHeadReg = regexp.MustCompile("^([*] )?(An?|The) +(`[a-zA-Z0-9_]+`[, and]*)+blocks?.*$")
	blockHeadReg    = regexp.MustCompile("^(an?|An?|The)[^`]+(`[a-zA-Z0-9_]+`[, and]*)+.*blocks?.*$")
	tryBlockPropReg = regexp.MustCompile("[*] `.*` .*(One or more |A |A list of ) [^.,]* blocks?.*")
	defaultValueReg = regexp.MustCompile(`[. ]+[d|D]efaults to (\w+)[,.]`)
	spaceReg        = regexp.MustCompile(` \s+`)
)

// replaceNBSP replaces non-breaking spaces and other non-standard spaces with regular spaces
func replaceNBSP(line string) string {
	isNormalSpace := func(r rune) bool {
		switch r {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			return true
		}
		return false
	}
	var res []rune
	for _, ch := range line {
		if unicode.IsSpace(ch) && !isNormalSpace(ch) {
			res = append(res, rune(' '))
		} else {
			res = append(res, ch)
		}
	}
	return string(res)
}

// removeRedundantSpace removes redundant spaces in property lines
func removeRedundantSpace(line string) string {
	// only process property line
	if !strings.HasPrefix(line, "*") {
		return line
	}

	// some line with multi-space in beginning by attention
	if len(line) < 4 {
		return line
	}
	return line[:4] + spaceReg.ReplaceAllString(line[4:], " ")
}

// tryBlockHeadDetect detects if a line is a block header
func tryBlockHeadDetect(line string) bool {
	if tryBlockHeadReg.MatchString(line) {
		return true
	}

	// sometimes a line with supports
	if strings.Contains(line, "supports") && strings.HasSuffix(line, ":") {
		return true
	}
	return false
}

// tryFixBlockHead fixes block header formatting
func tryFixBlockHead(line string) string {
	if strings.HasPrefix(line, "*") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
	}

	if strings.HasPrefix(line, "`") {
		line = "The " + line
	}

	// add a block string after the second `
	if !strings.Contains(line, "block") {
		if idx := strings.LastIndexByte(line, '`'); idx > 0 {
			idx += 1
			line = line[:idx] + " block" + line[idx:]
		}
	}
	return line
}

// tryFixProp performs detailed property line fixes
func tryFixProp(line string) string {
	// Convert attention markers
	if strings.HasPrefix(line, "**") {
		return "~> " + line
	}

	// Fix spacing around dash with Required/Optional
	if reqIdx := requireIndex(line); reqIdx > 0 {
		if !strings.HasSuffix(strings.TrimSpace(line[:reqIdx]), "-") {
			line = line[:reqIdx] + "- " + line[reqIdx:]
			reqIdx += 2
		}

		// a blank character before requiredness
		if divIdx := strings.Index(line[:reqIdx], "-"); divIdx > 0 {
			if ch := line[divIdx-1]; ch != ' ' {
				if unicode.IsSpace(rune(ch)) {
					line = line[:divIdx-1] + " " + line[divIdx:]
				} else {
					line = line[:divIdx] + " " + line[divIdx:]
				}
				divIdx += 1
			}
			if ch := line[divIdx+1]; ch != ' ' {
				if unicode.IsSpace(rune(ch)) {
					line = line[:divIdx+1] + " " + line[divIdx+2:]
				} else {
					line = line[:divIdx+1] + " " + line[divIdx+1:]
				}
			}
		}
	}

	// need a dash after property name
	idx := strings.Index(line, "`")
	if idx >= 0 {
		if idx2 := strings.Index(line[idx+1:], "`"); idx2 >= 0 {
			idx = idx + idx2 + 1
			for idx2 := idx + 1; idx2 < len(line); idx2++ {
				if line[idx2] == ' ' {
					continue
				}
				if line[idx2] != '-' {
					line2 := line[:idx+2]
					if idx2 > 0 && line[idx2-1] != ' ' {
						line2 += " "
					}
					line2 += "- "
					line2 += line[idx2:]
					line = line2
					break
				} else {
					// if line[idx2] == '-' exists
					if idx2 > 0 && line[idx2-1] != ' ' {
						line = line[:idx2] + " " + line[idx2:]
					}
					if idx2+1 < len(line) && line[idx2+1] != ' ' {
						line = line[:idx2+1] + " " + line[idx2+1:]
					}
					break
				}
			}
		}
	}

	// add ` to default value
	if vals := defaultValueReg.FindStringSubmatchIndex(line); len(vals) > 0 {
		// add a backquote to value
		valStr := line[vals[2]:vals[3]]
		line = line[:vals[2]] + "`" + valStr + "`" + line[vals[3]:]
	}

	// need a blank for block properties
	if tryBlockPropReg.MatchString(line) {
		if !strings.Contains(line, "below") && !strings.Contains(line, "above") {
			line = strings.TrimSuffix(line, ".")
			line += " as defined below."
		}
	}

	return line
}

// requireIndex returns the index of (Optional) or (Required) in the line
func requireIndex(line string) int {
	if idx := strings.Index(line, "(Optional)"); idx > 0 {
		return idx
	}

	if idx := strings.Index(line, "(Required)"); idx > 0 {
		return idx
	}
	return -1
}
