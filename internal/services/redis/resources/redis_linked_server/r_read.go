// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_linked_server

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/linkedserver"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisLinkedServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := linkedserver.ParseLinkedServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LinkedServerName)
	d.Set("target_redis_cache_name", id.RedisName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			cacheId, err := linkedserver.ParseRediIDInsensitively(props.LinkedRedisCacheId)
			if err != nil {
				return fmt.Errorf("parsing `linkedRedisCacheId` %q: %+v", props.LinkedRedisCacheId, err)
			}
			d.Set("linked_redis_cache_id", cacheId.ID())

			d.Set("linked_redis_cache_location", location.Normalize(props.LinkedRedisCacheLocation))
			d.Set("server_role", string(props.ServerRole))
			d.Set("geo_replicated_primary_host_name", pointer.From(props.GeoReplicatedPrimaryHostName))
		}
	}

	return nil
}

func redisLinkedServerStateRefreshFunc(ctx context.Context, client *linkedserver.LinkedServerClient, id linkedserver.LinkedServerId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving status of %s: %+v", id, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil && props.ProvisioningState != nil {
				return resp, *props.ProvisioningState, nil
			}
		}

		return nil, "", fmt.Errorf("retrieving %s: `model` was nil", id)
	}
}
