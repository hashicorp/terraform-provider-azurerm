// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_tag_description

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apitag"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func apiManagementApiTagDescriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_tag_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apitag.ValidateApiTagID,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"external_documentation_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"external_documentation_description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}
