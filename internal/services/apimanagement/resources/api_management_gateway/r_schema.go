// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementGatewaySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_id": {
			Description:  "The ID of the API Management Resource in which the gateway will be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"description": {
			Description: "The description of the API Management Gateway.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"location_data": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "A canonical name for the geographic or physical location.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"city": {
						Description: "The city or locality where the resource is located.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"district": {
						Description: "The district, state, or province where the resource is located.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"region": {
						Description: "The country or region where the resource is located.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},
	}
}
