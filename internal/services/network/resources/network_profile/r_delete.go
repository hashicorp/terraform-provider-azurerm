// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkProfileDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).NetworkProfiles
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := networkprofiles.ParseNetworkProfileID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, networkprofiles.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model.Properties` was nil", *id)
	}

	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(existing.Model.Properties.ContainerNetworkInterfaceConfigurations)
	if err != nil {
		return fmt.Errorf("extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	locks.MultipleByID(vnetsToLock)
	defer locks.UnlockMultipleByID(vnetsToLock)

	locks.MultipleByID(subnetsToLock)
	defer locks.UnlockMultipleByID(subnetsToLock)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
