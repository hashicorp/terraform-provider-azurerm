// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/resourcegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceResourceGroupFlatten(d *pluginsdk.ResourceData, id *commonids.ResourceGroupId, group *resourcegroups.ResourceGroup) error {
	d.Set("name", id.ResourceGroupName)

	if group != nil {
		d.Set("location", location.Normalize(group.Location))
		d.Set("managed_by", pointer.From(group.ManagedBy))

		if err := tags.FlattenAndSet(d, group.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
