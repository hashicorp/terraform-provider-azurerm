// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func frontdoorCustomHttpsConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"frontend_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(frontdoors.ParseFrontendEndpointIDInsensitively),
		},

		"custom_https_provisioning_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"custom_https_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: schemaCustomHttpsConfiguration(),
			},
		},
	}
}
