// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualWan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualWanCreate,
		Read:   resourceVirtualWanRead,
		Update: resourceVirtualWanUpdate,
		Delete: resourceVirtualWanDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVirtualWANID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"disable_vpn_encryption": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allow_branch_to_branch_traffic": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"office365_local_breakout_category": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualwans.OfficeTrafficCategoryAll),
					string(virtualwans.OfficeTrafficCategoryNone),
					string(virtualwans.OfficeTrafficCategoryOptimize),
					string(virtualwans.OfficeTrafficCategoryOptimizeAndAllow),
				}, false),
				Default: string(virtualwans.OfficeTrafficCategoryNone),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Standard",
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVirtualWanCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := virtualwans.NewVirtualWANID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.VirtualWansGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_wan", id.ID())
	}

	wan := virtualwans.VirtualWAN{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &virtualwans.VirtualWanProperties{
			DisableVpnEncryption:           pointer.To(d.Get("disable_vpn_encryption").(bool)),
			AllowBranchToBranchTraffic:     pointer.To(d.Get("allow_branch_to_branch_traffic").(bool)),
			Office365LocalBreakoutCategory: pointer.To(virtualwans.OfficeTrafficCategory(d.Get("office365_local_breakout_category").(string))),
			Type:                           pointer.To(d.Get("type").(string)),
		},
	}

	if err := client.VirtualWansCreateOrUpdateThenPoll(ctx, id, wan); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualWanRead(d, meta)
}

func resourceVirtualWanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualWANID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VirtualWansGet(ctx, *id)
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

	if d.HasChange("disable_vpn_encryption") {
		payload.Properties.DisableVpnEncryption = pointer.To(d.Get("disable_vpn_encryption").(bool))
	}

	if d.HasChange("allow_branch_to_branch_traffic") {
		payload.Properties.AllowBranchToBranchTraffic = pointer.To(d.Get("allow_branch_to_branch_traffic").(bool))
	}

	if d.HasChange("office365_local_breakout_category") {
		payload.Properties.Office365LocalBreakoutCategory = pointer.To(virtualwans.OfficeTrafficCategory(d.Get("office365_local_breakout_category").(string)))
	}

	if d.HasChange("type") {
		payload.Properties.Type = pointer.To(d.Get("type").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.VirtualWansCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualWanRead(d, meta)
}

func resourceVirtualWanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualWANID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualWansGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualWanName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
			d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
			d.Set("office365_local_breakout_category", pointer.From(props.Office365LocalBreakoutCategory))
			d.Set("type", props.Type)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceVirtualWanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualWANID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VirtualWansDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
