// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/connections/azuresdkhacks"
)

func flattenParameterValueSet(input *azuresdkhacks.ParameterValueSet) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":   pointer.From(input.Name),
			"values": flattenParameterValueSetValues(input.Values),
		},
	}
}

// flattenParameterValueSetValues extracts values from the API's parameter value set format.
// The API returns values in the format {"key": {"value": "actualValue"}}
// This function extracts the "value" field for each key.
func flattenParameterValueSetValues(input *map[string]any) map[string]string {
	values := make(map[string]string)
	if input == nil {
		return values
	}

	for key, val := range *input {
		if valueMap, ok := val.(map[string]any); ok {
			if v, exists := valueMap["value"]; exists {
				values[key] = fmt.Sprintf("%v", v)
			}
		} else {
			values[key] = fmt.Sprintf("%v", val)
		}
	}
	return values
}
