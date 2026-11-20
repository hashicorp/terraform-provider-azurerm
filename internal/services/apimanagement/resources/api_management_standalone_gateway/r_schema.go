// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_standalone_gateway

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "Specifies the name which should be used for this API Management Standalone Gateway.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z](?:[a-zA-Z0-9-]{0,43}[a-zA-Z0-9])?$"),
				"The `name` must be between 1 and 45 characters long and can only include letters, numbers, and hyphens. The first character must be a letter and last character must be a letter or a number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the SKU.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							// `Standard`: Deprecated - will be removed in a future version.
							// `WorkspaceGatewayStandard`: Private Preview – not currently for general use.
							string(apigateway.ApiGatewaySkuTypeWorkspaceGatewayPremium),
						}, false),
					},
					"capacity": {
						Description:  "The number of deployed units of the SKU.",
						Type:         pluginsdk.TypeInt,
						Default:      1,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"backend_subnet_id": {
			Description:  "Specifies the subnet ID in which the backend systems are hosted.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
			RequiredWith: []string{"virtual_network_type"},
		},

		"tags": commonschema.TagsForceNew(),

		"virtual_network_type": {
			Description: "Specifies the type of VPN in which API Management gateway needs to be configured.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				// Note: Whilst the `None` value exists it's handled in the Create/Update and Read functions.
				string(apigateway.VirtualNetworkTypeExternal),
				string(apigateway.VirtualNetworkTypeInternal),
			}, false),
			RequiredWith: []string{"backend_subnet_id"},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
