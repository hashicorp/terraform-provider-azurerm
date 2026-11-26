// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
)

func FlattenTransformSlice(input *[]webapplicationfirewallpolicies.TransformType) []any {
	result := make([]any, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func FlattenFrontendEndpointLinkSlice(input *[]webapplicationfirewallpolicies.FrontendEndpointLink) []any {
	result := make([]any, 0)

	if input != nil {
		for _, item := range *input {
			if item.Id == nil {
				continue
			}

			result = append(result, *item.Id)
		}
	}
	return result
}
