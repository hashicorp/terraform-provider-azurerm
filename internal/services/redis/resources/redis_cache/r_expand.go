// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandRedisConfiguration(d *pluginsdk.ResourceData) (*redisresources.RedisCommonPropertiesRedisConfiguration, error) {
	output := &redisresources.RedisCommonPropertiesRedisConfiguration{}

	input := d.Get("redis_configuration").([]any)
	if len(input) == 0 || input[0] == nil {
		return output, nil
	}
	raw := input[0].(map[string]any)
	skuName := d.Get("sku_name").(string)

	if v := raw["maxclients"].(int); v > 0 {
		output.Maxclients = pointer.To(strconv.Itoa(v))
	}

	if d.Get("sku_name").(string) != string(redisresources.SkuNameBasic) {
		if v := raw["maxmemory_delta"].(int); v > 0 {
			output.MaxmemoryDelta = pointer.To(strconv.Itoa(v))
		}

		if v := raw["maxmemory_reserved"].(int); v > 0 {
			output.MaxmemoryReserved = pointer.To(strconv.Itoa(v))
		}

		if v := raw["maxfragmentationmemory_reserved"].(int); v > 0 {
			output.MaxfragmentationmemoryReserved = pointer.To(strconv.Itoa(v))
		}
	}

	if v := raw["maxmemory_policy"].(string); v != "" {
		output.MaxmemoryPolicy = pointer.To(v)
	}

	output.PreferredDataPersistenceAuthMethod = pointer.To(raw["data_persistence_authentication_method"].(string))

	// AAD/Entra support
	v, valExists := d.GetOk("redis_configuration.0.active_directory_authentication_enabled")
	if valExists {
		entraEnabled := v.(bool)
		output.AadEnabled = pointer.To(strconv.FormatBool(entraEnabled))
	} else {
		output.AadEnabled = pointer.To("false")
	}

	// RDB Backup
	//nolint:staticcheck
	v, valExists = d.GetOkExists("redis_configuration.0.rdb_backup_enabled")
	if valExists {
		rdbBackupEnabled := v.(bool)

		// rdb_backup_enabled is available when SKU is Premium
		if strings.EqualFold(skuName, string(redisresources.SkuNamePremium)) {
			if rdbBackupEnabled {
				if connStr := raw["rdb_storage_connection_string"].(string); connStr == "" {
					return nil, fmt.Errorf("the rdb_storage_connection_string property must be set when rdb_backup_enabled is true")
				}
			}
			output.RdbBackupEnabled = pointer.To(strconv.FormatBool(rdbBackupEnabled))
		} else if rdbBackupEnabled && !strings.EqualFold(skuName, string(redisresources.SkuNamePremium)) {
			return nil, fmt.Errorf("the `rdb_backup_enabled` property requires a `Premium` sku to be set")
		}
	}

	if v := raw["rdb_backup_frequency"].(int); v > 0 {
		output.RdbBackupFrequency = pointer.To(strconv.Itoa(v))
	}

	if v := raw["rdb_backup_max_snapshot_count"].(int); v > 0 {
		output.RdbBackupMaxSnapshotCount = pointer.To(strconv.Itoa(v))
	}

	if v := raw["rdb_storage_connection_string"].(string); v != "" {
		output.RdbStorageConnectionString = pointer.To(v)
	}

	if v := raw["notify_keyspace_events"].(string); v != "" {
		output.NotifyKeyspaceEvents = pointer.To(v)
	}

	// AOF Backup
	//nolint:staticcheck
	v, valExists = d.GetOkExists("redis_configuration.0.aof_backup_enabled")
	if valExists {
		// aof_backup_enabled is available when SKU is Premium
		if strings.EqualFold(skuName, string(redisresources.SkuNamePremium)) {
			output.AofBackupEnabled = pointer.To(strconv.FormatBool(v.(bool)))
		}
	}

	if v := raw["aof_storage_connection_string_0"].(string); v != "" {
		output.AofStorageConnectionString0 = pointer.To(v)
	}

	if v := raw["aof_storage_connection_string_1"].(string); v != "" {
		output.AofStorageConnectionString1 = pointer.To(v)
	}

	authEnabled := raw["authentication_enabled"].(bool)

	// Redis authentication can only be disabled if it is launched inside a VNET.
	if _, isPrivate := d.GetOk("subnet_id"); !isPrivate {
		if !authEnabled {
			return nil, fmt.Errorf("cannot set `authentication_enabled` or `enable_authentication` to `false` when `subnet_id` is not set")
		}
	} else {
		value := isAuthNotRequiredAsString(authEnabled)
		output.Authnotrequired = pointer.To(value)
	}

	if v := raw["storage_account_subscription_id"].(string); v != "" {
		output.StorageSubscriptionId = pointer.To(v)
	}
	return output, nil
}

func expandRedisPatchSchedule(d *pluginsdk.ResourceData) *redispatchschedules.RedisPatchSchedule {
	v, ok := d.GetOk("patch_schedule")
	if !ok {
		return nil
	}

	scheduleValues := v.([]any)
	entries := make([]redispatchschedules.ScheduleEntry, 0)
	for _, scheduleValue := range scheduleValues {
		vals := scheduleValue.(map[string]any)
		dayOfWeek := vals["day_of_week"].(string)
		maintenanceWindow := vals["maintenance_window"].(string)
		startHourUtc := vals["start_hour_utc"].(int)

		entries = append(entries, redispatchschedules.ScheduleEntry{
			DayOfWeek:         redispatchschedules.DayOfWeek(dayOfWeek),
			MaintenanceWindow: pointer.To(maintenanceWindow),
			StartHourUtc:      int64(startHourUtc),
		})
	}

	schedule := redispatchschedules.RedisPatchSchedule{
		Properties: redispatchschedules.ScheduleEntries{
			ScheduleEntries: entries,
		},
	}
	return &schedule
}

func expandTenantSettings(input map[string]any) *map[string]string {
	output := make(map[string]string, len(input))

	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}
