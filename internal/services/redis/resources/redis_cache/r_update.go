// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisCacheUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Redis.RedisResourcesClient
	patchClient := meta.(*clients.Client).Redis.PatchSchedulesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redisresources.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	enableNonSslPort := d.Get("non_ssl_port_enabled")

	t := d.Get("tags").(map[string]any)
	expandedTags := tags.Expand(t)

	parameters := redisresources.RedisUpdateParameters{
		Properties: &redisresources.RedisUpdateProperties{
			DisableAccessKeyAuthentication: pointer.To(!d.Get("access_keys_authentication_enabled").(bool)),
			MinimumTlsVersion:              pointer.ToEnum[redisresources.TlsVersion](d.Get("minimum_tls_version").(string)),
			EnableNonSslPort:               pointer.To(enableNonSslPort.(bool)),
			Sku: &redisresources.Sku{
				Capacity: int64(d.Get("capacity").(int)),
				Family:   redisresources.SkuFamily(d.Get("family").(string)),
				Name:     redisresources.SkuName(d.Get("sku_name").(string)),
			},
		},
		Tags: expandedTags,
	}

	if d.HasChange("shard_count") {
		parameters.Properties.ShardCount = pointer.To(int64(d.Get("shard_count").(int)))
	}

	if d.HasChange("replicas_per_master") {
		parameters.Properties.ReplicasPerMaster = pointer.To(int64(d.Get("replicas_per_master").(int)))
	}

	if d.HasChange("replicas_per_primary") {
		parameters.Properties.ReplicasPerPrimary = pointer.To(int64(d.Get("replicas_per_primary").(int)))
	}

	if d.HasChange("redis_version") {
		parameters.Properties.RedisVersion = pointer.To(d.Get("redis_version").(string))
	}

	if d.HasChange("tenant_settings") {
		parameters.Properties.TenantSettings = expandTenantSettings(d.Get("tenant_settings").(map[string]any))
	}

	if d.HasChange("public_network_access_enabled") {
		parameters.Properties.PublicNetworkAccess = pointer.To(redisresources.PublicNetworkAccessDisabled)
		if d.Get("public_network_access_enabled").(bool) {
			parameters.Properties.PublicNetworkAccess = pointer.To(redisresources.PublicNetworkAccessEnabled)
		}
	}

	if d.HasChange("redis_configuration") {
		redisConfiguration, err := expandRedisConfiguration(d)
		if err != nil {
			return fmt.Errorf("parsing Redis Configuration: %+v", err)
		}
		parameters.Properties.RedisConfiguration = redisConfiguration
	}

	if _, err := client.RedisUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to become available", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating", "UpgradingRedisServerVersion", "ConfiguringAAD", "UpdatingManagedIdentity"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	// identity cannot be updated with sku,publicNetworkAccess,redisVersion etc.
	if d.HasChange("identity") {
		redisIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf(`expanding "identity": %v`, err)
		}

		identityParameter := redisresources.RedisUpdateParameters{
			Identity: redisIdentity,
		}
		if _, err := client.RedisUpdate(ctx, *id, identityParameter); err != nil {
			return fmt.Errorf("updating identity for %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Waiting for %s to become available", id)
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to become available: %+v", id, err)
		}
	}

	if d.HasChange("patch_schedule") {
		patchSchedule := expandRedisPatchSchedule(d)

		patchSchedulesRedisId := redispatchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
		if patchSchedule == nil || len(patchSchedule.Properties.ScheduleEntries) == 0 {
			if _, err = patchClient.PatchSchedulesDelete(ctx, patchSchedulesRedisId); err != nil {
				return fmt.Errorf("deleting Patch Schedule for %s: %+v", *id, err)
			}
		} else {
			if _, err = patchClient.PatchSchedulesCreateOrUpdate(ctx, patchSchedulesRedisId, *patchSchedule); err != nil {
				return fmt.Errorf("setting Patch Schedule for %s: %+v", *id, err)
			}
		}
	}

	return resourceRedisCacheRead(d, meta)
}
