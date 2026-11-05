// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dedicated_host_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/dedicatedhostgroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceDedicatedHostGroupFlatten(d *pluginsdk.ResourceData, id *commonids.DedicatedHostGroupId, model *dedicatedhostgroups.DedicatedHostGroup) error {
	d.Set("name", id.HostGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		zone := ""
		if model.Zones != nil && len(*model.Zones) > 0 {
			z := *model.Zones
			zone = z[0]
		}
		d.Set("zone", zone)

		if props := model.Properties; props != nil {
			d.Set("platform_fault_domain_count", props.PlatformFaultDomainCount)
			d.Set("automatic_placement_enabled", props.SupportAutomaticPlacement)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
