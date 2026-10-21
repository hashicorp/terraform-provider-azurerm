// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway_public_ip_prefix_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func natGatewayPublicIpPrefixAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"nat_gateway_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: natgateways.ValidateNatGatewayID,
		},

		"public_ip_prefix_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
		},
	}
}
