// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceLocalNetworkGatewayCreate,
		Read:     resourceLocalNetworkGatewayRead,
		Update:   resourceLocalNetworkGatewayUpdate,
		Delete:   resourceLocalNetworkGatewayDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&localnetworkgateways.LocalNetworkGatewayId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&localnetworkgateways.LocalNetworkGatewayId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: localNetworkGatewaySchema(),
	}
}
