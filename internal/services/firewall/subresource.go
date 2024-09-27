// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/firewallpolicies"
)

func flattenNetworkSubResourceID(input *[]firewallpolicies.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.Id != nil {
			results = append(results, *item.Id)
		}
	}

	return results
}
