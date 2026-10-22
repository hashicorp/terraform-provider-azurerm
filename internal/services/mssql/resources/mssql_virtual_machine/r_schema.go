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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mssqlVirtualMachineSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"sql_license_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForSqlServerLicenseType(), false),
		},

		"auto_backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"encryption_password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"manual_schedule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"full_backup_frequency": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.CaseDifference,
									ValidateFunc:     validation.StringInSlice(sqlvirtualmachines.PossibleValuesForFullBackupFrequencyType(), false),
								},

								"full_backup_start_hour": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 23),
								},

								"full_backup_window_in_hours": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 23),
								},

								"log_backup_frequency_in_minutes": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(5, 60),
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForAutoBackupDaysOfWeek(), false),
									},
								},
							},
						},
					},

					"retention_period_in_days": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 30),
					},

					"storage_blob_endpoint": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},

					"storage_account_access_key": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsBase64,
					},

					"system_databases_backup_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
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
						Type:     pluginsdk.TypeString,
						Required: true,
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
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(30, 180),
					},

					"maintenance_window_starting_hour": {
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
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"schedule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"weekly_interval": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ExactlyOneOf: []string{"assessment.0.schedule.0.monthly_occurrence"},
									ValidateFunc: validation.IntBetween(1, 6),
								},
								"monthly_occurrence": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ExactlyOneOf: []string{"assessment.0.schedule.0.weekly_interval"},
									ValidateFunc: validation.IntBetween(1, 5),
								},
								"day_of_week": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForAssessmentDayOfWeek(), false),
								},
								"start_time": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$"),
										"duration must match the format HH:mm",
									),
								},
							},
						},
					},

					"run_immediately": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						// api will add updated credential name, and return "sqlvmName:name1,sqlvmName:name2"
						DiffSuppressFunc: mssqlVMCredentialNameDiffSuppressFunc,
					},

					"key_vault_url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},

					"service_principal_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"service_principal_secret": {
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"sql_connectivity_port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1433,
			ValidateFunc: validation.IntBetween(1024, 65535),
		},

		"sql_connectivity_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(sqlvirtualmachines.ConnectivityTypePRIVATE),
			ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForConnectivityType(), false),
		},

		"sql_connectivity_update_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sql_connectivity_update_username": {
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
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"collation": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      "SQL_Latin1_General_CP1_CI_AS",
						ValidateFunc: validate.DatabaseCollation(),
					},

					"instant_file_initialization_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"lock_pages_in_memory_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},

					"max_dop": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 32767),
					},

					"max_server_memory_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      math.MaxInt32,
						ValidateFunc: validation.IntBetween(128, math.MaxInt32),
					},

					"min_server_memory_mb": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForDiskConfigurationType(), false),
					},
					"storage_workload_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(sqlvirtualmachines.PossibleValuesForSqlWorkloadType(), false),
					},
					"system_db_on_data_disk_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"data_settings":    helper.StorageSettingSchema(),
					"log_settings":     helper.StorageSettingSchema(),
					"temp_db_settings": helper.SQLTempDBStorageSettingSchema(),
				},
			},
		},

		"sql_virtual_machine_group_id": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"cluster_operator_account_password": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"sql_service_account_password": {
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
