// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"encoding/json"
)

func flattenApplicationInsightsSourceMapURI(input *map[string]string) string {
	if input == nil {
		return ""
	}
	raw, ok := (*input)[appInsightsSourceMapStorageTag]
	if !ok || raw == "" {
		return ""
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return ""
	}
	return parsed["Uri"]
}

func stripApplicationInsightsSourceMapTag(input *map[string]string) *map[string]string {
	if input == nil {
		return nil
	}
	// Avoid surfacing the provider-managed hidden-link as a user tag (prevents drift
	// when the first-class attribute is used).
	out := make(map[string]string, len(*input))
	for k, v := range *input {
		if k == appInsightsSourceMapStorageTag {
			continue
		}
		out[k] = v
	}
	return &out
}
