// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_pipeline

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/pipelines"
)

func expandDataFactoryPipelineParameters(input map[string]any) *map[string]pipelines.ParameterSpecification {
	output := make(map[string]pipelines.ParameterSpecification)

	for k, v := range input {
		output[k] = pipelines.ParameterSpecification{
			Type:         pipelines.ParameterTypeString,
			DefaultValue: pointer.To(v),
		}
	}

	return &output
}

func expandDataFactoryPipelineVariables(input map[string]any) *map[string]pipelines.VariableSpecification {
	output := make(map[string]pipelines.VariableSpecification)

	for k, v := range input {
		output[k] = pipelines.VariableSpecification{
			Type:         pipelines.VariableTypeString,
			DefaultValue: pointer.To(v),
		}
	}

	return &output
}
