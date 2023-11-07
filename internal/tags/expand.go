// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

func Expand(tagsMap map[string]interface{}) map[string]*string {
	output := make(map[string]*string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[i] = &value
	}

	return output
}
