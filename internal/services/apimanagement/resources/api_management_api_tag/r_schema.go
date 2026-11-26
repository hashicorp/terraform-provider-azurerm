// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementApiTagSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_id": {
			Description:  "The ID of the API Management API.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: api.ValidateApiID,
		},

		"name": {
			Description: "The name of the tag.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
	}
}
