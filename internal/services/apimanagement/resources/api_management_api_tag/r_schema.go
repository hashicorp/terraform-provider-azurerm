// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func apiManagementApiTagSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: api.ValidateApiID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}
