// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisCacheDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Redis.RedisResourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redisresources.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.RedisGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if subnetID := read.Model.Properties.SubnetId; subnetID != nil {
		parsed, parseErr := commonids.ParseSubnetIDInsensitively(*subnetID)
		if parseErr != nil {
			return err
		}

		locks.ByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(parsed.SubnetName, network.SubnetResourceName)
		defer locks.UnlockByName(parsed.SubnetName, network.SubnetResourceName)
	}

	if err := client.RedisDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
