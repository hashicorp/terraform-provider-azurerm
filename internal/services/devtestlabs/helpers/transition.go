// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import "github.com/hashicorp/go-azure-helpers/lang/pointer"

func ExpandTags(input map[string]any) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}

func FlattenTags(input *map[string]string) map[string]*string {
	output := make(map[string]*string)

	if input != nil {
		for k, v := range *input {
			output[k] = pointer.To(v)
		}
	}

	return output
}
