// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rule

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	upgradeGuideContent  string
	loadUpgradeGuideOnce sync.Once
)

// HasVersionChanges checks if the field has version-related changes by parsing the upgrade guide
func HasVersionChanges(documentPath, resourceType, fieldPath string) bool {
	content := getUpgradeGuideContent(documentPath)
	if content == "" {
		return false
	}

	// Look for the resource section
	resourceSectionPattern := fmt.Sprintf("### `%s`", resourceType)

	// Find the resource section
	lines := strings.Split(content, "\n")
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

// getUpgradeGuideContent loads and caches the upgrade guide content
func getUpgradeGuideContent(documentPath string) string {
	loadUpgradeGuideOnce.Do(func() {
		// Extract base docs directory from document path
		var docsBasePath string
		if idx := strings.Index(documentPath, "website/docs/"); idx >= 0 {
			docsBasePath = documentPath[:idx+len("website/docs")]
		} else if idx := strings.Index(documentPath, "website\\docs\\"); idx >= 0 {
			docsBasePath = documentPath[:idx+len("website\\docs")]
		}

		if docsBasePath == "" {
			upgradeGuideContent = ""
			return
		}

		pattern := filepath.Join(docsBasePath, "*-upgrade-guide.html.markdown")
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			if content, err := os.ReadFile(matches[0]); err == nil {
				upgradeGuideContent = string(content)
				return
			}
		}

		upgradeGuideContent = ""
	})

	return upgradeGuideContent
}
