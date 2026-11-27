// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenAzureRmCosmosDBAccountConsistencyPolicy(policy *cosmosdb.ConsistencyPolicy) []any {
	result := map[string]any{}
	result["consistency_level"] = string(policy.DefaultConsistencyLevel)
	if policy.MaxIntervalInSeconds != nil {
		result["max_interval_in_seconds"] = int(*policy.MaxIntervalInSeconds)
	}
	if policy.MaxStalenessPrefix != nil {
		result["max_staleness_prefix"] = int(*policy.MaxStalenessPrefix)
	}

	return []any{result}
}

func flattenAzureRmCosmosDBAccountGeoLocations(account *cosmosdb.DatabaseAccountGetProperties) *pluginsdk.Set {
	locationSet := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountGeoLocationHash,
	}
	if account == nil || account.FailoverPolicies == nil {
		return &locationSet
	}

	for _, l := range *account.FailoverPolicies {
		if l.Id == nil {
			continue
		}

		id := *l.Id
		lb := map[string]any{
			"id":                id,
			"location":          location.NormalizeNilable(l.LocationName),
			"failover_priority": int(pointer.From(l.FailoverPriority)),
			// there is no zone redundancy information in FailoverPolicies currently, we have to search it by `id` in the Locations property.
			"zone_redundant": findZoneRedundant(account.Locations, id),
		}

		locationSet.Add(lb)
	}

	return &locationSet
}

func flattenAzureRmCosmosDBAccountCapabilities(capabilities *[]cosmosdb.Capability) *pluginsdk.Set {
	s := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountCapabilitiesHash,
	}

	for _, c := range *capabilities {
		if v := c.Name; v != nil {
			e := map[string]any{
				"name": *v,
			}
			s.Add(e)
		}
	}

	return &s
}

func flattenAzureRmCosmosDBAccountVirtualNetworkRules(rules *[]cosmosdb.VirtualNetworkRule) *pluginsdk.Set {
	results := pluginsdk.Set{
		F: resourceAzureRMCosmosDBAccountVirtualNetworkRuleHash,
	}

	if rules != nil {
		for _, r := range *rules {
			rule := map[string]any{
				"id":                                   *r.Id,
				"ignore_missing_vnet_service_endpoint": *r.IgnoreMissingVNetServiceEndpoint,
			}
			results.Add(rule)
		}
	}

	return &results
}

func flattenCosmosdbAccountBackup(input cosmosdb.BackupPolicy) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	switch backupPolicy := input.(type) {
	case cosmosdb.ContinuousModeBackupPolicy:
		var tier cosmosdb.ContinuousTier
		if v := backupPolicy.ContinuousModeProperties; v != nil {
			tier = pointer.From(v.Tier)
		}
		return []any{
			map[string]any{
				"type": string(cosmosdb.BackupPolicyTypeContinuous),
				"tier": string(tier),
			},
		}, nil

	case cosmosdb.PeriodicModeBackupPolicy:
		var interval, retention int
		if v := backupPolicy.PeriodicModeProperties.BackupIntervalInMinutes; v != nil {
			interval = int(*v)
		}

		if v := backupPolicy.PeriodicModeProperties.BackupRetentionIntervalInHours; v != nil {
			retention = int(*v)
		}

		var storageRedundancy cosmosdb.BackupStorageRedundancy
		if backupPolicy.PeriodicModeProperties.BackupStorageRedundancy != nil {
			storageRedundancy = pointer.From(backupPolicy.PeriodicModeProperties.BackupStorageRedundancy)
		}

		return []any{
			map[string]any{
				"type":                string(cosmosdb.BackupPolicyTypePeriodic),
				"interval_in_minutes": interval,
				"retention_in_hours":  retention,
				"storage_redundancy":  storageRedundancy,
			},
		}, nil

	default:
		return nil, fmt.Errorf("unknown `type` in `backup`: %+v", input)
	}
}

func flattenCosmosDBAccountAnalyticalStorageConfiguration(input *cosmosdb.AnalyticalStorageConfiguration) []any {
	if input == nil {
		return make([]any, 0)
	}

	var schemaType cosmosdb.AnalyticalStorageSchemaType
	if input.SchemaType != nil {
		schemaType = pointer.From(input.SchemaType)
	}

	return []any{
		map[string]any{
			"schema_type": schemaType,
		},
	}
}

func flattenCosmosDBAccountCapacity(input *cosmosdb.Capacity) []any {
	if input == nil {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"total_throughput_limit": pointer.From(input.TotalThroughputLimit),
		},
	}
}

func flattenCosmosdbAccountRestoreParameters(input *cosmosdb.RestoreParameters) []any {
	if input == nil {
		return make([]any, 0)
	}
	restoreSource := pointer.From(input.RestoreSource)

	var restoreTimestampInUtc string
	if input.RestoreTimestampInUtc != nil {
		restoreTimestampInUtc = pointer.From(input.RestoreTimestampInUtc)
	}

	return []any{
		map[string]any{
			"database":                   flattenCosmosdbAccountDatabasesToRestore(input.DatabasesToRestore),
			"gremlin_database":           flattenCosmosdbAccountGremlinDatabasesToRestore(input.GremlinDatabasesToRestore),
			"source_cosmosdb_account_id": restoreSource,
			"restore_timestamp_in_utc":   restoreTimestampInUtc,
			"tables_to_restore":          pointer.From(input.TablesToRestore),
		},
	}
}

func flattenCosmosdbAccountDatabasesToRestore(input *[]cosmosdb.DatabaseRestoreResource) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		databaseName := pointer.From(item.DatabaseName)

		results = append(results, map[string]any{
			"collection_names": pluginsdk.FlattenSlice(item.CollectionNames),
			"name":             databaseName,
		})
	}
	return results
}

func flattenCosmosdbAccountGremlinDatabasesToRestore(input *[]cosmosdb.GremlinDatabaseRestoreResource) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]any{
			"graph_names": pluginsdk.FlattenSlice(item.GraphNames),
			"name":        pointer.From(item.DatabaseName),
		})
	}

	return results
}

func flattenCosmosdbAccountReadWriteEndpoints(input *[]cosmosdb.Location) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, l := range *input {
		if l.DocumentEndpoint == nil {
			continue
		}
		result = append(result, *l.DocumentEndpoint)
	}

	return result
}
