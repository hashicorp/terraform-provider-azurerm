// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCapacityReservationCreate(d *pluginsdk.ResourceData, meta any) error {
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

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, capacityreservations.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_capacity_reservation", id.ID())
		}
	}

	payload := capacityreservations.CapacityReservation{
		Location: location.Normalize(capacityReservationGroup.Model.Location),
		Sku:      expandCapacityReservationSku(d.Get("sku").([]any)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}
	if v, ok := d.GetOk("zone"); ok {
		payload.Zones = &[]string{
			v.(string),
		}
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceCapacityReservationRead(d, meta)
}

func expandCapacityReservationSku(input []any) capacityreservations.Sku {
	v := input[0].(map[string]any)
	return capacityreservations.Sku{
		Name:     pointer.To(v["name"].(string)),
		Capacity: pointer.To(int64(v["capacity"].(int))),
	}
}
