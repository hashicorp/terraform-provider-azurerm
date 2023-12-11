// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCapacityReservation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCapacityReservationCreate,
		Read:   resourceCapacityReservationRead,
		Update: resourceCapacityReservationUpdate,
		Delete: resourceCapacityReservationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := capacityreservations.ParseCapacityReservationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CapacityReservationName(),
			},

			"capacity_reservation_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
			},

			"zone": commonschema.ZoneSingleOptionalForceNew(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceCapacityReservationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	groupsClient := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	capacityReservationGroupId, err := capacityreservationgroups.ParseCapacityReservationGroupID(d.Get("capacity_reservation_group_id").(string))
	if err != nil {
		return err
	}

	capacityReservationGroup, err := groupsClient.Get(ctx, *capacityReservationGroupId, capacityreservationgroups.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *capacityReservationGroupId, err)
	}
	if capacityReservationGroup.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *capacityReservationGroupId)
	}

	id := capacityreservations.NewCapacityReservationID(subscriptionId, capacityReservationGroupId.ResourceGroupName, capacityReservationGroupId.CapacityReservationGroupName, d.Get("name").(string))
	existing, err := client.Get(ctx, id, capacityreservations.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_capacity_reservation", id.ID())
	}

	payload := capacityreservations.CapacityReservation{
		Location: location.Normalize(capacityReservationGroup.Model.Location),
		Sku:      expandCapacityReservationSku(d.Get("sku").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if v, ok := d.GetOk("zone"); ok {
		payload.Zones = &[]string{
			v.(string),
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCapacityReservationRead(d, meta)
}

func resourceCapacityReservationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservations.ParseCapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, capacityreservations.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CapacityReservationName)
	groupId := capacityreservationgroups.NewCapacityReservationGroupID(id.SubscriptionId, id.ResourceGroupName, id.CapacityReservationGroupName)
	d.Set("capacity_reservation_group_id", groupId.ID())

	if model := resp.Model; model != nil {
		if err := d.Set("sku", flattenCapacityReservationSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		zone := ""
		if model.Zones != nil && len(*model.Zones) > 0 {
			z := *model.Zones
			zone = z[0]
		}
		d.Set("zone", zone)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceCapacityReservationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservations.ParseCapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	payload := capacityreservations.CapacityReservationUpdate{}
	if d.HasChange("sku") {
		payload.Sku = pointer.To(expandCapacityReservationSku(d.Get("sku").([]interface{})))
	}
	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCapacityReservationRead(d, meta)
}

func resourceCapacityReservationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservations.ParseCapacityReservationID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandCapacityReservationSku(input []interface{}) capacityreservations.Sku {
	v := input[0].(map[string]interface{})
	return capacityreservations.Sku{
		Name:     utils.String(v["name"].(string)),
		Capacity: utils.Int64(int64(v["capacity"].(int))),
	}
}

func flattenCapacityReservationSku(input capacityreservations.Sku) []interface{} {
	var name string
	if input.Name != nil {
		name = *input.Name
	}

	var capacity int64
	if input.Capacity != nil {
		capacity = *input.Capacity
	}

	return []interface{}{
		map[string]interface{}{
			"name":     name,
			"capacity": capacity,
		},
	}
}
