// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nat_gateway

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceNatGatewayFlatten(d *pluginsdk.ResourceData, id *natgateways.NatGatewayId, model *natgateways.NatGateway) error {
	d.Set("name", id.NatGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		sku := ""
		if model.Sku != nil {
			sku = pointer.FromEnum(model.Sku.Name)
		}
		d.Set("sku_name", sku)
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		if props := model.Properties; props != nil {
			d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
			d.Set("resource_guid", props.ResourceGuid)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
