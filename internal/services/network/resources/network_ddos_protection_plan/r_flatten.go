// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_ddos_protection_plan

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceNetworkDDoSProtectionPlanFlatten(d *pluginsdk.ResourceData, id *ddosprotectionplans.DdosProtectionPlanId, model *ddosprotectionplans.DdosProtectionPlan) error {
	d.Set("name", id.DdosProtectionPlanName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("virtual_network_ids", flattenNetworkDDoSProtectionPlanVirtualNetworkIDs(props.VirtualNetworks)); err != nil {
				return fmt.Errorf("setting `virtual_network_ids`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenNetworkDDoSProtectionPlanVirtualNetworkIDs(input *[]ddosprotectionplans.SubResource) []string {
	vnetIDs := make([]string, 0)
	if input == nil {
		return vnetIDs
	}

	// if-continue is used to simplify the deeply nested if-else statement.
	for _, subRes := range *input {
		if subRes.Id != nil {
			vnetIDs = append(vnetIDs, *subRes.Id)
		}
	}

	return vnetIDs
}
