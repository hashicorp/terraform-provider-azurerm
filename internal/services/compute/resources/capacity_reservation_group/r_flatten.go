// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package capacity_reservation_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceCapacityReservationGroupFlatten(d *pluginsdk.ResourceData, id *capacityreservationgroups.CapacityReservationGroupId, model *capacityreservationgroups.CapacityReservationGroup) error {
	d.Set("name", id.CapacityReservationGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", pluginsdk.FlattenSlice(model.Zones))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
