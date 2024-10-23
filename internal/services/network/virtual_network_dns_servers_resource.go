// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualNetworkDnsServers() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkDnsServersCreate,
		Read:   resourceVirtualNetworkDnsServersRead,
		Update: resourceVirtualNetworkDnsServersUpdate,
		Delete: resourceVirtualNetworkDnsServersDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkDnsServersID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
			},

			"dns_servers": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceVirtualNetworkDnsServersCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetId, err := commonids.ParseVirtualNetworkID(d.Get("virtual_network_id").(string))
	if err != nil {
		return err
	}

	// This is a virtual resource so the last segment is hardcoded
	id := parse.NewVirtualNetworkDnsServersID(vnetId.SubscriptionId, vnetId.ResourceGroupName, vnetId.VirtualNetworkName, "default")

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	vnet, err := client.Get(ctx, *vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vnet.HttpResponse) {
			return fmt.Errorf("%s could not be found: %+v", vnetId, err)
		}
		return fmt.Errorf("retrieving %s: %+v", vnetId, err)
	}

	if vnet.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", vnetId)
	}
	if vnet.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", vnetId)
	}

	if vnet.Model.Properties.DhcpOptions == nil {
		vnet.Model.Properties.DhcpOptions = &virtualnetworks.DhcpOptions{}
	}

	vnet.Model.Properties.DhcpOptions.DnsServers = utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))

	if err := client.CreateOrUpdateThenPoll(ctx, *vnetId, *vnet.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, client, *vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkDnsServersRead(d, meta)
}

func resourceVirtualNetworkDnsServersRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkDnsServersID(d.Id())
	if err != nil {
		return err
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkName)

	resp, err := client.Get(ctx, vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("virtual_network_id", vnetId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("dns_servers", flattenVirtualNetworkDNSServers(props.DhcpOptions)); err != nil {
				return fmt.Errorf("setting `dns_servers`: %+v", err)
			}
		}
	}

	return nil
}

func resourceVirtualNetworkDnsServersUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetId, err := commonids.ParseVirtualNetworkID(d.Get("virtual_network_id").(string))
	if err != nil {
		return err
	}

	// This is a virtual resource so the last segment is hardcoded
	id := parse.NewVirtualNetworkDnsServersID(vnetId.SubscriptionId, vnetId.ResourceGroupName, vnetId.VirtualNetworkName, "default")

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	vnet, err := client.Get(ctx, *vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vnet.HttpResponse) {
			return fmt.Errorf("%s could not be found: %+v", vnetId, err)
		}
		return fmt.Errorf("retrieving %s: %+v", vnetId, err)
	}

	if vnet.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", vnetId)
	}
	if vnet.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", vnetId)
	}

	if vnet.Model.Properties.DhcpOptions == nil {
		vnet.Model.Properties.DhcpOptions = &virtualnetworks.DhcpOptions{}
	}

	if d.HasChange("dns_servers") {
		vnet.Model.Properties.DhcpOptions.DnsServers = utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *vnetId, *vnet.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, client, *vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkDnsServersRead(d, meta)
}

func resourceVirtualNetworkDnsServersDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkDnsServersID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkName)

	vnet, err := client.Get(ctx, vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vnet.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing %s from state", vnetId.ID(), id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", vnetId, err)
	}

	if vnet.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", vnetId)
	}
	if vnet.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", vnetId)
	}

	if vnet.Model.Properties.DhcpOptions == nil {
		log.Printf("[INFO] dhcpOptions for %s was nil, dnsServers already deleted - removing %s from state", vnetId.ID(), id)
		return nil
	}

	vnet.Model.Properties.DhcpOptions.DnsServers = utils.ExpandStringSlice(make([]interface{}, 0))

	if err := client.CreateOrUpdateThenPoll(ctx, vnetId, *vnet.Model); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
