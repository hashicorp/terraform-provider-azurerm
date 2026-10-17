// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gateway"
)

func expandApiManagementGatewayLocationData(input []any) *gateway.ResourceLocationDataContract {
	if len(input) == 0 {
		return nil
	}

	locationData := gateway.ResourceLocationDataContract{}

	vs := input[0].(map[string]any)
	for k, v := range vs {
		switch k {
		case "name":
			locationData.Name = v.(string)
		case "city":
			locationData.City = pointer.To(v.(string))
		case "district":
			locationData.District = pointer.To(v.(string))
		case "region":
			locationData.CountryOrRegion = pointer.To(v.(string))
		}
	}

	return &locationData
}
