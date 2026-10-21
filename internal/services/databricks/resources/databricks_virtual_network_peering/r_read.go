// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks_virtual_network_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/vnetpeering"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Databricks %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("name", id.VirtualNetworkPeeringName)
	d.Set("workspace_id", vnetpeering.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())

	if model := resp.Model; model != nil {
		d.Set("allow_virtual_network_access", model.Properties.AllowVirtualNetworkAccess)
		d.Set("allow_forwarded_traffic", model.Properties.AllowForwardedTraffic)
		d.Set("allow_gateway_transit", model.Properties.AllowGatewayTransit)
		d.Set("use_remote_gateways", model.Properties.UseRemoteGateways)

		addressSpacePrefixes := make([]string, 0)
		if model.Properties.DatabricksAddressSpace != nil && model.Properties.DatabricksAddressSpace.AddressPrefixes != nil {
			addressSpacePrefixes = *model.Properties.DatabricksAddressSpace.AddressPrefixes
		}
		d.Set("address_space_prefixes", addressSpacePrefixes)

		remoteAddressSpacePrefixes := make([]string, 0)
		if model.Properties.RemoteAddressSpace != nil && model.Properties.RemoteAddressSpace.AddressPrefixes != nil {
			remoteAddressSpacePrefixes = *model.Properties.RemoteAddressSpace.AddressPrefixes
		}
		d.Set("remote_address_space_prefixes", remoteAddressSpacePrefixes)

		databricksVirtualNetworkId := ""
		if model.Properties.DatabricksVirtualNetwork != nil && model.Properties.DatabricksVirtualNetwork.Id != nil {
			databricksVirtualNetworkId = *model.Properties.DatabricksVirtualNetwork.Id
		}
		d.Set("virtual_network_id", databricksVirtualNetworkId)

		d.Set("remote_virtual_network_id", pointer.From(model.Properties.RemoteVirtualNetwork.Id))
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

const databricksVnetPeeringsResourceType string = "azurerm_databricks_virtual_network_peering"
