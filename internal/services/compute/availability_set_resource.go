// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/availabilitysets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAvailabilitySet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAvailabilitySetCreateUpdate,
		Read:   resourceAvailabilitySetRead,
		Update: resourceAvailabilitySetCreateUpdate,
		Delete: resourceAvailabilitySetDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseAvailabilitySetID(id)
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
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,78}[a-zA-Z0-9_])?$"),
					"The Availability set name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 80 characters, and it must begin a letter or number and end with a letter, number or underscore.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"platform_update_domain_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      5,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 20),
			},

			"platform_fault_domain_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      3,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"managed": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"proximity_placement_group_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,

				// We have to ignore case due to incorrect capitalisation of resource group name in
				// proximity placement group ID in the response we get from the API request
				//
				// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceAvailabilitySetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewAvailabilitySetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_availability_set", id.ID())
		}
	}

	updateDomainCount := d.Get("platform_update_domain_count").(int)
	faultDomainCount := d.Get("platform_fault_domain_count").(int)
	managed := d.Get("managed").(bool)
	t := d.Get("tags").(map[string]interface{})

	payload := availabilitysets.AvailabilitySet{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &availabilitysets.AvailabilitySetProperties{
			PlatformFaultDomainCount:  utils.Int64(int64(faultDomainCount)),
			PlatformUpdateDomainCount: utils.Int64(int64(updateDomainCount)),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("proximity_placement_group_id"); ok {
		payload.Properties.ProximityPlacementGroup = &availabilitysets.SubResource{
			Id: utils.String(v.(string)),
		}
	}

	if managed {
		n := "Aligned"
		payload.Sku = &availabilitysets.Sku{
			Name: &n,
		}
	}

	_, err := client.CreateOrUpdate(ctx, id, payload)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAvailabilitySetRead(d, meta)
}

func resourceAvailabilitySetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseAvailabilitySetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AvailabilitySetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		managed := false
		if model.Sku != nil && model.Sku.Name != nil {
			managed = strings.EqualFold(*model.Sku.Name, "Aligned")
		}
		d.Set("managed", managed)

		if props := model.Properties; props != nil {
			d.Set("platform_update_domain_count", props.PlatformUpdateDomainCount)
			d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)

			if proximityPlacementGroup := props.ProximityPlacementGroup; proximityPlacementGroup != nil {
				d.Set("proximity_placement_group_id", proximityPlacementGroup.Id)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceAvailabilitySetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.AvailabilitySetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseAvailabilitySetID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
