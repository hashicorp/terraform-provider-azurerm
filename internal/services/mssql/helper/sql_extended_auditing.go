// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func FlattenSqlServerBlobAuditingPolicies(extendedServerBlobAuditingPolicy *sql.ExtendedServerBlobAuditingPolicy, d *pluginsdk.ResourceData) []interface{} {
	if extendedServerBlobAuditingPolicy == nil || extendedServerBlobAuditingPolicy.State == sql.BlobAuditingPolicyStateDisabled {
		return []interface{}{}
	}
	var storageEndpoint, storageAccessKey string
	// storage_account_access_key will not be returned, so we transfer the schema value
	if v, ok := d.GetOk("extended_auditing_policy.0.storage_account_access_key"); ok {
		storageAccessKey = v.(string)
	}
	if extendedServerBlobAuditingPolicy.StorageEndpoint != nil {
		storageEndpoint = *extendedServerBlobAuditingPolicy.StorageEndpoint
	}

	var secondKeyInUse bool
	if extendedServerBlobAuditingPolicy.IsStorageSecondaryKeyInUse != nil {
		secondKeyInUse = *extendedServerBlobAuditingPolicy.IsStorageSecondaryKeyInUse
	}
	var retentionDays int32
	if extendedServerBlobAuditingPolicy.RetentionDays != nil {
		retentionDays = *extendedServerBlobAuditingPolicy.RetentionDays
	}
	var monitor bool
	if extendedServerBlobAuditingPolicy.IsAzureMonitorTargetEnabled != nil {
		monitor = *extendedServerBlobAuditingPolicy.IsAzureMonitorTargetEnabled
	}

	var storageAccountSubscriptionId string
	if extendedServerBlobAuditingPolicy.ExtendedServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID != nil && extendedServerBlobAuditingPolicy.StorageAccountSubscriptionID.String() != "00000000-0000-0000-0000-000000000000" {
		storageAccountSubscriptionId = extendedServerBlobAuditingPolicy.StorageAccountSubscriptionID.String()
	}

	return []interface{}{
		map[string]interface{}{
			"storage_account_access_key":              storageAccessKey,
			"storage_endpoint":                        storageEndpoint,
			"storage_account_access_key_is_secondary": secondKeyInUse,
			"retention_in_days":                       retentionDays,
			"log_monitoring_enabled":                  monitor,
			"storage_account_subscription_id":         storageAccountSubscriptionId,
		},
	}
}

func ExpandMsSqlDBBlobAuditingPolicies(input []interface{}) *sql.ExtendedDatabaseBlobAuditingPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return &sql.ExtendedDatabaseBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		}
	}
	dbBlobAuditingPolicies := input[0].(map[string]interface{})

	ExtendedDatabaseBlobAuditingPolicyProperties := sql.ExtendedDatabaseBlobAuditingPolicyProperties{
		State:                       sql.BlobAuditingPolicyStateEnabled,
		StorageAccountAccessKey:     utils.String(dbBlobAuditingPolicies["storage_account_access_key"].(string)),
		StorageEndpoint:             utils.String(dbBlobAuditingPolicies["storage_endpoint"].(string)),
		IsAzureMonitorTargetEnabled: utils.Bool(dbBlobAuditingPolicies["log_monitoring_enabled"].(bool)),
	}
	if v, ok := dbBlobAuditingPolicies["storage_account_access_key_is_secondary"]; ok {
		ExtendedDatabaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = utils.Bool(v.(bool))
	}
	if v, ok := dbBlobAuditingPolicies["retention_in_days"]; ok {
		ExtendedDatabaseBlobAuditingPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ExtendedDatabaseBlobAuditingPolicyProperties
}

func FlattenMsSqlDBBlobAuditingPolicies(extendedDatabaseBlobAuditingPolicy *sql.ExtendedDatabaseBlobAuditingPolicy, d *pluginsdk.ResourceData) []interface{} {
	if extendedDatabaseBlobAuditingPolicy == nil || extendedDatabaseBlobAuditingPolicy.State == sql.BlobAuditingPolicyStateDisabled {
		return []interface{}{}
	}

	var storageAccessKey, storageEndpoint string
	// storage_account_access_key will not be returned, so we transfer the schema value
	if v, ok := d.GetOk("extended_auditing_policy.0.storage_account_access_key"); ok {
		storageAccessKey = v.(string)
	}

	if extendedDatabaseBlobAuditingPolicy.StorageEndpoint != nil {
		storageEndpoint = *extendedDatabaseBlobAuditingPolicy.StorageEndpoint
	}
	var secondKeyInUse bool
	if extendedDatabaseBlobAuditingPolicy.IsStorageSecondaryKeyInUse != nil {
		secondKeyInUse = *extendedDatabaseBlobAuditingPolicy.IsStorageSecondaryKeyInUse
	}
	var retentionDays int32
	if extendedDatabaseBlobAuditingPolicy.RetentionDays != nil {
		retentionDays = *extendedDatabaseBlobAuditingPolicy.RetentionDays
	}
	var monitor bool
	if extendedDatabaseBlobAuditingPolicy.IsAzureMonitorTargetEnabled != nil {
		monitor = *extendedDatabaseBlobAuditingPolicy.IsAzureMonitorTargetEnabled
	}

	return []interface{}{
		map[string]interface{}{
			"storage_account_access_key":              storageAccessKey,
			"storage_endpoint":                        storageEndpoint,
			"storage_account_access_key_is_secondary": secondKeyInUse,
			"retention_in_days":                       retentionDays,
			"log_monitoring_enabled":                  monitor,
		},
	}
}
