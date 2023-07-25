// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCapacityReservationGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCapacityReservationGroupCreate,
		Read:   resourceCapacityReservationGroupRead,
		Update: resourceCapacityReservationGroupUpdate,
		Delete: resourceCapacityReservationGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := capacityreservationgroups.ParseCapacityReservationGroupID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CapacityReservationGroupName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceCapacityReservationGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := capacityreservationgroups.NewCapacityReservationGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id, capacityreservationgroups.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_capacity_reservation_group", id.ID())
	}

	parameters := capacityreservationgroups.CapacityReservationGroup{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCapacityReservationGroupRead(d, meta)
}

func resourceCapacityReservationGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservationgroups.ParseCapacityReservationGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, capacityreservationgroups.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CapacityReservationGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", utils.FlattenStringSlice(model.Zones))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceCapacityReservationGroupUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservationgroups.ParseCapacityReservationGroupID(d.Id())
	if err != nil {
		return err
	}

	parameters := capacityreservationgroups.CapacityReservationGroupUpdate{}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCapacityReservationGroupRead(d, meta)
}

func resourceCapacityReservationGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.CapacityReservationGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacityreservationgroups.ParseCapacityReservationGroupID(d.Id())
	if err != nil {
		return err
	}

	// It takes several seconds to sync the cache of reservations list in Capacity Reservation Group. Delete operation requires the list to be empty, and fails before the cache sync is completed.
	// Retry the delete operation after a minute as a workaround. Issue is tracked by: https://github.com/Azure/azure-rest-api-specs/issues/18767
	if _, err := client.Delete(ctx, *id); err != nil {
		stateConf := &pluginsdk.StateChangeConf{
			Pending: []string{"Deleting"},
			Target:  []string{"Deleted"},
			Refresh: func() (interface{}, string, error) {
				res, err := client.Delete(ctx, *id)
				if err != nil {
					return res, "Deleting", nil
				}
				return res, "Deleted", nil
			},
			MinTimeout: 15 * time.Second,
			Timeout:    1 * time.Minute,
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
		}
	}

	return nil
}
