// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
)

func ExpandIsDeviceMatchValue(input []interface{}) *[]rules.IsDeviceMatchValue {
	result := make([]rules.IsDeviceMatchValue, 0)
	for _, item := range input {
		if item != nil {
			result = append(result, rules.IsDeviceMatchValue(item.(string)))
		}
	}
	return &result
}

func FlattenIsDeviceMatchValue(input *[]rules.IsDeviceMatchValue) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}
