// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceMsSqlDatabaseSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the MS SQL Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlDatabaseName,
		},

		"server_id": {
			Description:  "The ID of the MS SQL Server on which to create the database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlServerIDInsensitively),
		},

		"auto_pause_delay_in_minutes": {
			Description:  "Time in minutes after which database is automatically paused. A value of `-1` means that automatic pause is disabled. This property is only settable for Serverless databases.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.DatabaseAutoPauseDelay,
		},

		"create_mode": {
			Description: "The create mode of the database. Mutually exclusive with `import`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(databases.CreateModeDefault),
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForCreateMode(),
				false),
			ConflictsWith: []string{"import"},
		},

		"import": {
			Description:   "Mutually exclusive with `create_mode`.",
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"create_mode"}, // it needs `create_mode` to be `Default` to work, so make them conflict.
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_uri": {
						Description: "The blob URI of the .bacpac file.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"storage_key": {
						Description: "The access key for the Storage Account.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},
					"storage_key_type": {
						Description: "The type of access key for the Storage Account.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(databases.StorageKeyTypeSharedAccessKey),
							string(databases.StorageKeyTypeStorageAccessKey),
						}, false),
					},
					"administrator_login": {
						Description: "The name of the SQL administrator.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"administrator_login_password": {
						Description: "The password of the SQL administrator.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},
					"authentication_type": {
						Description: "The type of authentication used to access the server.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"ADPassword",
							"Sql",
						}, false),
					},
					"storage_account_id": {
						Description:  "The resource ID for the Storage Account used to store BACPAC file. If set, Private Endpoint connection will be created for the Storage Account. Must match Storage Account used for storage_uri parameter.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
				},
			},
		},

		"collation": {
			Description:  "The collation of the database.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.DatabaseCollation(),
		},

		"elastic_pool_id": {
			Description:  "The ID of the elastic pool containing this database.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseSqlElasticPoolIDInsensitively),
		},

		"enclave_type": {
			Description:  "The type of enclave to be used by the elastic pool. When `enclave_type` is not specified (e.g., the default) enclaves are not enabled on the database. Once enabled (e.g., by specifying `Default` or `VBS`) removing the `enclave_type` field from the configuration file will force the creation of a new resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForAlwaysEncryptedEnclaveType(), false),
		},

		"license_type": {
			Description: "The license type applied to this database.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForDatabaseLicenseType(),
				false),
		},

		"long_term_retention_policy": LongTermRetentionPolicySchema(),

		"short_term_retention_policy": ShortTermRetentionPolicySchema(),

		"max_size_gb": {
			Description:  "The max size of the database in gigabytes.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatBetween(0.1, 4096),
		},

		"min_capacity": {
			Description:  "Minimal capacity that database will always have allocated, if not paused. This property is only settable for Serverless databases.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.FloatInSlice([]float64{0, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40}),
		},

		"restore_point_in_time": {
			Description:      "The point in time (ISO8601 format) of the source database that will be restored to create the new database. This property is only settable for `create_mode`= `PointInTimeRestore` databases.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ValidateFunc:     validation.IsRFC3339Time,
			DiffSuppressFunc: suppress.RFC3339Time,
		},

		"recover_database_id": {
			Description:  "The ID of the database to be recovered. This property is only applicable when the `create_mode` is `Recovery`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RecoverableDatabaseID,
		},

		"recovery_point_id": {
			Description:  "The ID of the Recovery Services Recovery Point ID to be restored. This property is only applicable when the `create_mode` is `Recovery`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"restore_dropped_database_id": {
			Description:  "The ID of the database to be restored. This property is only applicable when the `create_mode` is `Restore`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RestorableDatabaseID,
		},

		"restore_long_term_retention_backup_id": {
			Description:  "The ID of the long term retention backup to be restored. This property is only applicable when the `create_mode` is `RestoreLongTermRetentionBackup`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"read_replica_count": {
			Description:  "The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed. This property is only settable for Hyperscale edition databases.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"read_scale": {
			Description: "If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica. This property is only settable for Premium and Business Critical databases.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"sample_name": {
			Description: "The name of the sample schema to apply when creating this database.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(databases.SampleNameAdventureWorksLT),
			}, false),
		},

		"sku_name": {
			Description:  "The name of the SKU used by the database. For example, `GP_S_Gen5_2`, `HS_Gen4_1`, `BC_Gen5_2`, `ElasticPool`, `Basic`, `S0`, `P2`, `DW100c`, `DS100`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.DatabaseSkuName(),
		},

		"creation_source_database_id": {
			Description:  "The ID of the source database from which to create the new database. This should only be used for databases with `create_mode` values that use another database as reference.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"storage_account_type": {
			Description: "The Storage Account type used to store backups for this database.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(databases.BackupStorageRedundancyGeo),
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForBackupStorageRedundancy(),
				false),
		},

		"zone_redundant": {
			Description: "Whether this database is zone redundant, which means the replicas of this database will be spread across multiple Availability Zones. This property is only settable for Premium and Business Critical databases.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"threat_detection_policy": {
			Description: "Threat detection policy configuration.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disabled_alerts": {
						Description: "A list of alerts which should be disabled.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashString,
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
						Description: "Whether the account administrators should receive an email when this alert is triggered.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"email_addresses": {
						Description: "A list of email addresses which alerts should be sent to.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"retention_days": {
						Description:  "The number of days to keep in the Threat Detection audit logs.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"enabled": {
						Description: "Whether the policy is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"storage_account_access_key": {
						Description:  "The identifier key of the Threat Detection audit Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_endpoint": {
						Description:  "The blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"geo_backup_enabled": {
			Description: "A boolean that specifies if the Geo Backup Policy is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"maintenance_configuration_name": {
			Description:   "The name of the Public Maintenance Configuration window to apply to the database.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ValidateFunc:  validation.StringInSlice(helpers.ResourceMsSqlDatabaseMaintenanceNames(), false),
			ConflictsWith: []string{"elastic_pool_id"},
		},

		"ledger_enabled": {
			Description: "A boolean that specifies if this is a ledger database.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"transparent_data_encryption_enabled": {
			Description: "If set to true, Transparent Data Encryption will be enabled on the database.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"transparent_data_encryption_key_vault_key_id": {
			Description:  "The fully versioned `Key Vault` `Key` URL (e.g. `'https://<YourVaultName>.vault.azure.net/keys/<YourKeyName>/<YourKeyVersion>`) to be used as the `Customer Managed Key` (CMK/BYOK) for the `Transparent Data Encryption` (TDE) layer.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
		},

		"transparent_data_encryption_key_automatic_rotation_enabled": {
			Description:  "Boolean flag to specify whether TDE automatically rotates the encryption Key to latest version.",
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			RequiredWith: []string{"transparent_data_encryption_key_vault_key_id"},
		},

		"secondary_type": {
			Description: "How the replica of this database is made.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because it has defaulted to Geo for replicas but not all databases are replicas.
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForSecondaryType(), false),
		},

		"tags": commonschema.Tags(),
	}

	if !features.SixPointOh() {
		threatDetection := s["threat_detection_policy"].Elem.(*pluginsdk.Resource).Schema

		threatDetection["state"] = &pluginsdk.Schema{
			Description:   "The State of the Policy.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ValidateFunc:  validation.StringInSlice(serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState(), false),
			ConflictsWith: []string{"threat_detection_policy.0.enabled"},
			Deprecated:    "`threat_detection_policy.state` has been deprecated in favour of `threat_detection_policy.enabled` and will be removed in v6.0 of the AzureRM Provider",
		}

		threatDetection["enabled"] = &pluginsdk.Schema{
			Description:   "Whether the policy is enabled.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"threat_detection_policy.0.state"},
		}
	}

	return s
}
