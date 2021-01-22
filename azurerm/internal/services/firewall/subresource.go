package firewall

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
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
