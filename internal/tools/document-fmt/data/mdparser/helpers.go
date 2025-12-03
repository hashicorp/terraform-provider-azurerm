// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdparser

import (
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
)

// Regular expressions for parsing
var (
	fieldReg     = regexp.MustCompile("^[*-] *`(.*?)`" + ` +\- +(\(Required\)|\(Optional\))? ?(.*)`)
	codeReg      = regexp.MustCompile("`([^`]+)`")
	blockHeadReg = regexp.MustCompile("^(an?|An?|The)[^`]+(`[a-zA-Z0-9_]+`[, and]*)+.*blocks?.*$")

	// DefaultsReg matches various default value patterns:
	// - "[.,?;] defaults to X" (with punctuation)
	// - "This value defaults to X" (without leading punctuation)
	// - "It defaults to X"
	// - "The default value is X"
	// - "Default value is X"
	// - "xxx default value is X"
	DefaultsReg     = regexp.MustCompile("[.,?;]?(?:(?: *[Tt]he)? *[Dd]efault value|[Ii]t defaults| *[Dd]efaults?)[^`'\".]*(?:to|is) ('[^']+'|`[^`]+`|\"[^\"]+\")[ .,]?")
	forceNewReg     = regexp.MustCompile(` ?Changing.*forces? a [^.]*(\.|$)`)
	partforceNewReg = regexp.MustCompile(` ?Changing.*forces? a [^.]* created when [^.]*(\.|$)`)
	blockPropRegs   = []*regexp.Regexp{
		regexp.MustCompile("(?:[Oo]ne|[Ee]ach|more(?: \\(.*\\))?|[Tt]he|as|of|[Aa]n?) ['\"`]([^ ]+)['\"`] (?:block|object)[^.]+(?:below|above)"),
	}
	blockTypeReg = blockPropRegs[0]
)

// ===== Field parsing helpers =====

// getDefaultValue extracts the default value from a field description line
func getDefaultValue(line string) string {
	if vals := DefaultsReg.FindStringSubmatch(line); len(vals) > 0 {
		if val := vals[1]; len(val) > 2 {
			return val[1 : len(val)-1]
		}
	}
	return ""
}

// isForceNew determines if a field triggers resource recreation when changed
func isForceNew(line string) bool {
	if forceNewReg.MatchString(line) && !partforceNewReg.MatchString(line) {
		return true
	}
	return false
}

// extractFieldFromLine parses a field definition line and extracts its properties
func extractFieldFromLine(line string) *models.DocumentProperty {
	// remove redundant /n
	cleanedLine := strings.TrimRight(line, "/n")
	field := &models.DocumentProperty{
		Content: cleanedLine,
	}

	if defaultVal := getDefaultValue(line); defaultVal != "" {
		field.DefaultValue = defaultVal
	}
	field.ForceNew = isForceNew(line)

	res := fieldReg.FindStringSubmatch(line)
	if len(res) <= 1 || res[1] == "" {
		field.Name = util.FirstCodeValue(line)
		if field.ParseErrors == nil {
			field.ParseErrors = []string{}
		}
		field.ParseErrors = append(field.ParseErrors, "no field name found")
		return field
	}
	field.Name = res[1]
	if field.Name == "" {
		log.Printf("field name is empty")
	}
	if len(res) > 2 {
		switch {
		case strings.Contains(line, "(Required)"):
			field.Required = true
		case strings.Contains(line, "(Optional)"):
			field.Optional = true
		case strings.Contains(line, "Required"):
			field.Required = true
		case strings.Contains(line, "Optional"):
			field.Optional = true
		}
	}

	possibleValueSep := func(line string) int {
		line = strings.ToLower(line)
		for _, sep := range []string{
			"possible value", "must be one of", "be one of", "one of", "allowed value", "valid value", "be set to",
			"supported value", "valid option", "accepted value", "acceptable value", "allowable value",
		} {
			if sepIdx := strings.Index(line, sep); sepIdx >= 0 {
				return sepIdx
			}
		}
		return -1
	}

	var enums []string
	if len(res) > 3 {
		// extract enums from code part
		// from possible value to first '.'
		// skip if there are more than one sep exists
		// do not check the possible part
		if sepIdx := possibleValueSep(line); sepIdx > 0 {
			subStr := line[sepIdx:]
			field.EnumStart = sepIdx
			// end with dot may not work in values like `7.2` ....
			// should be . not in ` mark
			// Possible values are `a`, `b`, `a.b` and `def`.
			pointEnd := strings.Index(subStr, ".")
			if pointEnd < 0 {
				pointEnd = len(subStr)
			}

			// Track parentheses depth to skip values inside explanatory text
			parenDepth := 0
			enumIndex := codeReg.FindAllStringIndex(subStr, -1)
			for _, val := range enumIndex {
				start, end := val[0], val[1]

				// Update paren depth for content before this code block
				for i := 0; i < start && i < len(subStr); i++ {
					if subStr[i] == '(' {
						parenDepth++
					} else if subStr[i] == ')' {
						parenDepth--
					}
				}

				// Skip values that are inside parentheses (explanatory text)
				if parenDepth > 0 {
					continue
				}

				if pointEnd > start && pointEnd < end {
					// point inside the code block
					if pointEnd = strings.Index(subStr[end:], "."); pointEnd < 0 {
						pointEnd = len(subStr)
					} else {
						pointEnd += end
					}
				}
				if pointEnd < start {
					break
				}
				enums = append(enums, strings.Trim(subStr[start:end], "`'\""))
				field.EnumEnd = sepIdx + end
			}
			// breaks if there are more than 1 possible value
			checkFromIdx := field.EnumEnd
			if checkFromIdx < len(line) {
				if sepIdx = possibleValueSep(line[checkFromIdx:]); sepIdx >= 0 {
					field.Skip = true
				}
			}
		}
		if len(enums) == 0 && strings.Index(res[3], "`") > 0 {
			guessValues := codeReg.FindAllString(res[3], -1)
			field.GuessEnums = setGuessEnums(guessValues)
		}
	}
	field.AddEnum(enums...)
	return field
}

// setGuessEnums deduplicates and cleans up guessed enum values
func setGuessEnums(values []string) []string {
	hys := make(map[string]struct{}, len(values))
	var res []string
	for _, val := range values {
		val = strings.Trim(val, "`\"'")
		if _, ok := hys[val]; !ok {
			hys[val] = struct{}{}
			res = append(res, val)
		}
	}
	return res
}

// newFieldFromLine creates a DocumentProperty from a field definition line
func newFieldFromLine(line string) *models.DocumentProperty {
	f := extractFieldFromLine(line)
	if guessBlockProperty(line) {
		// extract real block type
		f.BlockTypeName = f.Name
		if match := blockTypeReg.FindAllStringSubmatchIndex(strings.ToLower(line), -1); len(match) > 0 {
			f.BlockTypeName = line[match[0][2]:match[0][3]]
		}
		f.Block = true
	}
	return f
}

// ===== Block parsing helpers =====

// extractBlockNames extracts all block names from a block header line
func extractBlockNames(line string) []string {
	if blockHeadReg.MatchString(line) {
		idx := strings.Index(line, "block")
		names := codeReg.FindAllString(line[:idx], -1)
		for idx, val := range names {
			names[idx] = strings.Trim(val, "`'")
		}
		return names
	}
	return nil
}

// isBlockHead determines if a line is a block header
func isBlockHead(line string) bool {
	return blockHeadReg.MatchString(line)
}

// guessBlockProperty determines if a field description suggests it's a block type
func guessBlockProperty(line string) bool {
	for _, reg := range blockPropRegs {
		if reg.MatchString(line) {
			return true
		}
	}
	return strings.Contains(line, "A block to")
}

// blocksHaveSameDefinition checks if two blocks have the same definition/content
func blocksHaveSameDefinition(b1, b2 *markBlock) bool {
	if len(b1.Fields) != len(b2.Fields) {
		return false
	}

	for i, f1 := range b1.Fields {
		if i >= len(b2.Fields) {
			return false
		}
		f2 := b2.Fields[i]
		if f1.Name != f2.Name || f1.Required != f2.Required {
			return false
		}
	}

	return true
}
