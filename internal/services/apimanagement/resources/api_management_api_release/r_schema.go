// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_api_release

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiReleaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this API Management API Release.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApiManagementChildName,
		},

		"api_id": {
			Description:  "The ID of the API Management API.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: api.ValidateApiID,
		},

		"notes": {
			Description:  "The Release Notes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
