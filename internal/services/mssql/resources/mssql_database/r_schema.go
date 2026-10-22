// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceMsSqlDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlDatabaseName,
		},

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlServerIDInsensitively),
		},

		"auto_pause_delay_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.DatabaseAutoPauseDelay,
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(databases.CreateModeDefault),
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForCreateMode(),
				false),
			ConflictsWith: []string{"import"},
		},

		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"storage_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"storage_key_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(databases.StorageKeyTypeSharedAccessKey),
							string(databases.StorageKeyTypeStorageAccessKey),
						}, false),
					},
					"administrator_login": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"administrator_login_password": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ADPassword",
							"Sql",
						}, false),
					},
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
				},
			},
			ConflictsWith: []string{"create_mode"}, // it needs `create_mode` to be `Default` to work, so make them conflict.
		},

		"collation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.DatabaseCollation(),
		},

		"elastic_pool_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlElasticPoolIDInsensitively),
		},

		"enclave_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForAlwaysEncryptedEnclaveType(), false),
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForDatabaseLicenseType(),
				false),
		},

		"long_term_retention_policy": helper.LongTermRetentionPolicySchema(),

		"short_term_retention_policy": helper.ShortTermRetentionPolicySchema(),

		"max_size_gb": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatBetween(0.1, 4096),
		},

		"min_capacity": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatInSlice([]float64{0, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40}),
		},

		"restore_point_in_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"recover_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RecoverableDatabaseID,
		},

		"recovery_point_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"restore_dropped_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RestorableDatabaseID,
		},

		"restore_long_term_retention_backup_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"read_replica_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"read_scale": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"sample_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(databases.SampleNameAdventureWorksLT),
			}, false),
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.DatabaseSkuName(),
		},

		"creation_source_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(databases.BackupStorageRedundancyGeo),
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForBackupStorageRedundancy(),
				false),
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"threat_detection_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disabled_alerts": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Sql_Injection",
								"Sql_Injection_Vulnerability",
								"Access_Anomaly",
							}, false),
						},
					},

					"email_account_admins_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"email_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"retention_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					// NOTE: I believe that this was originally implemented incorrect, where it exposed the
					// 'serveradvancedthreatprotectionsettings.PossibleValuesForAdvancedThreatProtectionState'
					// which contains the values of 'Enabled', 'Disabled', and 'New'
					// where 'serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState'
					// only contains 'Enabled' and 'Disabled'
					"state": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(serversecurityalertpolicies.SecurityAlertsPolicyStateDisabled),
						ValidateFunc: validation.StringInSlice(serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState(),
							false),
					},

					"storage_account_access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"geo_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"maintenance_configuration_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{"elastic_pool_id"},
			ValidateFunc:  validation.StringInSlice(helpers.ResourceMsSqlDatabaseMaintenanceNames(), false),
		},

		"ledger_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"transparent_data_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"transparent_data_encryption_key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"transparent_data_encryption_key_automatic_rotation_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			RequiredWith: []string{"transparent_data_encryption_key_vault_key_id"},
		},

		"secondary_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because it has defaulted to Geo for replicas but not all databases are replicas.
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForSecondaryType(), false),
		},

		"tags": commonschema.Tags(),
	}
}
