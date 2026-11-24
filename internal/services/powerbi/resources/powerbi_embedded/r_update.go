// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package powerbi_embedded

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePowerBIEmbeddedUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).PowerBI.CapacityClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := capacities.ParseCapacityID(d.Id())
	if err != nil {
		return err
	}

	parameters := capacities.DedicatedCapacityUpdateParameters{}

	if d.HasChanges("administrators", "mode") {
		administrators := d.Get("administrators").(*pluginsdk.Set).List()
		mode := capacities.Mode(d.Get("mode").(string))

		parameters.Properties = &capacities.DedicatedCapacityMutableProperties{
			Administration: &capacities.DedicatedCapacityAdministrators{
				Members: pluginsdk.ExpandStringSlice(administrators),
			},
			Mode: &mode,
		}
	}

	if d.HasChange("sku_name") {
		parameters.Sku = &capacities.CapacitySku{
			Name: d.Get("sku_name").(string),
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourcePowerBIEmbeddedRead(d, meta)
}
