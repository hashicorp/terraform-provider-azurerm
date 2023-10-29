// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func Expand(tagsMap map[string]interface{}) map[string]*string {
	output := make(map[string]*string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[i] = &value
	}

	return output
}

func ExpandFrom(tagsMap *map[string]string) map[string]interface{} {
	output := make(map[string]interface{}, len(pointer.From(tagsMap)))

	for i, v := range pointer.From(tagsMap) {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[i] = &value
	}

	return output
}

func PointerTo(tagsMap map[string]interface{}) *map[string]string {
	output := make(map[string]string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[i] = value
	}

	return pointer.To(output)
}
