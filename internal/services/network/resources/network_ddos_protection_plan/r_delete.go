// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_ddos_protection_plan

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkDDoSProtectionPlanDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).DdosProtectionPlans
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}
	// if there's no VirtualNetworks configured, it's possible for this to be nil
	subResources := existing.Model.Properties.VirtualNetworks
	virtualNetworkIDsToLock, err := extractVnetIDs(subResources)
	if err != nil {
		return fmt.Errorf("extracting IDs of Virtual Network: %+v", err)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	locks.MultipleByID(virtualNetworkIDsToLock)
	defer locks.UnlockMultipleByID(virtualNetworkIDsToLock)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
