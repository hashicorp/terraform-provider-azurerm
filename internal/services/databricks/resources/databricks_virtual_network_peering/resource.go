// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_virtual_network_peering

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/vnetpeering"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name databricks_virtual_network_peering -service-package-name databricks -properties "name" -compare-values "subscription_id:workspace_id,resource_group_name:workspace_id,workspace_name:workspace_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabricksVirtualNetworkPeeringCreate,
		Read:   resourceDatabricksVirtualNetworkPeeringRead,
		Update: resourceDatabricksVirtualNetworkPeeringUpdate,
		Delete: resourceDatabricksVirtualNetworkPeeringDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&vnetpeering.VirtualNetworkPeeringId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&vnetpeering.VirtualNetworkPeeringId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: databricksVirtualNetworkPeeringSchema(),
	}
}
