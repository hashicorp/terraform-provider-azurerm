// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logic_app_action_http

// NOTE: this file is not a recommended way of developing Terraform resources; this exists to work around the fact that this API is dynamic (by its nature)
func flattenLogicAppActionRunAfter(input map[string]any) []any {
	if len(input) == 0 {
		return []any{}
	}
	output := []any{}
	for k, v := range input {
		output = append(output, map[string]any{
			"action_name":   k,
			"action_result": v.([]any)[0],
		})
	}

	return output
}

func expandLogicAppActionRunAfter(input []any) map[string]any {
	if len(input) == 0 {
		return nil
	}
	output := map[string]any{}
	for _, v := range input {
		b := v.(map[string]any)
		output[b["action_name"].(string)] = []string{b["action_result"].(string)}
	}

	return output
}
