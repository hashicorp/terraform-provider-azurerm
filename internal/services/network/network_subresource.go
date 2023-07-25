// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func flattenNetworkSubResourceID(input *[]network.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.ID != nil {
			results = append(results, *item.ID)
		}
	}

	return results
}
