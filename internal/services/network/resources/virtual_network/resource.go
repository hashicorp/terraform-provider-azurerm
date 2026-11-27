// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Virtual Network including any configured subnets. Each subnet can optionally be configured with a security group to be associated with the subnet.",
		Create:      resourceVirtualNetworkCreate,
		Read:        resourceVirtualNetworkRead,
		Update:      resourceVirtualNetworkUpdate,
		Delete:      resourceVirtualNetworkDelete,
		Importer:    pluginsdk.ImporterValidatingIdentity(&commonids.VirtualNetworkId{}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceVirtualNetworkSchema(),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.VirtualNetworkId{}),
		},
	}
}
