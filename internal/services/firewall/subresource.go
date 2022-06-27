package firewall

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/sdk/2022-01-01/network"
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
