// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip_prefix

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpPrefixRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := publicipprefixes.ParsePublicIPPrefixID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, publicipprefixes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PublicIPPrefixName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		if sku := model.Sku; sku != nil {
			d.Set("sku", pointer.FromEnum(sku.Name))
			d.Set("sku_tier", pointer.FromEnum(sku.Tier))
		}
		if props := model.Properties; props != nil {
			d.Set("prefix_length", props.PrefixLength)
			d.Set("ip_prefix", props.IPPrefix)
			if version := props.PublicIPAddressVersion; version != nil {
				d.Set("ip_version", string(*version))
			}

			customIpPrefixId := ""
			if props.CustomIPPrefix != nil {
				id, err := customipprefixes.ParseCustomIPPrefixID(pointer.From(props.CustomIPPrefix.Id))
				if err != nil {
					return err
				}
				customIpPrefixId = id.ID()
			}
			d.Set("custom_ip_prefix_id", customIpPrefixId)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
