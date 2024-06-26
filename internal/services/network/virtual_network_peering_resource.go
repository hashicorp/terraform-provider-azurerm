// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkpeerings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const virtualNetworkPeeringResourceType = "azurerm_virtual_network_peering"

func resourceVirtualNetworkPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkPeeringCreate,
		Read:   resourceVirtualNetworkPeeringRead,
		Update: resourceVirtualNetworkPeeringUpdate,
		Delete: resourceVirtualNetworkPeeringDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"virtual_network_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remote_virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
			},

			"allow_virtual_network_access": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allow_forwarded_traffic": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allow_gateway_transit": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"local_subnet_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"only_ipv6_peering_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"peer_complete_virtual_networks_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"remote_subnet_names": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"use_remote_gateways": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"triggers": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceVirtualNetworkPeeringCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualnetworkpeerings.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_network_peering", id.ID())
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
		peer.Properties.LocalSubnetNames = utils.ExpandStringSlice(v.([]interface{}))
	}

	if v, ok := d.GetOk("remote_subnet_names"); ok {
		peer.Properties.RemoteSubnetNames = utils.ExpandStringSlice(v.([]interface{}))
	}

	locks.ByID(virtualNetworkPeeringResourceType)
	defer locks.UnlockByID(virtualNetworkPeeringResourceType)

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Created"},
		Refresh: func() (interface{}, string, error) {
			future, err := client.CreateOrUpdate(ctx, id, peer, virtualnetworkpeerings.CreateOrUpdateOperationOptions{SyncRemoteAddressSpace: pointer.To(virtualnetworkpeerings.SyncRemoteAddressSpaceTrue)})
			if err != nil {
				if utils.ResponseErrorIsRetryable(err) {
					return future.HttpResponse, "Pending", err
				} else {
					if resp := future.HttpResponse; resp != nil && response.WasBadRequest(resp) && strings.Contains(err.Error(), "ReferencedResourceNotProvisioned") {
						// Resource is not yet ready, this may be the case if the Vnet was just created or another peering was just initiated.
						return future.HttpResponse, "Pending", err
					}
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

	return resourceVirtualNetworkPeeringRead(d, meta)
}

func resourceVirtualNetworkPeeringUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(virtualNetworkPeeringResourceType)
	defer locks.UnlockByID(virtualNetworkPeeringResourceType)

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("allow_forwarded_traffic") {
		existing.Model.Properties.AllowForwardedTraffic = pointer.To(d.Get("allow_forwarded_traffic").(bool))
	}
	if d.HasChange("allow_gateway_transit") {
		existing.Model.Properties.AllowGatewayTransit = pointer.To(d.Get("allow_gateway_transit").(bool))
	}
	if d.HasChange("allow_virtual_network_access") {
		existing.Model.Properties.AllowVirtualNetworkAccess = pointer.To(d.Get("allow_virtual_network_access").(bool))
	}
	if d.HasChange("local_subnet_names") {
		existing.Model.Properties.LocalSubnetNames = utils.ExpandStringSlice(d.Get("local_subnet_names").([]interface{}))
	}
	if d.HasChange("remote_subnet_names") {
		existing.Model.Properties.RemoteSubnetNames = utils.ExpandStringSlice(d.Get("remote_subnet_names").([]interface{}))
	}
	if d.HasChange("use_remote_gateways") {
		existing.Model.Properties.UseRemoteGateways = pointer.To(d.Get("use_remote_gateways").(bool))
	}
	if d.HasChange("remote_virtual_network_id") {
		existing.Model.Properties.RemoteVirtualNetwork = &virtualnetworkpeerings.SubResource{
			Id: pointer.To(d.Get("remote_virtual_network_id").(string)),
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, virtualnetworkpeerings.CreateOrUpdateOperationOptions{SyncRemoteAddressSpace: pointer.To(virtualnetworkpeerings.SyncRemoteAddressSpaceTrue)}); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceVirtualNetworkPeeringRead(d, meta)
}

func resourceVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkPeeringName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("virtual_network_name", id.VirtualNetworkName)

	if model := resp.Model; model != nil {
		if peer := model.Properties; peer != nil {
			d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
			d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
			d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
			d.Set("peer_complete_virtual_networks_enabled", pointer.From(peer.PeerCompleteVnets))
			d.Set("only_ipv6_peering_enabled", pointer.From(peer.EnableOnlyIPv6Peering))
			d.Set("local_subnet_names", pointer.From(peer.LocalSubnetNames))
			d.Set("remote_subnet_names", pointer.From(peer.RemoteSubnetNames))
			d.Set("use_remote_gateways", peer.UseRemoteGateways)

			remoteVirtualNetworkId := ""
			if network := peer.RemoteVirtualNetwork; network != nil {
				parsed, err := commonids.ParseVirtualNetworkIDInsensitively(*network.Id)
				if err != nil {
					return err
				}
				remoteVirtualNetworkId = parsed.ID()
			}
			d.Set("remote_virtual_network_id", remoteVirtualNetworkId)
		}
	}

	return nil
}

func resourceVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkPeerings
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkpeerings.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(virtualNetworkPeeringResourceType)
	defer locks.UnlockByID(virtualNetworkPeeringResourceType)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
