// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlVirtualMachine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlVirtualMachineCreateUpdate,
		Read:   resourceMsSqlVirtualMachineRead,
		Update: resourceMsSqlVirtualMachineCreateUpdate,
		Delete: resourceMsSqlVirtualMachineDelete,

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceMsSqlVirtualMachineCustomDiff),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := sqlvirtualmachines.ParseSqlVirtualMachineID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"virtual_machine_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualMachineID,
			},

			"sql_license_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachines.SqlServerLicenseTypePAYG),
					string(sqlvirtualmachines.SqlServerLicenseTypeAHUB),
					string(sqlvirtualmachines.SqlServerLicenseTypeDR),
				}, false),
			},

			"auto_backup": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

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
										ValidateFunc: validation.StringInSlice([]string{
											string(sqlvirtualmachines.FullBackupFrequencyTypeDaily),
											string(sqlvirtualmachines.FullBackupFrequencyTypeWeekly),
										}, false),
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
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												string(sqlvirtualmachines.AutoBackupDaysOfWeekMonday),
												string(sqlvirtualmachines.AutoBackupDaysOfWeekTuesday),
												string(sqlvirtualmachines.AutoBackupDaysOfWeekWednesday),
												string(sqlvirtualmachines.AutoBackupDaysOfWeekThursday),
												string(sqlvirtualmachines.AutoBackupDaysOfWeekFriday),
												string(sqlvirtualmachines.AutoBackupDaysOfWeekSaturday),
												string(sqlvirtualmachines.AutoBackupDaysOfWeekSunday),
											}, false),
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
							ValidateFunc: validation.StringIsNotEmpty,
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(sqlvirtualmachines.ConnectivityTypePRIVATE),
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachines.ConnectivityTypeLOCAL),
					string(sqlvirtualmachines.ConnectivityTypePRIVATE),
					string(sqlvirtualmachines.ConnectivityTypePUBLIC),
				}, false),
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
							Default:      2147483647,
							ValidateFunc: validation.IntBetween(128, 2147483647),
						},

						"min_server_memory_mb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 2147483647),
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachines.DiskConfigurationTypeNEW),
								string(sqlvirtualmachines.DiskConfigurationTypeEXTEND),
								string(sqlvirtualmachines.DiskConfigurationTypeADD),
							}, false),
						},
						"storage_workload_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachines.SqlWorkloadTypeGENERAL),
								string(sqlvirtualmachines.SqlWorkloadTypeOLTP),
								string(sqlvirtualmachines.SqlWorkloadTypeDW),
							}, false),
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
		},
	}
}

func resourceMsSqlVirtualMachineCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	// ForceNew when removing the auto_backup block.
	// See https://github.com/Azure/azure-rest-api-specs/issues/12818#issuecomment-773727756
	old, new := d.GetChange("auto_backup")
	if len(old.([]interface{})) == 1 && len(new.([]interface{})) == 0 {
		return d.ForceNew("auto_backup")
	}

	encryptionEnabled := d.Get("auto_backup.0.encryption_enabled")
	v, ok := d.GetOk("auto_backup.0.encryption_password")

	if encryptionEnabled.(bool) && (!ok || v.(string) == "") {
		return fmt.Errorf("auto_backup: `encryption_password` is required when `encryption_enabled` is true")
	}

	if !encryptionEnabled.(bool) && ok && v.(string) != "" {
		return fmt.Errorf("auto_backup: `encryption_enabled` must be true when `encryption_password` is set")
	}

	return nil
}

func resourceMsSqlVirtualMachineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	vmclient := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := virtualmachines.ParseVirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}
	id := sqlvirtualmachines.NewSqlVirtualMachineID(vmId.SubscriptionId, vmId.ResourceGroupName, vmId.VirtualMachineName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_machine", id.ID())
		}
	}

	// get location from vm
	respvm, err := vmclient.Get(ctx, *vmId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", vmId, err)
	}

	if respvm.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", vmId)
	}
	if respvm.Model.Location == "" {
		return fmt.Errorf("retrieving %s: `location` is empty", vmId)
	}
	sqlVmGroupId := ""
	if sqlVmGroupId = d.Get("sql_virtual_machine_group_id").(string); sqlVmGroupId != "" {
		parsedVmGroupId, err := sqlvirtualmachines.ParseSqlVirtualMachineGroupIDInsensitively(sqlVmGroupId)
		if err != nil {
			return err
		}
		sqlVmGroupId = parsedVmGroupId.ID()
	}

	sqlInstance, err := expandSqlVirtualMachineSQLInstance(d.Get("sql_instance").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `sql_instance`: %+v", err)
	}

	connectivityType := sqlvirtualmachines.ConnectivityType(d.Get("sql_connectivity_type").(string))
	sqlManagement := sqlvirtualmachines.SqlManagementModeFull
	sqlServerLicenseType := sqlvirtualmachines.SqlServerLicenseType(d.Get("sql_license_type").(string))
	autoBackupSettings, err := expandSqlVirtualMachineAutoBackupSettings(d.Get("auto_backup").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `auto_backup`: %+v", err)
	}

	parameters := sqlvirtualmachines.SqlVirtualMachine{
		Location: respvm.Model.Location,
		Properties: &sqlvirtualmachines.SqlVirtualMachineProperties{
			AutoBackupSettings:               autoBackupSettings,
			AutoPatchingSettings:             expandSqlVirtualMachineAutoPatchingSettings(d.Get("auto_patching").([]interface{})),
			AssessmentSettings:               expandSqlVirtualMachineAssessmentSettings(d.Get("assessment").([]interface{})),
			KeyVaultCredentialSettings:       expandSqlVirtualMachineKeyVaultCredential(d.Get("key_vault_credential").([]interface{})),
			WsfcDomainCredentials:            expandSqlVirtualMachineWsfcDomainCredentials(d.Get("wsfc_domain_credential").([]interface{})),
			SqlVirtualMachineGroupResourceId: pointer.To(sqlVmGroupId),
			ServerConfigurationsManagementSettings: &sqlvirtualmachines.ServerConfigurationsManagementSettings{
				AdditionalFeaturesServerConfigurations: &sqlvirtualmachines.AdditionalFeaturesServerConfigurations{
					IsRServicesEnabled: utils.Bool(d.Get("r_services_enabled").(bool)),
				},
				SqlConnectivityUpdateSettings: &sqlvirtualmachines.SqlConnectivityUpdateSettings{
					ConnectivityType:      &connectivityType,
					Port:                  utils.Int64(int64(d.Get("sql_connectivity_port").(int))),
					SqlAuthUpdatePassword: utils.String(d.Get("sql_connectivity_update_password").(string)),
					SqlAuthUpdateUserName: utils.String(d.Get("sql_connectivity_update_username").(string)),
				},
				SqlInstanceSettings: sqlInstance,
			},
			SqlManagement:                &sqlManagement,
			SqlServerLicenseType:         &sqlServerLicenseType,
			StorageConfigurationSettings: expandSqlVirtualMachineStorageConfigurationSettings(d.Get("storage_configuration").([]interface{})),
			VirtualMachineResourceId:     utils.String(d.Get("virtual_machine_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// Wait for the auto backup settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoBackup := d.Get("auto_backup"); (d.IsNewResource() && len(autoBackup.([]interface{})) > 0) || (!d.IsNewResource() && d.HasChange("auto_backup")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoBackupSettings to take effect", d.Id())
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoBackupSettings to take effect: %+v", d.Id(), err)
		}
	}

	// Wait for the auto patching settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoPatching := d.Get("auto_patching"); (d.IsNewResource() && len(autoPatching.([]interface{})) > 0) || (!d.IsNewResource() && d.HasChange("auto_patching")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoPatchingSettings to take effect", d.Id())
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoPatchingSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoPatchingSettings to take effect: %+v", d.Id(), err)
		}
	}

	return resourceMsSqlVirtualMachineRead(d, meta)
}

func resourceMsSqlVirtualMachineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("virtual_machine_id", props.VirtualMachineResourceId)
			sqlLicenseType := ""
			if licenceType := props.SqlServerLicenseType; licenceType != nil {
				sqlLicenseType = string(*licenceType)
			}
			d.Set("sql_license_type", sqlLicenseType)
			if err := d.Set("auto_backup", flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)); err != nil {
				return fmt.Errorf("setting `auto_backup`: %+v", err)
			}

			if err := d.Set("auto_patching", flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)); err != nil {
				return fmt.Errorf("setting `auto_patching`: %+v", err)
			}

			if err := d.Set("assessment", flattenSqlVirtualMachineAssessmentSettings(props.AssessmentSettings)); err != nil {
				return fmt.Errorf("setting `assessment`: %+v", err)
			}

			if err := d.Set("key_vault_credential", flattenSqlVirtualMachineKeyVaultCredential(props.KeyVaultCredentialSettings, d)); err != nil {
				return fmt.Errorf("setting `key_vault_credential`: %+v", err)
			}

			if mgmtSettings := props.ServerConfigurationsManagementSettings; mgmtSettings != nil {
				if cfgs := mgmtSettings.AdditionalFeaturesServerConfigurations; cfgs != nil {
					d.Set("r_services_enabled", mgmtSettings.AdditionalFeaturesServerConfigurations.IsRServicesEnabled)
				}
				if scus := mgmtSettings.SqlConnectivityUpdateSettings; scus != nil {
					d.Set("sql_connectivity_port", mgmtSettings.SqlConnectivityUpdateSettings.Port)
					d.Set("sql_connectivity_type", pointer.From(mgmtSettings.SqlConnectivityUpdateSettings.ConnectivityType))
				}

				d.Set("sql_instance", flattenSqlVirtualMachineSQLInstance(mgmtSettings.SqlInstanceSettings))
			}

			// `storage_configuration.0.storage_workload_type` is in a different spot than the rest of the `storage_configuration`
			// so we'll grab that here and pass it along
			storageWorkloadType := ""
			if props.ServerConfigurationsManagementSettings != nil && props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings != nil && props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings.SqlWorkloadType != nil {
				storageWorkloadType = string(*props.ServerConfigurationsManagementSettings.SqlWorkloadTypeUpdateSettings.SqlWorkloadType)
			}

			if err := d.Set("storage_configuration", flattenSqlVirtualMachineStorageConfigurationSettings(props.StorageConfigurationSettings, storageWorkloadType)); err != nil {
				return fmt.Errorf("setting `storage_configuration`: %+v", err)
			}

			sqlVirtualMachineGroupId := ""
			if props.SqlVirtualMachineGroupResourceId != nil {
				parsedId, err := sqlvirtualmachines.ParseSqlVirtualMachineGroupIDInsensitively(*props.SqlVirtualMachineGroupResourceId)
				if err != nil {
					return err
				}

				// get correct casing for subscription in id due to https://github.com/Azure/azure-rest-api-specs/issues/25211
				sqlVirtualMachineGroupId = sqlvirtualmachines.NewSqlVirtualMachineGroupID(id.SubscriptionId, parsedId.ResourceGroupName, parsedId.SqlVirtualMachineGroupName).ID()
			}
			d.Set("sql_virtual_machine_group_id", sqlVirtualMachineGroupId)

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceMsSqlVirtualMachineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachines.SqlVirtualMachinesClient, d *pluginsdk.ResourceData) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
		if err != nil {
			return nil, "Retry", err
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				autoBackupSettings := flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)

				if len(autoBackupSettings) == 0 {
					// auto backup was nil or disabled in the response
					if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]interface{})) == 0 {
						// also disabled in the config
						return resp, "Updated", nil
					}
					return resp, "Pending", nil
				}

				if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]interface{})) == 0 {
					// still waiting for it to be disabled
					return resp, "Pending", nil
				}

				// check each property in the auto_backup block for drift
				for prop, val := range autoBackupSettings[0].(map[string]interface{}) {
					v := d.Get(fmt.Sprintf("auto_backup.0.%s", prop))
					switch prop {
					case "manual_schedule":
						if m := val.([]interface{}); len(m) > 0 {
							if b, ok := d.GetOk("auto_backup.0.manual_schedule"); !ok || len(b.([]interface{})) == 0 {
								// manual schedule disabled in config but still showing in response
								return resp, "Pending", nil
							}
							// check each property in the manual_schedule block for drift
							for prop2, val2 := range m[0].(map[string]interface{}) {
								v2 := d.Get(fmt.Sprintf("auto_backup.0.manual_schedule.0.%s", prop2))
								switch prop2 {
								case "full_backup_frequency":
									if !strings.EqualFold(v2.(string), val2.(string)) {
										return resp, "Pending", nil
									}
								case "days_of_week":
									daysOfWeekLocal := make(map[string]bool, 0)
									if v2 != nil {
										for _, item := range v2.(*pluginsdk.Set).List() {
											daysOfWeekLocal[item.(string)] = true
										}
									}

									daysOfWeekRemote := make([]interface{}, 0)
									if val2 != nil {
										daysOfWeekRemote = val2.([]interface{})
									}

									if len(daysOfWeekRemote) != len(daysOfWeekLocal) {
										return resp, "Pending", nil
									}

									for _, item := range daysOfWeekRemote {
										if _, ok := daysOfWeekLocal[item.(string)]; !ok {
											return resp, "Pending", nil
										}
									}
								default:
									if v2 != val2 {
										return resp, "Pending", nil
									}
								}
							}
						} else if b, ok := d.GetOk("auto_backup.0.manual_schedule"); ok || len(b.([]interface{})) > 0 {
							// manual schedule set in config but not reflecting in response
							return resp, "Pending", nil
						}
					default:
						if v != val {
							return resp, "Pending", nil
						}
					}
				}

				return resp, "Updated", nil
			}
		}

		return resp, "Retry", nil
	}
}

func expandSqlVirtualMachineAutoBackupSettings(input []interface{}) (*sqlvirtualmachines.AutoBackupSettings, error) {
	ret := sqlvirtualmachines.AutoBackupSettings{
		Enable: utils.Bool(false),
	}

	if len(input) > 0 {
		config := input[0].(map[string]interface{})
		ret.Enable = utils.Bool(true)

		if v, ok := config["retention_period_in_days"]; ok {
			ret.RetentionPeriod = utils.Int64(int64(v.(int)))
		}
		if v, ok := config["storage_blob_endpoint"]; ok {
			ret.StorageAccountURL = utils.String(v.(string))
		}
		if v, ok := config["storage_account_access_key"]; ok {
			ret.StorageAccessKey = utils.String(v.(string))
		}

		v, ok := config["encryption_enabled"]
		enableEncryption := ok && v.(bool)
		ret.EnableEncryption = utils.Bool(enableEncryption)
		if v, ok := config["encryption_password"]; enableEncryption && ok {
			ret.Password = utils.String(v.(string))
		}

		if v, ok := config["system_databases_backup_enabled"]; ok {
			ret.BackupSystemDbs = utils.Bool(v.(bool))
		}

		backupScheduleTypeAutomated := sqlvirtualmachines.BackupScheduleTypeAutomated
		ret.BackupScheduleType = &backupScheduleTypeAutomated
		if v, ok := config["manual_schedule"]; ok && len(v.([]interface{})) > 0 {
			manualSchedule := v.([]interface{})[0].(map[string]interface{})
			backupScheduleTypeManual := sqlvirtualmachines.BackupScheduleTypeManual
			ret.BackupScheduleType = &backupScheduleTypeManual

			fullBackupFrequency := sqlvirtualmachines.FullBackupFrequencyType(manualSchedule["full_backup_frequency"].(string))

			daysOfWeek := manualSchedule["days_of_week"].(*pluginsdk.Set).List()
			if len(daysOfWeek) > 0 {
				if !strings.EqualFold(string(fullBackupFrequency), string(sqlvirtualmachines.FullBackupFrequencyTypeWeekly)) {
					return nil, fmt.Errorf("`manual_schedule.0.days_of_week` can only be specified when `manual_schedule.0.full_backup_frequency` is set to `Weekly`")
				}
				ret.DaysOfWeek = expandSqlVirtualMachineAutoBackupSettingsDaysOfWeek(daysOfWeek)
			}

			ret.FullBackupFrequency = &fullBackupFrequency
			ret.FullBackupStartTime = utils.Int64(int64(manualSchedule["full_backup_start_hour"].(int)))
			ret.FullBackupWindowHours = utils.Int64(int64(manualSchedule["full_backup_window_in_hours"].(int)))
			ret.LogBackupFrequency = utils.Int64(int64(manualSchedule["log_backup_frequency_in_minutes"].(int)))
		}
	}

	return &ret, nil
}

func flattenSqlVirtualMachineAutoBackup(autoBackup *sqlvirtualmachines.AutoBackupSettings, d *pluginsdk.ResourceData) []interface{} {
	if autoBackup == nil || autoBackup.Enable == nil || !*autoBackup.Enable {
		return []interface{}{}
	}

	manualSchedule := make([]interface{}, 0)
	if autoBackup.BackupScheduleType != nil && strings.EqualFold(string(*autoBackup.BackupScheduleType), string(sqlvirtualmachines.BackupScheduleTypeManual)) {
		var fullBackupStartHour int
		if autoBackup.FullBackupStartTime != nil {
			fullBackupStartHour = int(*autoBackup.FullBackupStartTime)
		}

		var fullBackupWindowHours int
		if autoBackup.FullBackupWindowHours != nil {
			fullBackupWindowHours = int(*autoBackup.FullBackupWindowHours)
		}

		var logBackupFrequency int
		if autoBackup.LogBackupFrequency != nil {
			logBackupFrequency = int(*autoBackup.LogBackupFrequency)
			// API returns 60 minutes as zero
			if logBackupFrequency == 0 {
				logBackupFrequency = 60
			}
		}

		var fullBackupFrequency string
		if autoBackup.FullBackupFrequency != nil {
			fullBackupFrequency = string(*autoBackup.FullBackupFrequency)
		}

		manualSchedule = []interface{}{
			map[string]interface{}{
				"full_backup_frequency":           fullBackupFrequency,
				"full_backup_start_hour":          fullBackupStartHour,
				"full_backup_window_in_hours":     fullBackupWindowHours,
				"log_backup_frequency_in_minutes": logBackupFrequency,
				"days_of_week":                    flattenSqlVirtualMachineAutoBackupDaysOfWeek(autoBackup.DaysOfWeek),
			},
		}
	}

	var retentionPeriod int
	if autoBackup.RetentionPeriod != nil {
		retentionPeriod = int(*autoBackup.RetentionPeriod)
	}

	// Password, StorageAccessKey, StorageAccountURL are not returned, so we try to copy them
	// from existing config as a best effort.
	encryptionPassword := d.Get("auto_backup.0.encryption_password").(string)
	storageKey := d.Get("auto_backup.0.storage_account_access_key").(string)
	blobEndpoint := d.Get("auto_backup.0.storage_blob_endpoint").(string)

	return []interface{}{
		map[string]interface{}{
			"encryption_enabled":              autoBackup.EnableEncryption != nil && *autoBackup.EnableEncryption,
			"encryption_password":             encryptionPassword,
			"manual_schedule":                 manualSchedule,
			"retention_period_in_days":        retentionPeriod,
			"storage_account_access_key":      storageKey,
			"storage_blob_endpoint":           blobEndpoint,
			"system_databases_backup_enabled": autoBackup.BackupSystemDbs != nil && *autoBackup.BackupSystemDbs,
		},
	}
}

func expandSqlVirtualMachineAutoBackupSettingsDaysOfWeek(input []interface{}) *[]sqlvirtualmachines.AutoBackupDaysOfWeek {
	result := make([]sqlvirtualmachines.AutoBackupDaysOfWeek, 0)
	for _, item := range input {
		result = append(result, sqlvirtualmachines.AutoBackupDaysOfWeek(item.(string)))
	}
	return &result
}

func flattenSqlVirtualMachineAutoBackupDaysOfWeek(daysOfWeek *[]sqlvirtualmachines.AutoBackupDaysOfWeek) []interface{} {
	output := make([]interface{}, 0)

	if daysOfWeek != nil {
		for _, v := range *daysOfWeek {
			output = append(output, string(v))
		}
	}

	return output
}

func resourceMsSqlVirtualMachineAutoPatchingSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachines.SqlVirtualMachinesClient, d *pluginsdk.ResourceData) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: utils.String("*")})
		if err != nil {
			return nil, "Retry", err
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				autoPatchingSettings := flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)

				if len(autoPatchingSettings) == 0 {
					if v, ok := d.GetOk("auto_patching"); !ok || len(v.([]interface{})) == 0 {
						return resp, "Updated", nil
					}
					return resp, "Pending", nil
				}

				if v, ok := d.GetOk("auto_patching"); !ok || len(v.([]interface{})) == 0 {
					return resp, "Pending", nil
				}

				for prop, val := range autoPatchingSettings[0].(map[string]interface{}) {
					v := d.Get(fmt.Sprintf("auto_patching.0.%s", prop))
					if v != val {
						return resp, "Pending", nil
					}
				}

				return resp, "Updated", nil
			}
		}

		return resp, "Retry", nil
	}
}

func expandSqlVirtualMachineAutoPatchingSettings(input []interface{}) *sqlvirtualmachines.AutoPatchingSettings {
	if len(input) == 0 {
		return nil
	}
	autoPatchingSetting := input[0].(map[string]interface{})

	dayOfWeek := sqlvirtualmachines.DayOfWeek(autoPatchingSetting["day_of_week"].(string))

	return &sqlvirtualmachines.AutoPatchingSettings{
		Enable:                        utils.Bool(true),
		MaintenanceWindowDuration:     utils.Int64(int64(autoPatchingSetting["maintenance_window_duration_in_minutes"].(int))),
		MaintenanceWindowStartingHour: utils.Int64(int64(autoPatchingSetting["maintenance_window_starting_hour"].(int))),
		DayOfWeek:                     &dayOfWeek,
	}
}

func flattenSqlVirtualMachineAutoPatching(autoPatching *sqlvirtualmachines.AutoPatchingSettings) []interface{} {
	if autoPatching == nil || autoPatching.Enable == nil || !*autoPatching.Enable {
		return []interface{}{}
	}

	var startHour int
	if autoPatching.MaintenanceWindowStartingHour != nil {
		startHour = int(*autoPatching.MaintenanceWindowStartingHour)
	}

	var duration int
	if autoPatching.MaintenanceWindowDuration != nil {
		duration = int(*autoPatching.MaintenanceWindowDuration)
	}

	var dayOfWeek string
	if autoPatching.DayOfWeek != nil {
		dayOfWeek = string(*autoPatching.DayOfWeek)
	}

	return []interface{}{
		map[string]interface{}{
			"day_of_week":                            dayOfWeek,
			"maintenance_window_starting_hour":       startHour,
			"maintenance_window_duration_in_minutes": duration,
		},
	}
}

func expandSqlVirtualMachineAssessmentSettings(input []interface{}) *sqlvirtualmachines.AssessmentSettings {
	if len(input) == 0 {
		return nil
	}
	assessmentSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachines.AssessmentSettings{
		Enable:         utils.Bool(true),
		RunImmediately: utils.Bool(assessmentSetting["run_immediately"].(bool)),
		Schedule:       expandSqlVirtualMachineAssessmentSettingsSchedule(assessmentSetting["schedule"].([]interface{})),
	}
}

func expandSqlVirtualMachineAssessmentSettingsSchedule(input []interface{}) *sqlvirtualmachines.Schedule {
	if len(input) == 0 {
		return &sqlvirtualmachines.Schedule{}
	}

	scheduleConfig := input[0].(map[string]interface{})

	dayOfWeek := sqlvirtualmachines.AssessmentDayOfWeek(scheduleConfig["day_of_week"].(string))

	schedule := &sqlvirtualmachines.Schedule{
		Enable:    utils.Bool(true),
		DayOfWeek: &dayOfWeek,
		StartTime: utils.String(scheduleConfig["start_time"].(string)),
	}

	if weeklyInterval := scheduleConfig["weekly_interval"].(int); weeklyInterval != 0 {
		schedule.WeeklyInterval = utils.Int64(int64(weeklyInterval))
	}

	if monthlyOccurrence := scheduleConfig["monthly_occurrence"].(int); monthlyOccurrence != 0 {
		schedule.MonthlyOccurrence = utils.Int64(int64(monthlyOccurrence))
	}

	return schedule
}

func flattenSqlVirtualMachineAssessmentSettings(assessmentSettings *sqlvirtualmachines.AssessmentSettings) []interface{} {
	if assessmentSettings == nil || assessmentSettings.Enable == nil || !*assessmentSettings.Enable {
		return []interface{}{}
	}

	var (
		runImmediately bool
		enabled        bool
	)
	if assessmentSettings.RunImmediately != nil {
		runImmediately = *assessmentSettings.RunImmediately
	}

	if assessmentSettings.Enable != nil {
		enabled = *assessmentSettings.Enable
	}

	var attr map[string]interface{}
	if schedule := assessmentSettings.Schedule; schedule != nil {
		var (
			weeklyInterval    int64
			monthlyOccurrence int64
			dayOfWeek         string
			startTime         string
		)

		if schedule.WeeklyInterval != nil {
			weeklyInterval = *schedule.WeeklyInterval
		}
		if schedule.MonthlyOccurrence != nil {
			monthlyOccurrence = *schedule.MonthlyOccurrence
		}
		if schedule.DayOfWeek != nil {
			dayOfWeek = string(*schedule.DayOfWeek)
		}
		if schedule.StartTime != nil {
			startTime = *schedule.StartTime
		}

		attr = map[string]interface{}{
			"weekly_interval":    weeklyInterval,
			"monthly_occurrence": monthlyOccurrence,
			"day_of_week":        dayOfWeek,
			"start_time":         startTime,
		}
	}

	return []interface{}{
		map[string]interface{}{
			"run_immediately": runImmediately,
			"enabled":         enabled,
			"schedule":        []interface{}{attr},
		},
	}
}

func expandSqlVirtualMachineKeyVaultCredential(input []interface{}) *sqlvirtualmachines.KeyVaultCredentialSettings {
	if len(input) == 0 {
		return nil
	}
	keyVaultCredentialSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachines.KeyVaultCredentialSettings{
		Enable:                 utils.Bool(true),
		CredentialName:         utils.String(keyVaultCredentialSetting["name"].(string)),
		AzureKeyVaultURL:       utils.String(keyVaultCredentialSetting["key_vault_url"].(string)),
		ServicePrincipalName:   utils.String(keyVaultCredentialSetting["service_principal_name"].(string)),
		ServicePrincipalSecret: utils.String(keyVaultCredentialSetting["service_principal_secret"].(string)),
	}
}

func flattenSqlVirtualMachineKeyVaultCredential(keyVault *sqlvirtualmachines.KeyVaultCredentialSettings, d *pluginsdk.ResourceData) []interface{} {
	if keyVault == nil || keyVault.Enable == nil || !*keyVault.Enable {
		return []interface{}{}
	}

	name := ""
	if keyVault.CredentialName != nil {
		name = *keyVault.CredentialName
	}

	keyVaultUrl := ""
	if v, ok := d.GetOk("key_vault_credential.0.key_vault_url"); ok {
		keyVaultUrl = v.(string)
	}

	servicePrincipalName := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_name"); ok {
		servicePrincipalName = v.(string)
	}

	servicePrincipalSecret := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_secret"); ok {
		servicePrincipalSecret = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"name":                     name,
			"key_vault_url":            keyVaultUrl,
			"service_principal_name":   servicePrincipalName,
			"service_principal_secret": servicePrincipalSecret,
		},
	}
}

func mssqlVMCredentialNameDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	oldNamelist := strings.Split(old, ",")
	for _, n := range oldNamelist {
		cur := strings.Split(n, ":")
		if len(cur) > 1 && cur[1] == new {
			return true
		}
	}
	return false
}

func expandSqlVirtualMachineStorageConfigurationSettings(input []interface{}) *sqlvirtualmachines.StorageConfigurationSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	storageSettings := input[0].(map[string]interface{})

	diskConfigurationType := sqlvirtualmachines.DiskConfigurationType(storageSettings["disk_type"].(string))
	storageWorkloadType := sqlvirtualmachines.StorageWorkloadType(storageSettings["storage_workload_type"].(string))

	return &sqlvirtualmachines.StorageConfigurationSettings{
		DiskConfigurationType: &diskConfigurationType,
		StorageWorkloadType:   &storageWorkloadType,
		SqlSystemDbOnDataDisk: utils.Bool(storageSettings["system_db_on_data_disk_enabled"].(bool)),
		SqlDataSettings:       expandSqlVirtualMachineDataStorageSettings(storageSettings["data_settings"].([]interface{})),
		SqlLogSettings:        expandSqlVirtualMachineDataStorageSettings(storageSettings["log_settings"].([]interface{})),
		SqlTempDbSettings:     expandSqlVirtualMachineTempDbSettings(storageSettings["temp_db_settings"].([]interface{})),
	}
}

func flattenSqlVirtualMachineStorageConfigurationSettings(input *sqlvirtualmachines.StorageConfigurationSettings, storageWorkloadType string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var diskType string
	if input.DiskConfigurationType != nil {
		diskType = string(*input.DiskConfigurationType)
	}

	systemDbOnDataDisk := false
	if input.SqlSystemDbOnDataDisk != nil {
		systemDbOnDataDisk = *input.SqlSystemDbOnDataDisk
	}

	output := map[string]interface{}{
		"storage_workload_type":          storageWorkloadType,
		"disk_type":                      diskType,
		"system_db_on_data_disk_enabled": systemDbOnDataDisk,
		"data_settings":                  flattenSqlVirtualMachineStorageSettings(input.SqlDataSettings),
		"log_settings":                   flattenSqlVirtualMachineStorageSettings(input.SqlLogSettings),
		"temp_db_settings":               flattenSqlVirtualMachineTempDbSettings(input.SqlTempDbSettings),
	}

	if output["storage_workload_type"].(string) == "" && output["disk_type"] == "" &&
		len(output["data_settings"].([]interface{})) == 0 &&
		len(output["log_settings"].([]interface{})) == 0 &&
		len(output["temp_db_settings"].([]interface{})) == 0 {
		return []interface{}{}
	}

	return []interface{}{output}
}

func expandSqlVirtualMachineDataStorageSettings(input []interface{}) *sqlvirtualmachines.SQLStorageSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	dataStorageSettings := input[0].(map[string]interface{})

	return &sqlvirtualmachines.SQLStorageSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(dataStorageSettings["luns"].([]interface{})),
		DefaultFilePath: utils.String(dataStorageSettings["default_file_path"].(string)),
	}
}

func expandSqlVirtualMachineStorageSettingsLuns(input []interface{}) *[]int64 {
	expandedLuns := make([]int64, 0)
	for i := range input {
		if input[i] != nil {
			expandedLuns = append(expandedLuns, int64(input[i].(int)))
		}
	}

	return &expandedLuns
}

func flattenSqlVirtualMachineStorageSettings(input *sqlvirtualmachines.SQLStorageSettings) []interface{} {
	if input == nil || input.Luns == nil {
		return []interface{}{}
	}
	attrs := make(map[string]interface{})

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	return []interface{}{attrs}
}

func expandSqlVirtualMachineTempDbSettings(input []interface{}) *sqlvirtualmachines.SQLTempDbSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	tempDbSettings := input[0].(map[string]interface{})

	return &sqlvirtualmachines.SQLTempDbSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(tempDbSettings["luns"].([]interface{})),
		DefaultFilePath: utils.String(tempDbSettings["default_file_path"].(string)),
		DataFileCount:   utils.Int64(int64(tempDbSettings["data_file_count"].(int))),
		DataFileSize:    utils.Int64(int64(tempDbSettings["data_file_size_mb"].(int))),
		DataGrowth:      utils.Int64(int64(tempDbSettings["data_file_growth_in_mb"].(int))),
		LogFileSize:     utils.Int64(int64(tempDbSettings["log_file_size_mb"].(int))),
		LogGrowth:       utils.Int64(int64(tempDbSettings["log_file_growth_mb"].(int))),
	}
}

func flattenSqlVirtualMachineTempDbSettings(input *sqlvirtualmachines.SQLTempDbSettings) []interface{} {
	if input == nil || input.Luns == nil {
		return []interface{}{}
	}
	attrs := make(map[string]interface{})

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DataFileCount != nil {
		attrs["data_file_count"] = *input.DataFileCount
	}

	if input.DataFileSize != nil {
		attrs["data_file_size_mb"] = *input.DataFileSize
	}

	if input.DataGrowth != nil {
		attrs["data_file_growth_in_mb"] = *input.DataGrowth
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	if input.LogFileSize != nil {
		attrs["log_file_size_mb"] = *input.LogFileSize
	}

	if input.LogGrowth != nil {
		attrs["log_file_growth_mb"] = *input.LogGrowth
	}

	return []interface{}{attrs}
}
func expandSqlVirtualMachineSQLInstance(input []interface{}) (*sqlvirtualmachines.SQLInstanceSettings, error) {
	if len(input) == 0 || input[0] == nil {
		return &sqlvirtualmachines.SQLInstanceSettings{}, nil
	}

	settings := input[0].(map[string]interface{})
	maxServerMemoryMB := settings["max_server_memory_mb"].(int)
	minServerMemoryMB := settings["min_server_memory_mb"].(int)

	if maxServerMemoryMB < minServerMemoryMB {
		return nil, fmt.Errorf("`max_server_memory_mb` must be greater than or equal to `min_server_memory_mb`")
	}

	result := sqlvirtualmachines.SQLInstanceSettings{
		Collation:                          utils.String(settings["collation"].(string)),
		IsIfiEnabled:                       utils.Bool(settings["instant_file_initialization_enabled"].(bool)),
		IsLpimEnabled:                      utils.Bool(settings["lock_pages_in_memory_enabled"].(bool)),
		IsOptimizeForAdHocWorkloadsEnabled: utils.Bool(settings["adhoc_workloads_optimization_enabled"].(bool)),
		MaxDop:                             utils.Int64(int64(settings["max_dop"].(int))),
		MaxServerMemoryMB:                  utils.Int64(int64(maxServerMemoryMB)),
		MinServerMemoryMB:                  utils.Int64(int64(minServerMemoryMB)),
	}

	return &result, nil
}

func flattenSqlVirtualMachineSQLInstance(input *sqlvirtualmachines.SQLInstanceSettings) []interface{} {
	if input == nil || input.Collation == nil {
		return []interface{}{}
	}

	collation := *input.Collation

	isIfiEnabled := false
	if input.IsIfiEnabled != nil {
		isIfiEnabled = *input.IsIfiEnabled
	}

	isLpimEnabled := false
	if input.IsLpimEnabled != nil {
		isLpimEnabled = *input.IsLpimEnabled
	}

	isOptimizeForAdhocWorkloadsEnabled := false
	if input.IsOptimizeForAdHocWorkloadsEnabled != nil {
		isOptimizeForAdhocWorkloadsEnabled = *input.IsOptimizeForAdHocWorkloadsEnabled
	}

	var maxDop int64 = 0
	if input.MaxDop != nil {
		maxDop = *input.MaxDop
	}

	var maxServerMemoryMB int64 = 2147483647
	if input.MaxServerMemoryMB != nil {
		maxServerMemoryMB = *input.MaxServerMemoryMB
	}

	var minServerMemoryMB int64 = 0
	if input.MinServerMemoryMB != nil {
		minServerMemoryMB = *input.MinServerMemoryMB
	}

	return []interface{}{
		map[string]interface{}{
			"adhoc_workloads_optimization_enabled": isOptimizeForAdhocWorkloadsEnabled,
			"collation":                            collation,
			"instant_file_initialization_enabled":  isIfiEnabled,
			"lock_pages_in_memory_enabled":         isLpimEnabled,
			"max_dop":                              maxDop,
			"max_server_memory_mb":                 maxServerMemoryMB,
			"min_server_memory_mb":                 minServerMemoryMB,
		},
	}
}

func expandSqlVirtualMachineWsfcDomainCredentials(input []interface{}) *sqlvirtualmachines.WsfcDomainCredentials {
	if len(input) == 0 {
		return nil
	}
	wsfcDomainCredentials := input[0].(map[string]interface{})

	return &sqlvirtualmachines.WsfcDomainCredentials{
		ClusterBootstrapAccountPassword: pointer.To(wsfcDomainCredentials["cluster_bootstrap_account_password"].(string)),
		ClusterOperatorAccountPassword:  pointer.To(wsfcDomainCredentials["cluster_operator_account_password"].(string)),
		SqlServiceAccountPassword:       pointer.To(wsfcDomainCredentials["sql_service_account_password"].(string)),
	}
}
