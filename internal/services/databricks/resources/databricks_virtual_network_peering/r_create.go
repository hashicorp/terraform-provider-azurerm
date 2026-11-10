// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package databricks_virtual_network_peering

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/vnetpeering"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksVirtualNetworkPeeringCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var id vnetpeering.VirtualNetworkPeeringId

	// I need to include the workspace ID in the properties because I need the name
	// of the workspace to create the peerings ID
	workspaceId, err := vnetpeering.ParseWorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("unable to parse 'workspace_id': %+v", err)
	}

	id = vnetpeering.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), workspaceId.WorkspaceName, d.Get("name").(string))

	locks.ByID(databricksVnetPeeringsResourceType)
	defer locks.UnlockByID(databricksVnetPeeringsResourceType)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Databricks %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databricks_virtual_network_peering", id.ID())
		}
	}

	allowForwardedTraffic := d.Get("allow_forwarded_traffic").(bool)
	allowGatewayTransit := d.Get("allow_gateway_transit").(bool)
	allowVirtualNetworkAccess := d.Get("allow_virtual_network_access").(bool)
	useRemoteGateways := d.Get("use_remote_gateways").(bool)
	remoteVirtualNetwork := d.Get("remote_virtual_network_id").(string)
	databricksAddressSpace := pluginsdk.ExpandStringSlice(d.Get("address_space_prefixes").([]any))
	remoteAddressSpace := pluginsdk.ExpandStringSlice(d.Get("remote_address_space_prefixes").([]any))

	props := vnetpeering.VirtualNetworkPeeringPropertiesFormat{
		DatabricksAddressSpace: &vnetpeering.AddressSpace{
			AddressPrefixes: databricksAddressSpace,
		},
		// The RP always creates the same vNet ID for the Databricks internal vNet in the below format:
		// '/subscriptions/{subscription}/resourceGroups/{group1}/providers/Microsoft.Network/virtualNetworks/workers-vnet'
		DatabricksVirtualNetwork: &vnetpeering.VirtualNetworkPeeringPropertiesFormatDatabricksVirtualNetwork{
			Id: pointer.To(commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, "workers-vnet").ID()),
		},
		RemoteAddressSpace: &vnetpeering.AddressSpace{
			AddressPrefixes: remoteAddressSpace,
		},
		RemoteVirtualNetwork: vnetpeering.VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork{
			Id: pointer.To(remoteVirtualNetwork),
		},
		AllowForwardedTraffic:     &allowForwardedTraffic,
		AllowGatewayTransit:       &allowGatewayTransit,
		AllowVirtualNetworkAccess: &allowVirtualNetworkAccess,
		UseRemoteGateways:         &useRemoteGateways,
	}

	peer := vnetpeering.VirtualNetworkPeering{
		Name:       &id.VirtualNetworkPeeringName,
		Properties: props,
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, peer, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating Databricks %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceDatabricksVirtualNetworkPeeringRead(d, meta)
}
