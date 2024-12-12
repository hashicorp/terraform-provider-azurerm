// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/datasets"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func expandDataSetParameters(input map[string]interface{}) map[string]*datafactory.ParameterSpecification {
	output := make(map[string]*datafactory.ParameterSpecification)

	for k, v := range input {
		output[k] = &datafactory.ParameterSpecification{
			Type:         datafactory.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func expandDataSetParametersGoAzureSdk(input *map[string]string) *map[string]datasets.ParameterSpecification {
	output := make(map[string]datasets.ParameterSpecification)

	for k, v := range *input {
		var value interface{} = v
		output[k] = datasets.ParameterSpecification{
			Type:         datasets.ParameterTypeString,
			DefaultValue: &value,
		}
	}

	return &output
}

func flattenDataSetParameters(input map[string]*datafactory.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func flattenDataSetParametersGoAzureSdk(input *map[string]datasets.ParameterSpecification) map[string]string {
	output := make(map[string]string)

	if input == nil {
		return output
	}

	for k, v := range *input {
		if v.DefaultValue != nil {
			val, ok := (*v.DefaultValue).(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since its default value is not a string", k)
				continue
			}

			output[k] = val
		}
	}

	return output
}
