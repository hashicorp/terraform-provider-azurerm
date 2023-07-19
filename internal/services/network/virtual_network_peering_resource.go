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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

const virtualNetworkPeeringResourceType = "azurerm_virtual_network_peering"

func resourceVirtualNetworkPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkPeeringCreate,
		Read:   resourceVirtualNetworkPeeringRead,
		Update: resourceVirtualNetworkPeeringUpdate,
		Delete: resourceVirtualNetworkPeeringDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkPeeringID(id)
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
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_virtual_network_peering", id.ID())
	}

	peer := network.VirtualNetworkPeering{
		VirtualNetworkPeeringPropertiesFormat: &network.VirtualNetworkPeeringPropertiesFormat{
			AllowVirtualNetworkAccess: pointer.To(d.Get("allow_virtual_network_access").(bool)),
			AllowForwardedTraffic:     pointer.To(d.Get("allow_forwarded_traffic").(bool)),
			AllowGatewayTransit:       pointer.To(d.Get("allow_gateway_transit").(bool)),
			UseRemoteGateways:         pointer.To(d.Get("use_remote_gateways").(bool)),
			RemoteVirtualNetwork: &network.SubResource{
				ID: pointer.To(d.Get("remote_virtual_network_id").(string)),
			},
		},
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
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name, peer, network.SyncRemoteAddressSpaceTrue)
			if err != nil {
				if utils.ResponseErrorIsRetryable(err) {
					return future.Response(), "Pending", err
				} else {
					if resp := future.Response(); resp != nil && response.WasBadRequest(resp) && strings.Contains(err.Error(), "ReferencedResourceNotProvisioned") {
						// Resource is not yet ready, this may be the case if the Vnet was just created or another peering was just initiated.
						return future.Response(), "Pending", err
					}
				}

				return future.Response(), "", err
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return future.Response(), "", err
			}

			return future.Response(), "Created", nil
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
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(virtualNetworkPeeringResourceType)
	defer locks.UnlockByID(virtualNetworkPeeringResourceType)

	existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if d.HasChange("allow_forwarded_traffic") {
		existing.VirtualNetworkPeeringPropertiesFormat.AllowForwardedTraffic = pointer.To(d.Get("allow_forwarded_traffic").(bool))
	}
	if d.HasChange("allow_gateway_transit") {
		existing.VirtualNetworkPeeringPropertiesFormat.AllowGatewayTransit = pointer.To(d.Get("allow_gateway_transit").(bool))
	}
	if d.HasChange("allow_virtual_network_access") {
		existing.VirtualNetworkPeeringPropertiesFormat.AllowVirtualNetworkAccess = pointer.To(d.Get("allow_virtual_network_access").(bool))
	}
	if d.HasChange("use_remote_gateways") {
		existing.VirtualNetworkPeeringPropertiesFormat.UseRemoteGateways = pointer.To(d.Get("use_remote_gateways").(bool))
	}
	if d.HasChange("remote_virtual_network_id") {
		existing.VirtualNetworkPeeringPropertiesFormat.RemoteVirtualNetwork = &network.SubResource{
			ID: pointer.To(d.Get("remote_virtual_network_id").(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name, existing, network.SyncRemoteAddressSpaceTrue)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceVirtualNetworkPeeringRead(d, meta)
}

func resourceVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("virtual_network_name", id.VirtualNetworkName)

	if peer := resp.VirtualNetworkPeeringPropertiesFormat; peer != nil {
		d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
		d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
		d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
		d.Set("use_remote_gateways", peer.UseRemoteGateways)

		remoteVirtualNetworkId := ""
		if network := peer.RemoteVirtualNetwork; network != nil {
			parsed, err := commonids.ParseVirtualNetworkIDInsensitively(*network.ID)
			if err != nil {
				return fmt.Errorf("parsing %q as a Virtual Network ID: %+v", *network.ID, err)
			}
			remoteVirtualNetworkId = parsed.ID()
		}
		d.Set("remote_virtual_network_id", remoteVirtualNetworkId)
	}

	return nil
}

func resourceVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(virtualNetworkPeeringResourceType)
	defer locks.UnlockByID(virtualNetworkPeeringResourceType)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return err
}
