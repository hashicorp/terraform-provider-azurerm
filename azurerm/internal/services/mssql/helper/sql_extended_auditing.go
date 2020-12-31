package helper

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ExtendedAuditingSchema() *schema.Schema {
	return &schema.Schema{
		Type:       schema.TypeList,
		Optional:   true,
		Computed:   true,
		Deprecated: "the `extended_auditing_policy` block has been moved to `azurerm_mssql_server_extended_auditing_policy` and `azurerm_mssql_database_extended_auditing_policy`. This block will be removed in version 3.0 of the provider.",
		ConfigMode: schema.SchemaConfigModeAttr,
		MaxItems:   1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"storage_account_access_key": {
					Type:         schema.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"storage_endpoint": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},

				"storage_account_access_key_is_secondary": {
					Type:     schema.TypeBool,
					Optional: true,
				},

				"retention_in_days": {
					Type:         schema.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(0, 3285),
				},
			},
		},
	}
}

func ExpandSqlServerBlobAuditingPolicies(input []interface{}) *sql.ExtendedServerBlobAuditingPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return &sql.ExtendedServerBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,

			// NOTE: this works around a regression in the Azure API detailed here:
			// https://github.com/Azure/azure-rest-api-specs/issues/11271
			IsAzureMonitorTargetEnabled: utils.Bool(true),
		}
	}
	serverBlobAuditingPolicies := input[0].(map[string]interface{})

	ExtendedServerBlobAuditingPolicyProperties := sql.ExtendedServerBlobAuditingPolicyProperties{
		State:                   sql.BlobAuditingPolicyStateEnabled,
		StorageAccountAccessKey: utils.String(serverBlobAuditingPolicies["storage_account_access_key"].(string)),
		StorageEndpoint:         utils.String(serverBlobAuditingPolicies["storage_endpoint"].(string)),
	}
	if v, ok := serverBlobAuditingPolicies["storage_account_access_key_is_secondary"]; ok {
		ExtendedServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = utils.Bool(v.(bool))
	}
	if v, ok := serverBlobAuditingPolicies["retention_in_days"]; ok {
		ExtendedServerBlobAuditingPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ExtendedServerBlobAuditingPolicyProperties
}

func FlattenSqlServerBlobAuditingPolicies(extendedServerBlobAuditingPolicy *sql.ExtendedServerBlobAuditingPolicy, d *schema.ResourceData) []interface{} {
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

	return []interface{}{
		map[string]interface{}{
			"storage_account_access_key":              storageAccessKey,
			"storage_endpoint":                        storageEndpoint,
			"storage_account_access_key_is_secondary": secondKeyInUse,
			"retention_in_days":                       retentionDays,
		},
	}
}

func ExpandMsSqlDBBlobAuditingPolicies(input []interface{}) *sql.ExtendedDatabaseBlobAuditingPolicyProperties {
	if len(input) == 0 || input[0] == nil {
		return &sql.ExtendedDatabaseBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,

			// NOTE: this works around a regression in the Azure API detailed here:
			// https://github.com/Azure/azure-rest-api-specs/issues/11271
			IsAzureMonitorTargetEnabled: utils.Bool(true),
		}
	}
	dbBlobAuditingPolicies := input[0].(map[string]interface{})

	ExtendedDatabaseBlobAuditingPolicyProperties := sql.ExtendedDatabaseBlobAuditingPolicyProperties{
		State:                   sql.BlobAuditingPolicyStateEnabled,
		StorageAccountAccessKey: utils.String(dbBlobAuditingPolicies["storage_account_access_key"].(string)),
		StorageEndpoint:         utils.String(dbBlobAuditingPolicies["storage_endpoint"].(string)),
	}
	if v, ok := dbBlobAuditingPolicies["storage_account_access_key_is_secondary"]; ok {
		ExtendedDatabaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = utils.Bool(v.(bool))
	}
	if v, ok := dbBlobAuditingPolicies["retention_in_days"]; ok {
		ExtendedDatabaseBlobAuditingPolicyProperties.RetentionDays = utils.Int32(int32(v.(int)))
	}

	return &ExtendedDatabaseBlobAuditingPolicyProperties
}

func FlattenMsSqlDBBlobAuditingPolicies(extendedDatabaseBlobAuditingPolicy *sql.ExtendedDatabaseBlobAuditingPolicy, d *schema.ResourceData) []interface{} {
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

	return []interface{}{
		map[string]interface{}{
			"storage_account_access_key":              storageAccessKey,
			"storage_endpoint":                        storageEndpoint,
			"storage_account_access_key_is_secondary": secondKeyInUse,
			"retention_in_days":                       retentionDays,
		},
	}
}
