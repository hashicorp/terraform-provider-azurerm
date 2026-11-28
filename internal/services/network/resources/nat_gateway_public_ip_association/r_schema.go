// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package nat_gateway_public_ip_association

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func natGatewayPublicIpAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": {
			Description:  "The ID of the NAT Gateway.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: natgateways.ValidateNatGatewayID,
		},

		"public_ip_address_id": {
			Description:  "The ID of the Public IP Address which this NAT Gateway should be connected to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidatePublicIPAddressID,
		},
	}
}
