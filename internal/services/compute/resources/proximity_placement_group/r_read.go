// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package proximity_placement_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceProximityPlacementGroupRead(d *pluginsdk.ResourceData, meta any) error {
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

	return pluginsdk.SetResourceIdentityData(d, id)
}
