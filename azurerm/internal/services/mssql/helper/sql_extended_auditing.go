package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ExtendedAuditingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:       pluginsdk.TypeList,
		Optional:   true,
		Computed:   true,
		Deprecated: "the `extended_auditing_policy` block has been moved to `azurerm_mssql_server_extended_auditing_policy` and `azurerm_mssql_database_extended_auditing_policy`. This block will be removed in version 3.0 of the provider.",
		ConfigMode: pluginsdk.SchemaConfigModeAttr,
		MaxItems:   1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"storage_account_access_key": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_endpoint": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},

				"storage_account_access_key_is_secondary": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},

				"retention_in_days": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 3285),
				},

				"log_monitoring_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},
			},
		},
	}
}

func ExpandSqlServerBlobAuditingPolicies(input []interface{}) *sql.ExtendedServerBlobAuditingPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return &sql.ExtendedServerBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		}
	}
	serverBlobAuditingPolicies := input[0].(map[string]interface{})

	ExtendedServerBlobAuditingPolicyProperties := sql.ExtendedServerBlobAuditingPolicyProperties{
		State:                       sql.BlobAuditingPolicyStateEnabled,
		StorageAccountAccessKey:     utils.String(serverBlobAuditingPolicies["storage_account_access_key"].(string)),
		StorageEndpoint:             utils.String(serverBlobAuditingPolicies["storage_endpoint"].(string)),
		IsAzureMonitorTargetEnabled: utils.Bool(serverBlobAuditingPolicies["log_monitoring_enabled"].(bool)),
	}
	if v, ok := serverBlobAuditingPolicies["storage_account_access_key_is_secondary"]; ok {
		ExtendedServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = utils.Bool(v.(bool))
	}
	if v, ok := serverBlobAuditingPolicies["retention_in_days"]; ok {
		ExtendedServerBlobAuditingPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ExtendedServerBlobAuditingPolicyProperties
}

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
