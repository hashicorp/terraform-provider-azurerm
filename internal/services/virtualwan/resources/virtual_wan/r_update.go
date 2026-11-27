// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_wan

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualWanUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualwans.ParseVirtualWANID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VirtualWansGet(ctx, *id)
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

	if d.HasChange("disable_vpn_encryption") {
		payload.Properties.DisableVpnEncryption = pointer.To(d.Get("disable_vpn_encryption").(bool))
	}

	if d.HasChange("allow_branch_to_branch_traffic") {
		payload.Properties.AllowBranchToBranchTraffic = pointer.To(d.Get("allow_branch_to_branch_traffic").(bool))
	}

	if d.HasChange("office365_local_breakout_category") {
		payload.Properties.Office365LocalBreakoutCategory = pointer.ToEnum[virtualwans.OfficeTrafficCategory](d.Get("office365_local_breakout_category").(string))
	}

	if d.HasChange("type") {
		payload.Properties.Type = pointer.To(d.Get("type").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.VirtualWansCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualWanRead(d, meta)
}
