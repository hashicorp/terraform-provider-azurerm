// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubIP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubIPCreate,
		Read:   resourceVirtualHubIPRead,
		Update: resourceVirtualHubIPUpdate,
		Delete: resourceVirtualHubIPDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseVirtualHubIPConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"public_ip_address_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidatePublicIPAddressID,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"private_ip_address": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"private_ip_allocation_method": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  virtualwans.IPAllocationMethodDynamic,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualwans.IPAllocationMethodDynamic),
					string(virtualwans.IPAllocationMethodStatic),
				}, false),
			},
		},
	}
}

func resourceVirtualHubIPCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtualHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubId.VirtualHubName, virtualHubResourceName)

	id := commonids.NewVirtualHubIPConfigurationID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, d.Get("name").(string))

	existing, err := client.VirtualHubIPConfigurationGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_hub_ip", id.ID())
	}

	parameters := virtualwans.HubIPConfiguration{
		Name: pointer.To(id.IpConfigName),
		Properties: &virtualwans.HubIPConfigurationPropertiesFormat{
			Subnet: &virtualwans.Subnet{
				Id: pointer.To(d.Get("subnet_id").(string)),
			},
		},
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		parameters.Properties.PrivateIPAddress = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("private_ip_allocation_method"); ok {
		parameters.Properties.PrivateIPAllocationMethod = pointer.To(virtualwans.IPAllocationMethod(v.(string)))
	}

	if v, ok := d.GetOk("public_ip_address_id"); ok {
		parameters.Properties.PublicIPAddress = &virtualwans.PublicIPAddress{
			Id: pointer.To(v.(string)),
		}
	}

	if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubIPRead(d, meta)
}

func resourceVirtualHubIPUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtualHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubId.VirtualHubName, virtualHubResourceName)

	id, err := commonids.ParseVirtualHubIPConfigurationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VirtualHubIPConfigurationGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("private_ip_address") {
		payload.Properties.PrivateIPAddress = pointer.To(d.Get("private_ip_address").(string))
	}

	if d.HasChange("private_ip_allocation_method") {
		payload.Properties.PrivateIPAllocationMethod = pointer.To(virtualwans.IPAllocationMethod(d.Get("private_ip_allocation_method").(string)))
	}

	if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubIPRead(d, meta)
}

func resourceVirtualHubIPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubIPConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualHubIPConfigurationGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Virtual Hub IP %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.IpConfigName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("private_ip_address", props.PrivateIPAddress)
			d.Set("private_ip_allocation_method", string(pointer.From(props.PrivateIPAllocationMethod)))

			if v := props.PublicIPAddress; v != nil {
				d.Set("public_ip_address_id", v.Id)
			}

			if v := props.Subnet; v != nil {
				d.Set("subnet_id", v.Id)
			}
		}
	}

	return nil
}

func resourceVirtualHubIPDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubIPConfigurationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	if err := client.VirtualHubIPConfigurationDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
