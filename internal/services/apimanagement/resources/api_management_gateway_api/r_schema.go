// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_api

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementGatewayApiSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_id": {
			Description:  "The Identifier of the API Management API within the API Management Service.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: api.ValidateApiID,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				oldApiTagId, err := api.ParseApiID(old)
				if err != nil {
					return false
				}

				newApiTagId, err := api.ParseApiID(new)
				if err != nil {
					return false
				}

				return oldApiTagId.ApiId == helpers.GetApiName(newApiTagId.ApiId)
			},
		},
		"gateway_id": {
			Description:  "The Identifier for the API Management Gateway.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}
