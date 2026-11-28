// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package static_web_app_function_app_registration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"static_web_app_id": {
			Description:  "The ID of the Static Web App to register the Function App to as a backend.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: staticsites.ValidateStaticSiteID,
		},

		"function_app_id": {
			Description:  "The ID of a Linux or Windows Function App to connect to the Static Web App as a Backend.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
