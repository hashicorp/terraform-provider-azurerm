// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkProfileUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkprofiles.ParseNetworkProfileID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, networkprofiles.DefaultGetOperationOptions())
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

	containerNetworkInterfaceConfigurations := expandNetworkProfileContainerNetworkInterface(d.Get("container_network_interface").([]any))
	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(containerNetworkInterfaceConfigurations)
	if err != nil {
		return fmt.Errorf("extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	locks.MultipleByID(vnetsToLock)
	defer locks.UnlockMultipleByID(vnetsToLock)

	locks.MultipleByID(subnetsToLock)
	defer locks.UnlockMultipleByID(subnetsToLock)

	if d.HasChange("container_network_interface") {
		payload.Properties.ContainerNetworkInterfaceConfigurations = containerNetworkInterfaceConfigurations
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkProfileRead(d, meta)
}
