// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/connections/azuresdkhacks"
)

func expandParameterValueSet(input []any) *azuresdkhacks.ParameterValueSet {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	result := &azuresdkhacks.ParameterValueSet{
		Name: pointer.To(v["name"].(string)),
	}

	if values, ok := v["values"].(map[string]any); ok && len(values) > 0 {
		expandedValues := make(map[string]any)
		for key, val := range values {
			expandedValues[key] = map[string]any{
				"value": val,
			}
		}
		result.Values = pointer.To(expandedValues)
	}

	return result
}
