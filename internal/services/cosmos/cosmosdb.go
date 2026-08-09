package cosmos

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
)

func isServerlessCapacityMode(input *openapis.DatabaseAccountGetResults) bool {
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
