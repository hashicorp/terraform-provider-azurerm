// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_peering

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkPeeringCreate,
		Read:   resourceVirtualNetworkPeeringRead,
		Update: resourceVirtualNetworkPeeringUpdate,
		Delete: resourceVirtualNetworkPeeringDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&virtualnetworkpeerings.VirtualNetworkPeeringId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualnetworkpeerings.VirtualNetworkPeeringId{}),
		},
		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: virtualNetworkPeeringSchema(),
	}
}
