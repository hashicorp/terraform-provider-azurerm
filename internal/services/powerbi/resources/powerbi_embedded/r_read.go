// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package powerbi_embedded

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePowerBIEmbeddedRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacities.ParseCapacityID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetDetails(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CapacityName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			var adminMembers *[]string
			if props.Administration != nil {
				adminMembers = props.Administration.Members
			}
			if err := d.Set("administrators", pluginsdk.FlattenSlice(adminMembers)); err != nil {
				return fmt.Errorf("setting `administration`: %+v", err)
			}

			mode := ""
			if props.Mode != nil {
				mode = string(*props.Mode)
			}
			d.Set("mode", mode)
		}

		d.Set("sku_name", model.Sku.Name)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
