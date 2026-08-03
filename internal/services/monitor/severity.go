// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var monitorSeverityNameToLevel = map[string]int{
	"critical":      0,
	"error":         1,
	"warning":       2,
	"informational": 3,
	"verbose":       4,
}

func normalizeMonitorSeverity(input string) (int, error) {
	value := strings.TrimSpace(strings.ToLower(input))

	if level, ok := monitorSeverityNameToLevel[value]; ok {
		return level, nil
	}

	switch value {
	case "0":
		return 0, nil
	case "1":
		return 1, nil
	case "2":
		return 2, nil
	case "3":
		return 3, nil
	case "4":
		return 4, nil
	default:
		return 0, fmt.Errorf("expected one of 0, 1, 2, 3, 4, critical, error, warning, informational, verbose")
	}
}

func normalizeMonitorSeverityState(input interface{}) string {
	value := fmt.Sprintf("%v", input)
	level, err := normalizeMonitorSeverity(value)
	if err != nil {
		return strings.TrimSpace(strings.ToLower(value))
	}

	return fmt.Sprintf("%d", level)
}

func validateMonitorSeverity(input interface{}, key string) (warnings []string, errors []error) {
	value, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", key))
		return warnings, errors
	}

	if _, err := normalizeMonitorSeverity(value); err != nil {
		errors = append(errors, fmt.Errorf("invalid %s %q: %s", key, value, err))
	}

	return warnings, errors
}

func suppressMonitorSeverityDiff(_, old, new string, _ *pluginsdk.ResourceData) bool {
	oldLevel, oldErr := normalizeMonitorSeverity(old)
	newLevel, newErr := normalizeMonitorSeverity(new)

	if oldErr != nil || newErr != nil {
		return false
	}

	return oldLevel == newLevel
}
