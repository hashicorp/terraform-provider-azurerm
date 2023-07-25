// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func interfaceValueToString(v interface{}) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int:
		return fmt.Sprintf("%d", value), nil
	default:
		return "", fmt.Errorf("unknown type %T in value", value)
	}
}

func expandStringInterfaceMap(strInterfaceMap map[string]interface{}) map[string]string {
	output := make(map[string]string, len(strInterfaceMap))

	for i, v := range strInterfaceMap {
		// Validate should have ignored this error already
		value, _ := interfaceValueToString(v)
		output[i] = value
	}

	return output
}

func flattenMap(strStrMap map[string]string) map[string]interface{} {
	// If strStrMap is nil, len(strStrMap) will be 0.
	output := make(map[string]interface{}, len(strStrMap))

	for i, v := range strStrMap {
		output[i] = v
	}

	return output
}

func flattenAndSetTags(d *pluginsdk.ResourceData, tagMap map[string]string) error {
	flattened := flattenMap(tagMap)
	if err := d.Set("tags", flattened); err != nil {
		return fmt.Errorf("setting `tags`: %s", err)
	}

	return nil
}
