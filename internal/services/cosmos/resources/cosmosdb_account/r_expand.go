// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandAzureRmCosmosDBAccountConsistencyPolicy(d *pluginsdk.ResourceData) *cosmosdb.ConsistencyPolicy {
	i := d.Get("consistency_policy").([]any)
	if len(i) == 0 || i[0] == nil {
		return nil
	}
	input := i[0].(map[string]any)

	consistencyLevel := input["consistency_level"].(string)
	policy := cosmosdb.ConsistencyPolicy{
		DefaultConsistencyLevel: cosmosdb.DefaultConsistencyLevel(consistencyLevel),
	}

	if stalenessPrefix, ok := input["max_staleness_prefix"].(int); ok {
		if stalenessPrefix == 0 {
			stalenessPrefix = 100
		}
		policy.MaxStalenessPrefix = pointer.To(int64(stalenessPrefix))
	}
	if maxInterval, ok := input["max_interval_in_seconds"].(int); ok {
		if maxInterval == 0 {
			maxInterval = 5
		}
		policy.MaxIntervalInSeconds = pointer.To(int64(maxInterval))
	}

	return &policy
}

func expandAzureRmCosmosDBAccountGeoLocations(d *pluginsdk.ResourceData) ([]cosmosdb.Location, error) {
	locations := make([]cosmosdb.Location, 0)
	for _, l := range d.Get("geo_location").(*pluginsdk.Set).List() {
		data := l.(map[string]any)

		location := cosmosdb.Location{
			LocationName:     pointer.To(location.Normalize(data["location"].(string))),
			FailoverPriority: pointer.To(int64(data["failover_priority"].(int))),
			IsZoneRedundant:  pointer.To(data["zone_redundant"].(bool)),
		}

		locations = append(locations, location)
	}

	// TODO: maybe this should be in a CustomizeDiff
	// all priorities & locations must be unique
	byPriorities := make(map[int]any, len(locations))
	byName := make(map[string]any, len(locations))
	locationsCount := len(locations)
	for _, location := range locations {
		priority := int(*location.FailoverPriority)
		name := *location.LocationName

		if _, ok := byPriorities[priority]; ok {
			return nil, fmt.Errorf("each `geo_location` needs to have a unique failover_priority. Multiple instances of '%d' found", priority)
		}

		if _, ok := byName[name]; ok {
			return nil, fmt.Errorf("each `geo_location` needs to be in unique location. Multiple instances of %q found", name)
		}

		if priority > locationsCount-1 {
			return nil, fmt.Errorf("the maximum value for a failover priority = (total number of regions - 1). '%d' was found", priority)
		}

		byPriorities[priority] = location
		byName[name] = location
	}

	// and one of them must have a priority of 0...
	if _, ok := byPriorities[0]; !ok {
		return nil, fmt.Errorf("there needs to be a `geo_location` with a `failover_priority` of `0`")
	}

	return locations, nil
}

func expandAzureRmCosmosDBAccountCapabilities(d *pluginsdk.ResourceData) *[]cosmosdb.Capability {
	capabilities := d.Get("capabilities").(*pluginsdk.Set).List()
	s := make([]cosmosdb.Capability, 0)

	for _, c := range capabilities {
		m := c.(map[string]any)
		s = append(s, cosmosdb.Capability{Name: pointer.To(m["name"].(string))})
	}

	return &s
}

func expandAzureRmCosmosDBAccountVirtualNetworkRules(d *pluginsdk.ResourceData) *[]cosmosdb.VirtualNetworkRule {
	virtualNetworkRules := d.Get("virtual_network_rule").(*pluginsdk.Set).List()

	s := make([]cosmosdb.VirtualNetworkRule, len(virtualNetworkRules))
	for i, r := range virtualNetworkRules {
		m := r.(map[string]any)
		s[i] = cosmosdb.VirtualNetworkRule{
			Id:                               pointer.To(m["id"].(string)),
			IgnoreMissingVNetServiceEndpoint: pointer.To(m["ignore_missing_vnet_service_endpoint"].(bool)),
		}
	}

	return &s
}

func expandCosmosdbAccountBackup(input []any, backupHasChange bool, createMode string) (cosmosdb.BackupPolicy, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}
	attr := input[0].(map[string]any)

	switch attr["type"].(string) {
	case string(cosmosdb.BackupPolicyTypeContinuous):
		if v := attr["interval_in_minutes"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`interval_in_minutes` cannot be defined when the `backup.type` is set to %q", cosmosdb.BackupPolicyTypeContinuous)
		}

		if v := attr["retention_in_hours"].(int); v != 0 && !backupHasChange {
			return nil, fmt.Errorf("`retention_in_hours` cannot be defined when the `backup.type` is set to %q", cosmosdb.BackupPolicyTypeContinuous)
		}

		if v := attr["storage_redundancy"].(string); v != "" && !backupHasChange {
			return nil, fmt.Errorf("`storage_redundancy` cannot be defined when the `backup.type` is set to %q", cosmosdb.BackupPolicyTypeContinuous)
		}

		result := cosmosdb.ContinuousModeBackupPolicy{}

		if v := attr["tier"].(string); v != "" {
			result.ContinuousModeProperties = &cosmosdb.ContinuousModeProperties{
				Tier: pointer.ToEnum[cosmosdb.ContinuousTier](v),
			}
		}

		return result, nil

	case string(cosmosdb.BackupPolicyTypePeriodic):
		if createMode != "" {
			return nil, fmt.Errorf("`create_mode` can only be defined when the `backup.type` is set to %q, got %q", cosmosdb.BackupPolicyTypeContinuous, cosmosdb.BackupPolicyTypePeriodic)
		}

		if v := attr["tier"].(string); v != "" && !backupHasChange {
			return nil, fmt.Errorf("`tier` can not be set when `type` in `backup` is `Periodic`")
		}

		periodicModeBackupPolicy := cosmosdb.PeriodicModeBackupPolicy{
			PeriodicModeProperties: &cosmosdb.PeriodicModeProperties{},
		}

		if v := attr["interval_in_minutes"].(int); v != 0 {
			periodicModeBackupPolicy.PeriodicModeProperties.BackupIntervalInMinutes = pointer.To(int64(v))
		}

		if v := attr["retention_in_hours"].(int); v != 0 {
			periodicModeBackupPolicy.PeriodicModeProperties.BackupRetentionIntervalInHours = pointer.To(int64(v))
		}

		if v := attr["storage_redundancy"].(string); v != "" {
			periodicModeBackupPolicy.PeriodicModeProperties.BackupStorageRedundancy = pointer.ToEnum[cosmosdb.BackupStorageRedundancy](attr["storage_redundancy"].(string))
		}

		return periodicModeBackupPolicy, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`:%+v", attr["type"].(string))
	}
}

func expandCosmosDBAccountAnalyticalStorageConfiguration(input []any) *cosmosdb.AnalyticalStorageConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &cosmosdb.AnalyticalStorageConfiguration{
		SchemaType: pointer.ToEnum[cosmosdb.AnalyticalStorageSchemaType](v["schema_type"].(string)),
	}
}

func expandCosmosDBAccountCapacity(input []any) *cosmosdb.Capacity {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &cosmosdb.Capacity{
		TotalThroughputLimit: pointer.To(int64(v["total_throughput_limit"].(int))),
	}
}

func expandCosmosdbAccountRestoreParameters(input []any) *cosmosdb.RestoreParameters {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	restoreParameters := cosmosdb.RestoreParameters{
		RestoreMode:               pointer.To(cosmosdb.RestoreModePointInTime),
		RestoreSource:             pointer.To(v["source_cosmosdb_account_id"].(string)),
		DatabasesToRestore:        expandCosmosdbAccountDatabasesToRestore(v["database"].(*pluginsdk.Set).List()),
		GremlinDatabasesToRestore: expandCosmosdbAccountGremlinDatabasesToRestore(v["gremlin_database"].([]any)),
	}

	restoreTimestampInUtc, _ := time.Parse(time.RFC3339, v["restore_timestamp_in_utc"].(string))
	restoreParameters.SetRestoreTimestampInUtcAsTime(restoreTimestampInUtc)

	if tablesToRestore := v["tables_to_restore"].([]any); len(tablesToRestore) > 0 {
		restoreParameters.TablesToRestore = pluginsdk.ExpandStringSlice(tablesToRestore)
	}

	return &restoreParameters
}

func expandCosmosdbAccountDatabasesToRestore(input []any) *[]cosmosdb.DatabaseRestoreResource {
	results := make([]cosmosdb.DatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, cosmosdb.DatabaseRestoreResource{
			DatabaseName:    pointer.To(v["name"].(string)),
			CollectionNames: pluginsdk.ExpandStringSlice(v["collection_names"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func expandCosmosdbAccountGremlinDatabasesToRestore(input []any) *[]cosmosdb.GremlinDatabaseRestoreResource {
	results := make([]cosmosdb.GremlinDatabaseRestoreResource, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, cosmosdb.GremlinDatabaseRestoreResource{
			DatabaseName: pointer.To(v["name"].(string)),
			GraphNames:   pluginsdk.ExpandStringSlice(v["graph_names"].([]any)),
		})
	}

	return &results
}

func expandCosmosdbAccountNetworkBypass(input bool) *cosmosdb.NetworkAclBypass {
	if input {
		return pointer.To(cosmosdb.NetworkAclBypassAzureServices)
	}
	return pointer.To(cosmosdb.NetworkAclBypassNone)
}

func expandCosmosdbAccountPublicNetworkAccess(input bool) *cosmosdb.PublicNetworkAccess {
	if input {
		return pointer.To(cosmosdb.PublicNetworkAccessEnabled)
	}
	return pointer.To(cosmosdb.PublicNetworkAccessDisabled)
}
