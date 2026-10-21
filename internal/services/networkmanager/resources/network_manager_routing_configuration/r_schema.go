// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_configuration

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagerroutingconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"network_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&networkmanagerroutingconfigurations.NetworkManagerId{}),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"route_table_usage_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(networkmanagerroutingconfigurations.RouteTableUsageModeManagedOnly),
			ValidateFunc: validation.StringInSlice(networkmanagerroutingconfigurations.PossibleValuesForRouteTableUsageMode(), false),
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
