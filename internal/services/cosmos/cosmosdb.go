// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
)

func isServerlessCapacityMode(input *cosmosdb.DatabaseAccountGetResults) bool {
	if input == nil || input.Properties == nil || input.Properties.Capabilities == nil {
		return false
	}

	for _, v := range *input.Properties.Capabilities {
		if pointer.From(v.Name) == "EnableServerless" {
			return true
		}
	}

	return false
}
