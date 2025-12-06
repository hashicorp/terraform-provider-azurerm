// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import "github.com/hashicorp/go-azure-helpers/lang/pointer"

func expandTags(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}

func flattenTags(input *map[string]string) map[string]*string {
	output := make(map[string]*string)

	if input != nil {
		for k, v := range *input {
			output[k] = pointer.To(v)
		}
	}

	return output
}
