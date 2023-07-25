// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

// Expand transforms the input Tags to a `*map[string]string`
func Expand(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		tagKey := k
		tagValue := v.(string)
		output[tagKey] = tagValue
	}

	return &output
}
