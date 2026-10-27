// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceRedisCacheFlatten(ctx context.Context, metaClient *clients.Client, d *pluginsdk.ResourceData, id *redisresources.RediId, redis *redisresources.RedisResource, fetchCompleteData bool) error {
	client := metaClient.Redis.RedisResourcesClient
	patchSchedulesClient := metaClient.Redis.PatchSchedulesClient

	d.Set("name", id.RedisName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if redis != nil {
		redisIdentity, err := identity.FlattenSystemAndUserAssignedMap(redis.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", redisIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("location", location.Normalize(redis.Location))
		d.Set("zones", zones.FlattenUntyped(redis.Zones))

		props := redis.Properties
		d.Set("capacity", int(props.Sku.Capacity))
		d.Set("family", string(props.Sku.Family))
		d.Set("sku_name", string(props.Sku.Name))
		d.Set("ssl_port", props.SslPort)
		d.Set("hostname", props.HostName)
		d.Set("port", props.Port)
		d.Set("non_ssl_port_enabled", props.EnableNonSslPort)

		minimumTlsVersion := string(redisresources.TlsVersionOnePointTwo)
		if props.MinimumTlsVersion != nil {
			minimumTlsVersion = string(*props.MinimumTlsVersion)
		}
		d.Set("minimum_tls_version", minimumTlsVersion)

		shardCount := 0
		if props.ShardCount != nil {
			shardCount = int(*props.ShardCount)
		}
		d.Set("shard_count", shardCount)
		d.Set("private_static_ip_address", props.StaticIP)

		subnetId := ""
		if props.SubnetId != nil {
			parsed, err := commonids.ParseSubnetIDInsensitively(*props.SubnetId)
			if err != nil {
				return err
			}

			subnetId = parsed.ID()
		}
		d.Set("subnet_id", subnetId)

		publicNetworkAccessEnabled := true
		if props.PublicNetworkAccess != nil {
			publicNetworkAccessEnabled = *props.PublicNetworkAccess == redisresources.PublicNetworkAccessEnabled
		}
		d.Set("public_network_access_enabled", publicNetworkAccessEnabled)
		d.Set("replicas_per_master", props.ReplicasPerMaster)
		d.Set("replicas_per_primary", props.ReplicasPerPrimary)
		d.Set("redis_version", props.RedisVersion)
		d.Set("tenant_settings", flattenTenantSettings(props.TenantSettings))

		redisConfiguration, err := flattenRedisConfiguration(d, props.RedisConfiguration)
		if err != nil {
			return fmt.Errorf("flattening `redis_configuration`: %+v", err)
		}
		if err := d.Set("redis_configuration", redisConfiguration); err != nil {
			return fmt.Errorf("setting `redis_configuration`: %+v", err)
		}

		if fetchCompleteData {
			keysResp, err := client.RedisListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			d.Set("primary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.Model.PrimaryKey, true))
			d.Set("secondary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.Model.SecondaryKey, true))
			d.Set("primary_access_key", keysResp.Model.PrimaryKey)
			d.Set("secondary_access_key", keysResp.Model.SecondaryKey)
			d.Set("access_keys_authentication_enabled", !pointer.From(props.DisableAccessKeyAuthentication))

			patchSchedulesRedisId := redispatchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
			schedule, err := patchSchedulesClient.PatchSchedulesGet(ctx, patchSchedulesRedisId)
			var patchSchedule []any
			if err == nil {
				patchSchedule = flattenRedisPatchSchedules(*schedule.Model)
			}
			if err = d.Set("patch_schedule", patchSchedule); err != nil {
				return fmt.Errorf("setting `patch_schedule`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, redis.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenTenantSettings(input *map[string]string) map[string]string {
	output := make(map[string]string)

	if input != nil {
		maps.Copy(output, *input)
	}

	return output
}

func flattenRedisConfiguration(d *pluginsdk.ResourceData, input *redisresources.RedisCommonPropertiesRedisConfiguration) ([]any, error) {
	outputs := make(map[string]any)

	if input.AadEnabled != nil {
		a, err := strconv.ParseBool(*input.AadEnabled)
		if err != nil {
			return nil, fmt.Errorf("parsing `aad-enabled` %q: %+v", *input.AadEnabled, err)
		}
		outputs["active_directory_authentication_enabled"] = a
	}

	if input.Maxclients != nil {
		i, err := strconv.Atoi(*input.Maxclients)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxclients` %q: %+v", *input.Maxclients, err)
		}
		outputs["maxclients"] = i
	}
	if input.MaxmemoryDelta != nil {
		i, err := strconv.Atoi(*input.MaxmemoryDelta)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxmemory-delta` %q: %+v", *input.MaxmemoryDelta, err)
		}
		outputs["maxmemory_delta"] = i
	}
	if input.MaxmemoryReserved != nil {
		i, err := strconv.Atoi(*input.MaxmemoryReserved)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxmemory-reserved` %q: %+v", *input.MaxmemoryReserved, err)
		}
		outputs["maxmemory_reserved"] = i
	}
	if input.MaxmemoryPolicy != nil {
		outputs["maxmemory_policy"] = *input.MaxmemoryPolicy
	}

	if input.PreferredDataPersistenceAuthMethod != nil {
		outputs["data_persistence_authentication_method"] = *input.PreferredDataPersistenceAuthMethod
	}

	if input.MaxfragmentationmemoryReserved != nil {
		i, err := strconv.Atoi(*input.MaxfragmentationmemoryReserved)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxfragmentationmemory-reserved` %q: %+v", *input.MaxfragmentationmemoryReserved, err)
		}
		outputs["maxfragmentationmemory_reserved"] = i
	}

	// delta, reserved, enabled, frequency,, count,
	if input.RdbBackupEnabled != nil {
		b, err := strconv.ParseBool(*input.RdbBackupEnabled)
		if err != nil {
			return nil, fmt.Errorf("parsing `rdb-backup-enabled` %q: %+v", *input.RdbBackupEnabled, err)
		}
		outputs["rdb_backup_enabled"] = b
	}
	if input.RdbBackupFrequency != nil {
		i, err := strconv.Atoi(*input.RdbBackupFrequency)
		if err != nil {
			return nil, fmt.Errorf("parsing `rdb-backup-frequency` %q: %+v", *input.RdbBackupFrequency, err)
		}
		outputs["rdb_backup_frequency"] = i
	}
	if input.RdbBackupMaxSnapshotCount != nil {
		i, err := strconv.Atoi(*input.RdbBackupMaxSnapshotCount)
		if err != nil {
			return nil, fmt.Errorf("parsing `rdb-backup-max-snapshot-count` %q: %+v", *input.RdbBackupMaxSnapshotCount, err)
		}
		outputs["rdb_backup_max_snapshot_count"] = i
	}
	if input.RdbStorageConnectionString != nil {
		// The API returns AccountKey=[key hidden] instead of the value being passed in so we'll just set that value to what Terraform thinks the value is
		if len(strings.Split(*input.RdbStorageConnectionString, "AccountKey=[key hidden]")) > 1 {
			outputs["rdb_storage_connection_string"] = d.Get("redis_configuration.0.rdb_storage_connection_string")
		} else {
			outputs["rdb_storage_connection_string"] = *input.RdbStorageConnectionString
		}
	}
	outputs["notify_keyspace_events"] = pointer.From(input.NotifyKeyspaceEvents)

	if v := input.AofBackupEnabled; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("parsing `aof-backup-enabled` %q: %+v", *v, err)
		}
		outputs["aof_backup_enabled"] = b
	}
	if input.AofStorageConnectionString0 != nil {
		// The API returns AccountKey=[key hidden] instead of the value being passed in so we'll just set that value to what Terraform thinks the value is
		if len(strings.Split(*input.AofStorageConnectionString0, "AccountKey=[key hidden]")) > 1 {
			outputs["aof_storage_connection_string_0"] = d.Get("redis_configuration.0.aof_storage_connection_string_0")
		} else {
			outputs["aof_storage_connection_string_0"] = *input.AofStorageConnectionString0
		}
	}
	if input.AofStorageConnectionString1 != nil {
		// The API returns AccountKey=[key hidden] instead of the value being passed in so we'll just set that value to what Terraform thinks the value is
		if len(strings.Split(*input.AofStorageConnectionString1, "AccountKey=[key hidden]")) > 1 {
			outputs["aof_storage_connection_string_1"] = d.Get("redis_configuration.0.aof_storage_connection_string_1")
		} else {
			outputs["aof_storage_connection_string_1"] = *input.AofStorageConnectionString1
		}
	}

	// `authnotrequired` is not set for instances launched outside a VNET
	outputs["authentication_enabled"] = true
	if v := input.Authnotrequired; v != nil {
		outputs["authentication_enabled"] = isAuthRequiredAsBool(*v)
	}

	outputs["storage_account_subscription_id"] = pointer.From(input.StorageSubscriptionId)

	return []any{outputs}, nil
}

func flattenRedisPatchSchedules(schedule redispatchschedules.RedisPatchSchedule) []any {
	outputs := make([]any, 0)

	for _, entry := range schedule.Properties.ScheduleEntries {
		maintenanceWindow := pointer.From(entry.MaintenanceWindow)

		outputs = append(outputs, map[string]any{
			"day_of_week":        string(entry.DayOfWeek),
			"maintenance_window": maintenanceWindow,
			"start_hour_utc":     int(entry.StartHourUtc),
		})
	}

	return outputs
}
