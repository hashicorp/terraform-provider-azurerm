// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

type possibleValueDiff struct {
	checkBase

	Want []string
	Got  []string

	Missed []string // value not exists in doc
	Spare  []string // value redundant in doc, not exists in code
}

func newPossibleValueDiff(checkBase checkBase, want []string, got []string, missed []string, spare []string) *possibleValueDiff {
	return &possibleValueDiff{
		checkBase: checkBase,
		Want:      want,
		Got:       got,
		Missed:    missed,
		Spare:     spare,
	}
}

func possibleValueStr(values []string) string {
	codes := make([]string, 0, len(values))
	for _, val := range values {
		codes = append(codes, util.ItalicCode(val))
	}
	return fmt.Sprintf("[%s]", strings.Join(codes, ", "))
}

func (p possibleValueDiff) String() string {
	var missInDoc, missInCode string
	if len(p.Missed) > 0 {
		missInDoc = fmt.Sprintf(" the following possible values are missing in the documentation: %s.", possibleValueStr(p.Missed))
	}
	if len(p.Spare) > 0 {
		missInCode = fmt.Sprintf(" the following possible values are missing in the schema: %v.", possibleValueStr(p.Spare))
	}
	return fmt.Sprintf(`%s:%s%s`,
		p.Str(),
		missInDoc,
		missInCode,
	)
}

func (p possibleValueDiff) Fix(line string) (result string, err error) {
	if len(p.Want) == 0 {
		return line, nil
	}
	result = line
	// replace from field.EnumStart to field.EnumEnd
	var bs strings.Builder
	if len(p.Got) == 0 || (p.MDField().EnumStart == 0 && len(p.Missed) > 0) {
		// skip this kind of field. may submit in a separate run
		// find default index
		idx := strings.Index(line, "Defaults to")
		if idx < 0 {
			idx = strings.Index(line, "Changing this forces")
		}
		if idx > 0 {
			bs.WriteString(line[:idx])
		} else {
			bs.WriteString(line)
			bs.WriteByte(' ')
		}
		if len(p.Want) == 1 {
			bs.WriteString("The only possible value is ")
		} else {
			bs.WriteString("Possible values are ")
		}
		bs.WriteString(patchWantEnums(p.Want))
		bs.WriteByte('.')
		if idx > 0 {
			bs.WriteByte(' ')
			bs.WriteString(line[idx:])
		}
		result = bs.String()
	} else if len(p.Missed) > 0 {
		// only replace missed values
		bs.WriteString(line[:p.MDField().EnumStart])
		if len(p.Want) == 1 {
			bs.WriteString("The only possible value is ")
		} else {
			bs.WriteString("Possible values are ")
		}
		bs.WriteString(patchWantEnums(p.Want))
		if end := p.MDField().EnumEnd; end > 0 && end < len(line) {
			bs.WriteString(line[p.MDField().EnumEnd:])
		} else {
			f := p.MDField()
			log.Printf("warning enum end %s:L%d len %dvs%d; %s", path.Base(f.Path), f.Line, f.EnumEnd, len(line), line)
		}
		result = bs.String()
	}
	return result, nil
}

var _ Checker = (*possibleValueDiff)(nil)

// sortPossibleValues sorts possible values in a sensible order
// For weekdays, it maintains chronological order (Monday -> Sunday)
// For other values, it sorts alphabetically (case-insensitive)
func sortPossibleValues(values []string) []string {
	if len(values) <= 1 {
		return values
	}

	// Check if all values are weekdays
	weekdayOrder := map[string]int{
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
		"sunday":    7,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
		"Sunday":    7,
	}

	allWeekdays := true
	for _, v := range values {
		if _, ok := weekdayOrder[v]; !ok {
			allWeekdays = false
			break
		}
	}

	sorted := make([]string, len(values))
	copy(sorted, values)

	if allWeekdays {
		sort.Slice(sorted, func(i, j int) bool {
			return weekdayOrder[sorted[i]] < weekdayOrder[sorted[j]]
		})
	} else {
		sort.Slice(sorted, func(i, j int) bool {
			return strings.ToLower(sorted[i]) < strings.ToLower(sorted[j])
		})
	}

	return sorted
}

func patchWantEnums(want []string) string {
	// Sort the values before formatting
	sorted := sortPossibleValues(want)

	res := make([]string, len(sorted))
	for idx, val := range sorted {
		res[idx] = "`" + val + "`"
	}
	if len(res) == 1 {
		return res[0]
	}

	s := res[0]
	if len(res) >= 3 {
		s = strings.Join(res[:len(res)-1], ", ")
	}
	s += " and " + res[len(res)-1]
	return s
}

// check possible values
func checkPossibleValues(r *schema.Resource, md *model.ResourceDoc) (res []Checker) {
	schemModel := r.Schema.Schema
	_ = schemModel
	if md == nil {
		log.Printf("%s no match document exists", r.ResourceType)
		return
	}
	// loop over document model
	for name, field := range md.Args {
		partRes := diffField(r, field, []string{name})
		res = append(res, partRes...)
	}
	return
}

// xPath property name for parent nodes
func diffField(r *schema.Resource, mdField *model.Field, xPath []string) (res []Checker) {
	fullPath := strings.Join(xPath, ".")
	if isSkipProp(r.ResourceType, fullPath) {
		return
	}

	// if end property
	if mdField.Subs == nil {
		wanted := r.PossibleValues[fullPath]
		docVal := mdField.PossibleValues()
		if missed, spare := SliceDiff(wanted, docVal, true); len(missed)+len(spare) > 0 {
			// Check if this property has version changes before reporting error
			// Skip possible value diff if this property has changes in new version upgrade.
			// This is because it's not consistent whether possible values change should be documented or not when it's in the upgrade feature flag
			if hasVersionChanges(r.ResourceType, fullPath) {
				log.Printf("[SKIP] %s.%s: Skipping possible values validation due to version changes detected in upgrade guide (missed: %v, spare: %v)",
					r.ResourceType, fullPath, missed, spare)
				return
			}

			if !mayExistsInDoc(mdField.Content, wanted) {
				base := newCheckBase(mdField.Line, fullPath, mdField)
				res = append(res, newPossibleValueDiff(base, wanted, docVal, missed, spare))
			}
		}
		return
	}
	// check if r has such path
	if !r.HasPathFor(xPath) {
		log.Printf("%s %s has no path [%s], there must be an error in markdwon", color.YellowString("[WARN]"), r.ResourceType, strings.Join(xPath, "."))
		return
	}
	for _, sub := range mdField.Subs {
		subRes := diffField(r, sub, append(xPath, sub.Name))
		res = append(res, subRes...)
	}
	return
}

func SliceDiff(want, got []string, caseInSensitive bool) (missed, spare []string) {
	// if `want` is nil then it may only write in doc, skip this
	if len(want) == 0 {
		return
	}
	// cross-check
	wantCpy, gotCpy := want, got
	if caseInSensitive {
		wantCpy = make([]string, len(want))
		gotCpy = make([]string, len(got))
		for idx := range want {
			wantCpy[idx] = strings.ToLower(want[idx])
		}
		for idx := range got {
			gotCpy[idx] = strings.ToLower(got[idx])
		}
	}
	wantMap := util.Slice2Map(wantCpy)
	gotMap := util.Slice2Map(gotCpy)

	for idx, k := range wantCpy {
		if _, ok := gotMap[k]; !ok {
			missed = append(missed, want[idx])
		}
	}

	for idx, k := range gotCpy {
		if _, ok := wantMap[k]; !ok {
			spare = append(spare, got[idx])
		}
	}

	return
}

// return true values exists in doc but may not with code quote
func mayExistsInDoc(docLine string, want []string) bool {
	for _, val := range want {
		if !strings.Contains(docLine, val) {
			return false
		}
	}
	return true
}

// hasVersionChanges checks if the field has version-related changes by parsing the upgrade guide
func hasVersionChanges(resourceType, fieldPath string) bool {
	upgradeGuideContent := getUpgradeGuideContent()
	if upgradeGuideContent == "" {
		return false
	}

	// Look for the resource section
	resourceSectionPattern := fmt.Sprintf("### `%s`", resourceType)

	// Find the resource section
	lines := strings.Split(upgradeGuideContent, "\n")
	resourceSectionStart := -1

	for i, line := range lines {
		if strings.Contains(line, resourceSectionPattern) {
			resourceSectionStart = i
			break
		}
	}

	// Resource not found in upgrade guide
	if resourceSectionStart == -1 {
		return false
	}

	// Find the end of this resource section (next resource or major section)
	resourceSectionEnd := len(lines)
	for i := resourceSectionStart + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "### `azurerm_") || strings.HasPrefix(line, "## ") {
			resourceSectionEnd = i
			break
		}
	}

	// Get the content of this resource section
	resourceSection := strings.Join(lines[resourceSectionStart:resourceSectionEnd], "\n")

	// Check if the property is mentioned in this resource section
	propertyPatterns := []string{
		fmt.Sprintf("`%s`", fieldPath), // `property_name`
		fmt.Sprintf(" %s ", fieldPath), // property_name with spaces
		fmt.Sprintf(".%s ", fieldPath), // .property_name
		fmt.Sprintf("`%s.", fieldPath), // `property_name.
	}

	// For nested properties like "site_config.remote_debugging_version", also check the last part
	if strings.Contains(fieldPath, ".") {
		parts := strings.Split(fieldPath, ".")
		lastPart := parts[len(parts)-1]
		propertyPatterns = append(propertyPatterns,
			fmt.Sprintf("`%s`", lastPart), // `last_part`
			fmt.Sprintf(" %s ", lastPart), // last_part with spaces
		)
	}

	// Check if any pattern is found in the resource section
	for _, pattern := range propertyPatterns {
		if strings.Contains(resourceSection, pattern) {
			return true
		}
	}

	return false
}

var (
	upgradeGuideContent  string
	loadUpgradeGuideOnce sync.Once
)

// getUpgradeGuideContent loads and caches the upgrade guide content
func getUpgradeGuideContent() string {
	loadUpgradeGuideOnce.Do(func() {
		docsDir := md.DocDir()

		if upgradeFile := findUpgradeGuide(docsDir); upgradeFile != "" {
			if content, err := os.ReadFile(upgradeFile); err == nil {
				upgradeGuideContent = string(content)
				log.Printf("Loaded upgrade guide from: %s", upgradeFile)
				return
			}
		}

		log.Printf("Warning: Could not find any *-upgrade-guide.html.markdown file in docs directory: %s", docsDir)
		upgradeGuideContent = ""
	})

	return upgradeGuideContent
}

// findUpgradeGuide searches for upgrade guide files in the given directory
func findUpgradeGuide(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return ""
	}

	pattern := filepath.Join(dir, "*-upgrade-guide.html.markdown")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}

	return matches[0]
}
