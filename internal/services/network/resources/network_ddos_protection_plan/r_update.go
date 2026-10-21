// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_ddos_protection_plan

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkDDoSProtectionPlanUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetsToLock, err := expandNetworkDDoSProtectionPlanVnetIDs(d.Get("virtual_network_ids").([]any))
	if err != nil {
		return fmt.Errorf("extracting IDs of Virtual Network: %+v", err)
	}

	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(d.Id())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())
	locks.MultipleByID(vnetsToLock)
	defer locks.UnlockMultipleByID(vnetsToLock)

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkDDoSProtectionPlanRead(d, meta)
}
