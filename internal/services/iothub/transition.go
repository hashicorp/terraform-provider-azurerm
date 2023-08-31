// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import "github.com/hashicorp/terraform-provider-azurerm/utils"

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
			output[k] = utils.String(v)
		}
	}

	return output
}
