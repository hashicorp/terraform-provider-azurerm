// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_endpoint

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionendpoints"
)

func flattenDataCollectionEndpointKind(input *datacollectionendpoints.KnownDataCollectionEndpointResourceKind) string {
	if input == nil {
		return ""
	}

	return string(*input)
}

func flattenDataCollectionEndpointDescription(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func flattenDataCollectionEndpointPublicNetworkAccess(input *datacollectionendpoints.KnownPublicNetworkAccessOptions) bool {
	if input == nil {
		return false
	}
	var result bool
	switch *input {
	case datacollectionendpoints.KnownPublicNetworkAccessOptionsEnabled:
		result = true
	case datacollectionendpoints.KnownPublicNetworkAccessOptionsDisabled:
		result = false
	}
	return result
}
