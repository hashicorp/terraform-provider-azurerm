// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_collector_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"
)

func flattenEmissionPoliciesPropertiesFormatModelArray(inputList *[]collectorpolicies.EmissionPoliciesPropertiesFormat) []IpfxEmissionModel {
	outputList := make([]IpfxEmissionModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := IpfxEmissionModel{
			DestinationTypes: flattenEmissionPolicyDestinationModelArray(input.EmissionDestinations),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenEmissionPolicyDestinationModelArray(inputList *[]collectorpolicies.EmissionPolicyDestination) []string {
	outputList := make([]string, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		if input.DestinationType != nil {
			outputList = append(outputList, string(*input.DestinationType))
		}
	}

	return outputList
}

func flattenIngestionSourcesPropertiesFormatModelArray(inputList *[]collectorpolicies.IngestionSourcesPropertiesFormat) []IpfxIngestionModel {
	outputList := make([]IpfxIngestionModel, 0)
	if inputList == nil {
		return outputList
	}

	output := IpfxIngestionModel{
		SourceResourceIds: make([]string, 0),
	}

	for _, input := range *inputList {
		if input.ResourceId != nil && input.SourceType != nil && *input.SourceType == collectorpolicies.SourceTypeResource {
			output.SourceResourceIds = append(output.SourceResourceIds, *input.ResourceId)
		}
	}

	if len(output.SourceResourceIds) > 0 {
		outputList = append(outputList, output)
	}

	return outputList
}
