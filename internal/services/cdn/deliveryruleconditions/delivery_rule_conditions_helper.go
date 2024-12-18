// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
)

func expandTransforms(input []interface{}) *[]rules.Transform {
	if len(input) == 0 {
		return nil
	}

	result := make([]rules.Transform, 0)

	if v := input; len(v) != 0 {
		for _, t := range v {
			result = append(result, rules.Transform(t.(string)))
		}

		return &result
	}

	return nil
}

func flattenTransforms(input *[]rules.Transform) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, transform := range *input {
		result = append(result, string(transform))
	}

	return result
}

func expandIsDeviceMatchValue(input []interface{}) *[]rules.IsDeviceMatchValue {
	result := make([]rules.IsDeviceMatchValue, 0)

	for _, item := range input {
		if item != nil {
			result = append(result, rules.IsDeviceMatchValue(item.(string)))
		}
	}

	return &result
}

func flattenIsDeviceMatchValue(input *[]rules.IsDeviceMatchValue) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}

func expandRequestMethodMatchValue(input []interface{}) *[]rules.RequestMethodMatchValue {
	result := make([]rules.RequestMethodMatchValue, 0)

	for _, item := range input {
		if item != nil {
			result = append(result, rules.RequestMethodMatchValue(item.(string)))
		}
	}

	return &result
}

func flattenRequestMethodMatchValue(input *[]rules.RequestMethodMatchValue) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}

func expandRequestSchemeMatchValue(input []interface{}) *[]rules.RequestSchemeMatchValue {
	result := make([]rules.RequestSchemeMatchValue, 0)

	for _, item := range input {
		if item != nil {
			result = append(result, rules.RequestSchemeMatchValue(item.(string)))
		}
	}

	return &result
}

func flattenRequestSchemeMatchValue(input *[]rules.RequestSchemeMatchValue) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}
