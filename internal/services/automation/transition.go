// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
)

func interfaceValueToString(v any) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown type %T in value", value)
	}
}

func expandStringInterfaceMap(strInterfaceMap map[string]any) map[string]string {
	output := make(map[string]string, len(strInterfaceMap))

	for i, v := range strInterfaceMap {
		// Validate should have ignored this error already
		value, _ := interfaceValueToString(v)
		output[i] = value
	}

	return output
}

func flattenMap(strStrMap map[string]string) map[string]any {
	// If strStrMap is nil, len(strStrMap) will be 0.
	output := make(map[string]any, len(strStrMap))

	for i, v := range strStrMap {
		output[i] = v
	}

	return output
}
