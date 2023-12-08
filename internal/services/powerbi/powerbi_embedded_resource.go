// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package powerbi

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePowerBIEmbedded() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePowerBIEmbeddedCreate,
		Read:   resourcePowerBIEmbeddedRead,
		Update: resourcePowerBIEmbeddedUpdate,
		Delete: resourcePowerBIEmbeddedDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := capacities.ParseCapacityID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EmbeddedName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A1",
					"A2",
					"A3",
					"A4",
					"A5",
					"A6",
				}, false),
			},

			"administrators": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.EmbeddedAdministratorName,
				},
			},

			"mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(capacities.ModeGenOne),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(capacities.ModeGenOne),
					string(capacities.ModeGenTwo),
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePowerBIEmbeddedCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := capacities.NewCapacityID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.GetDetails(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_powerbi_embedded", id.ID())
	}

	administrators := d.Get("administrators").(*pluginsdk.Set).List()
	mode := capacities.Mode(d.Get("mode").(string))

	parameters := capacities.DedicatedCapacity{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Properties: &capacities.DedicatedCapacityProperties{
			Administration: &capacities.DedicatedCapacityAdministrators{
				Members: utils.ExpandStringSlice(administrators),
			},
			Mode: &mode,
		},
		Sku: capacities.CapacitySku{
			Name: d.Get("sku_name").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePowerBIEmbeddedRead(d, meta)
}

func resourcePowerBIEmbeddedRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacities.ParseCapacityID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetDetails(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CapacityName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			var adminMembers *[]string
			if props.Administration != nil {
				adminMembers = props.Administration.Members
			}
			if err := d.Set("administrators", utils.FlattenStringSlice(adminMembers)); err != nil {
				return fmt.Errorf("setting `administration`: %+v", err)
			}

			mode := ""
			if props.Mode != nil {
				mode = string(*props.Mode)
			}
			d.Set("mode", mode)
		}

		d.Set("sku_name", model.Sku.Name)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourcePowerBIEmbeddedUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacities.ParseCapacityID(d.Id())
	if err != nil {
		return err
	}

	parameters := capacities.DedicatedCapacityUpdateParameters{}

	if d.HasChange("administrators") || d.HasChange("mode") {
		administrators := d.Get("administrators").(*pluginsdk.Set).List()
		mode := capacities.Mode(d.Get("mode").(string))

		parameters.Properties = &capacities.DedicatedCapacityMutableProperties{
			Administration: &capacities.DedicatedCapacityAdministrators{
				Members: utils.ExpandStringSlice(administrators),
			},
			Mode: &mode,
		}
	}

	if d.HasChange("sku_name") {
		parameters.Sku = &capacities.CapacitySku{
			Name: d.Get("sku_name").(string),
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourcePowerBIEmbeddedRead(d, meta)
}

func resourcePowerBIEmbeddedDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacities.ParseCapacityID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
