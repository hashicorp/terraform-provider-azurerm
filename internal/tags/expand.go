// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package tags

import (
	rmtags "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
)

func Expand(tagsMap map[string]interface{}) map[string]*string {
	output := make(map[string]*string, len(tagsMap))

	for i, v := range tagsMap {
		// Validate should have ignored this error already
		value, _ := TagValueToString(v)
		output[i] = &value
	}

	// exclude any provider-level ignored tag keys from the desired set
	return rmtags.Ignore().ApplyMap(output)
}
