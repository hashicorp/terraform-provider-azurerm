// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_peering

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkPeeringCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworkpeerings.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_network_peering", id.ID())
		}
	}

	peer := virtualnetworkpeerings.VirtualNetworkPeering{
		Properties: &virtualnetworkpeerings.VirtualNetworkPeeringPropertiesFormat{
			AllowVirtualNetworkAccess: pointer.To(d.Get("allow_virtual_network_access").(bool)),
			AllowForwardedTraffic:     pointer.To(d.Get("allow_forwarded_traffic").(bool)),
			AllowGatewayTransit:       pointer.To(d.Get("allow_gateway_transit").(bool)),
			PeerCompleteVnets:         pointer.To(d.Get("peer_complete_virtual_networks_enabled").(bool)),
			UseRemoteGateways:         pointer.To(d.Get("use_remote_gateways").(bool)),
			RemoteVirtualNetwork: &virtualnetworkpeerings.SubResource{
				Id: pointer.To(d.Get("remote_virtual_network_id").(string)),
			},
		},
	}

	if v, ok := d.GetOk("only_ipv6_peering_enabled"); ok {
		peer.Properties.EnableOnlyIPv6Peering = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("local_subnet_names"); ok {
		peer.Properties.LocalSubnetNames = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if v, ok := d.GetOk("remote_subnet_names"); ok {
		peer.Properties.RemoteSubnetNames = pluginsdk.ExpandStringSlice(v.([]any))
	}

	vnetId := commonids.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string))
	remoteVnetId, err := commonids.ParseVirtualNetworkID(d.Get("remote_virtual_network_id").(string))
	if err != nil {
		return err
	}

	vnetIDsToLock := []string{vnetId.ID(), remoteVnetId.ID()}
	locks.MultipleByID(&vnetIDsToLock)
	defer locks.UnlockMultipleByID(&vnetIDsToLock)

	// TODO: implement `CallbackThenPoll`, rework to remove StateChangeConf
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Created"},
		Refresh: func() (any, string, error) {
			future, err := client.CreateOrUpdate(ctx, id, peer, virtualnetworkpeerings.CreateOrUpdateOperationOptions{SyncRemoteAddressSpace: pointer.To(virtualnetworkpeerings.SyncRemoteAddressSpaceTrue)})
			if err != nil {
				if resp := future.HttpResponse; resp != nil && response.WasBadRequest(resp) && strings.Contains(err.Error(), "ReferencedResourceNotProvisioned") {
					// Resource is not yet ready, this may be the case if the Vnet was just created or another peering was just initiated.
					return future.HttpResponse, "Pending", err
				}

				return future.HttpResponse, "", err
			}

			if err = future.Poller.PollUntilDone(ctx); err != nil {
				return future.HttpResponse, "", err
			}

			return future.HttpResponse, "Created", nil
		},
		Timeout: time.Until(deadline),
		Delay:   15 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be created: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceVirtualNetworkPeeringRead(d, meta)
}
