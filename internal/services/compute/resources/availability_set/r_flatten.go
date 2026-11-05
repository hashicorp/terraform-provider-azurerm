// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package availability_set

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/availabilitysets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceAvailabilitySetFlatten(d *pluginsdk.ResourceData, id *commonids.AvailabilitySetId, model *availabilitysets.AvailabilitySet) error {
	d.Set("name", id.AvailabilitySetName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
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

	return pluginsdk.SetResourceIdentityData(d, id)
}
