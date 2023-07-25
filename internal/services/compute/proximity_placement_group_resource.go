// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceProximityPlacementGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceProximityPlacementGroupCreateUpdate,
		Read:   resourceProximityPlacementGroupRead,
		Update: resourceProximityPlacementGroupCreateUpdate,
		Delete: resourceProximityPlacementGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := proximityplacementgroups.ParseProximityPlacementGroupID(id)
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

			"allowed_vm_sizes": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"allowed_vm_sizes"},
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("allowed_vm_sizes", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.(*pluginsdk.Set).List()) > 0 && len(new.(*pluginsdk.Set).List()) == 0
			}),
		),
	}
}

func resourceProximityPlacementGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := proximityplacementgroups.NewProximityPlacementGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, proximityplacementgroups.DefaultGetOperationOptions())
	if d.IsNewResource() {
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_proximity_placement_group", id.ID())
		}
	}

	payload := proximityplacementgroups.ProximityPlacementGroup{
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &proximityplacementgroups.ProximityPlacementGroupProperties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("allowed_vm_sizes"); ok {
		if payload.Properties.Intent == nil {
			payload.Properties.Intent = &proximityplacementgroups.ProximityPlacementGroupPropertiesIntent{}
		}
		payload.Properties.Intent.VMSizes = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("zone"); ok {
		zones := zones.Expand([]string{v.(string)})
		payload.Zones = &zones
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceProximityPlacementGroupRead(d, meta)
}

func resourceProximityPlacementGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := proximityplacementgroups.ParseProximityPlacementGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, proximityplacementgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", id.ProximityPlacementGroupName)
		d.Set("resource_group_name", id.ResourceGroupName)

		d.Set("location", location.Normalize(model.Location))

		intentVmSizes := make([]string, 0)
		if props := model.Properties; props != nil {
			if intent := props.Intent; intent != nil {
				if intent.VMSizes != nil {
					intentVmSizes = *intent.VMSizes
				}
			}
		}
		d.Set("allowed_vm_sizes", intentVmSizes)

		zone := ""
		if v := zones.Flatten(model.Zones); len(v) != 0 {
			zone = v[0]
		}
		d.Set("zone", zone)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceProximityPlacementGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.ProximityPlacementGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := proximityplacementgroups.ParseProximityPlacementGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
