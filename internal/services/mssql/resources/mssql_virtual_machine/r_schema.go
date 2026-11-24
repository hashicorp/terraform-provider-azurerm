// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"math"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlVirtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"virtual_machine_id": {
			Description:  "The ID of the Virtual Machine.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"sql_license_type": {
			Description:  "The SQL Server license type.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForSqlServerLicenseType(), false),
		},

		"auto_backup": {
			Description: "This block can be added to an existing resource, but removing this block forces a new resource to be created.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"encryption_password": {
						Description:  "Encryption password to use. Setting a password will enable encryption.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"manual_schedule": {
						Description: "When this block is present, the schedule type is set to `Manual`. Without this block, the schedule type is set to `Automated`.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"full_backup_frequency": {
									Description:      "Frequency of full backups.",
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.CaseDifference,
									ValidateFunc:     validation.StringInSlice(sqlvirtualmachines.PossibleValuesForFullBackupFrequencyType(), false),
								},

								"full_backup_start_hour": {
									Description:  "Start hour of a given day during which full backups can take place.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 23),
								},

								"full_backup_window_in_hours": {
									Description:  "Duration of the time window of a given day during which full backups can take place, in hours.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 23),
								},

								"log_backup_frequency_in_minutes": {
									Description:  "Frequency of log backups, in minutes.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(5, 60),
								},

								"days_of_week": {
									Description: "A list of days on which backup can take place.",
									Type:        pluginsdk.TypeSet,
									Optional:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForAutoBackupDaysOfWeek(), false),
									},
								},
							},
						},
					},

					"retention_period_in_days": {
						Description:  "Retention period of backups, in days.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 30),
					},

					"storage_blob_endpoint": {
						Description:  "Blob endpoint for the Storage Account where backups will be kept.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},

					"storage_account_access_key": {
						Description:  "Access key for the Storage Account where backups will be kept.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsBase64,
					},

					"system_databases_backup_enabled": {
						Description: "Include or exclude system databases from auto backup.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"auto_patching": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Description: "The day of week to apply the patch on.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(sqlvirtualmachines.DayOfWeekMonday),
							string(sqlvirtualmachines.DayOfWeekTuesday),
							string(sqlvirtualmachines.DayOfWeekWednesday),
							string(sqlvirtualmachines.DayOfWeekThursday),
							string(sqlvirtualmachines.DayOfWeekFriday),
							string(sqlvirtualmachines.DayOfWeekSaturday),
							string(sqlvirtualmachines.DayOfWeekSunday),
						}, false),
					},

					"maintenance_window_duration_in_minutes": {
						Description:  "The size of the Maintenance Window in minutes.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(30, 180),
					},

					"maintenance_window_starting_hour": {
						Description:  "The Hour, in the Virtual Machine Time-Zone when the patching maintenance window should begin.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 23),
					},
				},
			},
		},

		"assessment": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Description: "Should Assessment be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"schedule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"weekly_interval": {
									Description:  "How many weeks between assessment runs.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ExactlyOneOf: []string{"assessment.0.schedule.0.monthly_occurrence"},
									ValidateFunc: validation.IntBetween(1, 6),
								},
								"monthly_occurrence": {
									Description:  "How many months between assessment runs.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ExactlyOneOf: []string{"assessment.0.schedule.0.weekly_interval"},
									ValidateFunc: validation.IntBetween(1, 5),
								},
								"day_of_week": {
									Description:  "What day of the week the assessment will be run.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForAssessmentDayOfWeek(), false),
								},
								"start_time": {
									Description: "What time the assessment will be run.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$"),
										"duration must match the format HH:mm",
									),
								},
							},
						},
					},

					"run_immediately": {
						Description: "Should Assessment be run immediately?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"key_vault_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The credential name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						// api will add updated credential name, and return "sqlvmName:name1,sqlvmName:name2"
						DiffSuppressFunc: mssqlVMCredentialNameDiffSuppressFunc,
					},

					"key_vault_url": {
						Description:  "The Azure Key Vault URL.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},

					"service_principal_name": {
						Description:  "The Service Principal name to access Key Vault.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"service_principal_secret": {
						Description:  "The Service Principal name secret to access Key Vault.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"r_services_enabled": {
			Description: "Should R Services be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"sql_connectivity_port": {
			Description:  "The SQL Server port.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1433,
			ValidateFunc: validation.IntBetween(1024, 65535),
		},

		"sql_connectivity_type": {
			Description:  "The connectivity type used for this SQL Server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(sqlvirtualmachines.ConnectivityTypePRIVATE),
			ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForConnectivityType(), false),
		},

		"sql_connectivity_update_password": {
			Description:  "The SQL Server sysadmin login password.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sql_connectivity_update_username": {
			Description:  "The SQL Server sysadmin login to create.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validate.SqlVirtualMachineLoginUserName,
		},

		"sql_instance": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"adhoc_workloads_optimization_enabled": {
						Description: "Whether the SQL Server is optimized for adhoc workloads.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"collation": {
						Description:  "Collation of the SQL Server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      "SQL_Latin1_General_CP1_CI_AS",
						ValidateFunc: validate.DatabaseCollation(),
					},

					"instant_file_initialization_enabled": {
						Description: "Whether Instant File Initialization is enabled for the SQL Server.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},

					"lock_pages_in_memory_enabled": {
						Description: "Whether Lock Pages in Memory is enabled for the SQL Server.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						ForceNew:    true,
						Default:     false,
					},

					"max_dop": {
						Description:  "Maximum Degree of Parallelism of the SQL Server.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 32767),
					},

					"max_server_memory_mb": {
						Description:  "Maximum amount memory that SQL Server Memory Manager can allocate to the SQL Server process.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      math.MaxInt32,
						ValidateFunc: validation.IntBetween(128, math.MaxInt32),
					},

					"min_server_memory_mb": {
						Description:  "Minimum amount memory that SQL Server Memory Manager can allocate to the SQL Server process.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},
				},
			},
		},

		"storage_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_type": {
						Description:  "The type of disk configuration to apply to the SQL Server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForDiskConfigurationType(), false),
					},
					"storage_workload_type": {
						Description:  "The type of storage workload.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForSqlWorkloadType(), false),
					},
					"system_db_on_data_disk_enabled": {
						Description: "Whether to set system databases (except tempDb) location to newly created data storage.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"data_settings":    StorageSettingSchema(),
					"log_settings":     StorageSettingSchema(),
					"temp_db_settings": SQLTempDBStorageSettingSchema(),
				},
			},
		},

		"sql_virtual_machine_group_id": {
			Description:  "The ID of the SQL Virtual Machine Group that the SQL Virtual Machine belongs to.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: sqlvirtualmachinegroups.ValidateSqlVirtualMachineGroupID,
		},

		"wsfc_domain_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cluster_bootstrap_account_password": {
						Description:  "The account password used for creating cluster.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"cluster_operator_account_password": {
						Description:  "The account password used for operating cluster.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sql_service_account_password": {
						Description:  "The account password under which SQL service will run on all participating SQL Virtual Machines in the cluster.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
