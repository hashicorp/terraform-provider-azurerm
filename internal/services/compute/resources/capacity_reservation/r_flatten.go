// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Note: when adding new properties, ensure the `(CapacityReservationListResource).capacityReservationToCapacityReservations` function is updated as well.
func resourceCapacityReservationFlatten(d *pluginsdk.ResourceData, id *capacityreservations.CapacityReservationId, model *capacityreservations.CapacityReservation) error {
	d.Set("name", id.CapacityReservationName)
	groupId := capacityreservationgroups.NewCapacityReservationGroupID(id.SubscriptionId, id.ResourceGroupName, id.CapacityReservationGroupName)
	d.Set("capacity_reservation_group_id", groupId.ID())

	if model != nil {
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

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenCapacityReservationSku(input capacityreservations.Sku) []any {
	return []any{
		map[string]any{
			"name":     pointer.From(input.Name),
			"capacity": pointer.From(input.Capacity),
		},
	}
}
