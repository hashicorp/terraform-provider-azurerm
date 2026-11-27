// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_interface_backend_address_pool_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func networkInterfaceBackendAddressPoolAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"network_interface_id": {
			Description:  "The ID of the Network Interface.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateNetworkInterfaceID,
		},

		"ip_configuration_name": {
			Description:  "The Name of the IP Configuration within the Network Interface which should be connected to the Backend Address Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"backend_address_pool_id": {
			Description:  "The ID of the Load Balancer Backend Address Pool which this Network Interface should be connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: loadbalancers.ValidateBackendAddressPoolID,
		},
	}
}
