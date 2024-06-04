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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/virtualwans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/securitypartnerproviders"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubSecurityPartnerProvider() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubSecurityPartnerProviderCreate,
		Read:   resourceVirtualHubSecurityPartnerProviderRead,
		Update: resourceVirtualHubSecurityPartnerProviderUpdate,
		Delete: resourceVirtualHubSecurityPartnerProviderDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securitypartnerproviders.ParseSecurityPartnerProviderID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"security_provider_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securitypartnerproviders.SecurityProviderNameZScaler),
					string(securitypartnerproviders.SecurityProviderNameIBoss),
					string(securitypartnerproviders.SecurityProviderNameCheckpoint),
				}, false),
			},

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceVirtualHubSecurityPartnerProviderCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviders
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := securitypartnerproviders.NewSecurityPartnerProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_hub_security_partner_provider", id.ID())
	}

	parameters := securitypartnerproviders.SecurityPartnerProvider{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &securitypartnerproviders.SecurityPartnerProviderPropertiesFormat{
			SecurityProviderName: pointer.To(securitypartnerproviders.SecurityProviderName(d.Get("security_provider_name").(string))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("virtual_hub_id"); ok {
		parameters.Properties.VirtualHub = &securitypartnerproviders.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubSecurityPartnerProviderRead(d, meta)
}

func resourceVirtualHubSecurityPartnerProviderRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviders
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypartnerproviders.ParseSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SecurityPartnerProviderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			d.Set("security_provider_name", string(pointer.From(props.SecurityProviderName)))

			if props.VirtualHub != nil && props.VirtualHub.Id != nil {
				d.Set("virtual_hub_id", props.VirtualHub.Id)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceVirtualHubSecurityPartnerProviderUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviders
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypartnerproviders.ParseSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	parameters := securitypartnerproviders.TagsObject{}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.UpdateTags(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualHubSecurityPartnerProviderRead(d, meta)
}

func resourceVirtualHubSecurityPartnerProviderDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SecurityPartnerProviders
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypartnerproviders.ParseSecurityPartnerProviderID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
