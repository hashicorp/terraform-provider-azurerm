// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type HandlerMappings struct {
	Extension           string `tfschema:"extension"`
	ScriptProcessorPath string `tfschema:"script_processor_path"`
	Arguments           string `tfschema:"arguments"`
}

func HandlerMappingSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"extension": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"script_processor_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"arguments": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func HandlerMappingSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"extension": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"script_processor_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"arguments": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type VirtualApplication struct {
	VirtualPath        string             `tfschema:"virtual_path"`
	PhysicalPath       string             `tfschema:"physical_path"`
	Preload            bool               `tfschema:"preload"`
	VirtualDirectories []VirtualDirectory `tfschema:"virtual_directory"`
}

type VirtualDirectory struct {
	VirtualPath  string `tfschema:"virtual_path"`
	PhysicalPath string `tfschema:"physical_path"`
}

func virtualApplicationsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"physical_path": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"preload": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},

				"virtual_directory": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"virtual_path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"physical_path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func virtualApplicationsSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"physical_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"preload": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"virtual_directory": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"virtual_path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"physical_path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

type StorageAccount struct {
	Name        string `tfschema:"name"`
	Type        string `tfschema:"type"`
	AccountName string `tfschema:"account_name"`
	ShareName   string `tfschema:"share_name"`
	AccessKey   string `tfschema:"access_key"`
	MountPath   string `tfschema:"mount_path"`
}

func StorageAccountSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAzureStorageType(), false),
				},

				"account_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"share_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"access_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func StorageAccountSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(webapps.AzureStorageTypeAzureFiles),
					}, false),
				},

				"account_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"share_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"access_key": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func StorageAccountSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"account_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"share_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"access_key": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},

				"mount_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

type Backup struct {
	Name              string           `tfschema:"name"`
	StorageAccountURL string           `tfschema:"storage_account_url"`
	Enabled           bool             `tfschema:"enabled"`
	Schedule          []BackupSchedule `tfschema:"schedule"`
}

type BackupSchedule struct {
	FrequencyInterval    int64  `tfschema:"frequency_interval"`
	FrequencyUnit        string `tfschema:"frequency_unit"`
	KeepAtLeastOneBackup bool   `tfschema:"keep_at_least_one_backup"`
	RetentionPeriodDays  int64  `tfschema:"retention_period_days"`
	StartTime            string `tfschema:"start_time"`
	LastExecutionTime    string `tfschema:"last_execution_time"`
}

func BackupSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name which should be used for this Backup.",
				},

				"storage_account_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.IsURLWithHTTPS,
					Description:  "The SAS URL to the container.",
				},

				"enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Should this backup job be enabled?",
				},

				"schedule": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"frequency_interval": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntBetween(0, 1000),
								Description:  "How often the backup should be executed (e.g. for weekly backup, this should be set to `7` and `frequency_unit` should be set to `Day`).",
							},

							"frequency_unit": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									"Day",
									"Hour",
								}, false),
								Description: "The unit of time for how often the backup should take place. Possible values include: `Day` and `Hour`.",
							},

							"keep_at_least_one_backup": {
								Type:        pluginsdk.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Should the service keep at least one backup, regardless of age of backup. Defaults to `false`.",
							},

							"retention_period_days": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								Default:      30,
								ValidateFunc: validation.IntBetween(0, 9999999),
								Description:  "After how many days backups should be deleted.",
							},

							"start_time": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
								Description:  "When the schedule should start working in RFC-3339 format.",
								ValidateFunc: validation.IsRFC3339Time,
							},

							"last_execution_time": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The time the backup was last attempted.",
							},
						},
					},
				},
			},
		},
	}
}

func BackupSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of this Backup.",
				},

				"storage_account_url": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The SAS URL to the container.",
				},

				"enabled": {
					Type:        pluginsdk.TypeBool,
					Computed:    true,
					Description: "Is this backup job enabled?",
				},

				"schedule": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"frequency_interval": {
								Type:        pluginsdk.TypeInt,
								Computed:    true,
								Description: "How often the backup should is executed in multiples of the `frequency_unit`.",
							},

							"frequency_unit": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The unit of time for how often the backup takes place.",
							},

							"keep_at_least_one_backup": {
								Type:        pluginsdk.TypeBool,
								Computed:    true,
								Description: "Does the service keep at least one backup, regardless of age of backup.",
							},

							"retention_period_days": {
								Type:        pluginsdk.TypeInt,
								Computed:    true,
								Description: "After how many days are backups deleted.",
							},

							"start_time": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "When the schedule should start working in RFC-3339 format.",
							},

							"last_execution_time": {
								Type:        pluginsdk.TypeString,
								Computed:    true,
								Description: "The time the backup was last attempted.",
							},
						},
					},
				},
			},
		},
	}
}

type ConnectionString struct {
	Name  string `tfschema:"name"`
	Type  string `tfschema:"type"`
	Value string `tfschema:"value"`
}

func ConnectionStringSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The name which should be used for this Connection.",
				},

				"type": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForConnectionStringType(), false),
					Description:  "Type of database. Possible values include: `MySQL`, `SQLServer`, `SQLAzure`, `Custom`, `NotificationHub`, `ServiceBus`, `EventHub`, `APIHub`, `DocDb`, `RedisCache`, and `PostgreSQL`.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The connection string value.",
				},
			},
		},
	}
}

func ConnectionStringSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The name of this Connection.",
				},

				"type": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Description: "The type of database.",
				},

				"value": {
					Type:        pluginsdk.TypeString,
					Computed:    true,
					Sensitive:   true,
					Description: "The connection string value.",
				},
			},
		},
	}
}

type LogsConfig struct {
	ApplicationLogs       []ApplicationLog `tfschema:"application_logs"`
	HttpLogs              []HttpLog        `tfschema:"http_logs"`
	DetailedErrorMessages bool             `tfschema:"detailed_error_messages"`
	FailedRequestTracing  bool             `tfschema:"failed_request_tracing"`
}

type ApplicationLog struct {
	FileSystemLevel  string             `tfschema:"file_system_level"`
	AzureBlobStorage []AzureBlobStorage `tfschema:"azure_blob_storage"`
}

type AzureBlobStorage struct {
	Level           string `tfschema:"level"`
	SasURL          string `tfschema:"sas_url"`
	RetentionInDays int64  `tfschema:"retention_in_days"`
}

type HttpLog struct {
	FileSystems      []LogsFileSystem       `tfschema:"file_system"`
	AzureBlobStorage []AzureBlobStorageHttp `tfschema:"azure_blob_storage"`
}

type AzureBlobStorageHttp struct {
	SasURL          string `tfschema:"sas_url"`
	RetentionInDays int64  `tfschema:"retention_in_days"`
}

type LogsFileSystem struct {
	RetentionMB   int64 `tfschema:"retention_in_mb"`
	RetentionDays int64 `tfschema:"retention_in_days"`
}

func LogsConfigSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"application_logs": applicationLogSchema(),

				"http_logs": httpLogSchema(),

				"failed_request_tracing": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"detailed_error_messages": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func LogsConfigSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"application_logs": applicationLogSchemaComputed(),

				"http_logs": httpLogSchemaComputed(),

				"failed_request_tracing": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"detailed_error_messages": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

func applicationLogSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system_level": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{ // webapps.LoglevelOff is the implied value when this block is removed.
						string(webapps.LogLevelError),
						string(webapps.LogLevelOff),
						string(webapps.LogLevelInformation),
						string(webapps.LogLevelVerbose),
						string(webapps.LogLevelWarning),
					}, false),
				},

				"azure_blob_storage": appLogBlobStorageSchema(),
			},
		},
		DiffSuppressFunc: func(k, _, _ string, d *schema.ResourceData) bool {
			stateLogs, planLogs := d.GetChange("logs.0.application_logs")
			if stateLogs == nil || planLogs == nil {
				return false
			}
			stateAttrs := stateLogs.([]interface{})
			planAttrs := planLogs.([]interface{})

			// If the plan wants to set default values and the state is empty; suppress diff
			if len(stateAttrs) == 0 && len(planAttrs) > 0 && planAttrs[0] != nil {
				planAttr := planAttrs[0].(map[string]interface{})
				newFileSystemLevel, ok := planAttr["file_system_level"].(string)
				if !ok {
					return false
				}

				// if something is in `azure_blob_storage`, then we don't suppress the diff as we don't allow the default values for `azure_blob_storage` to be passed in
				newAzureBlobStorage, ok := planAttr["azure_blob_storage"].([]interface{})
				if !ok || len(newAzureBlobStorage) != 0 {
					return false
				}

				if newFileSystemLevel == string(webapps.LogLevelOff) {
					return true
				}
			}

			return false
		},
	}
}

func applicationLogSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system_level": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"azure_blob_storage": appLogBlobStorageSchemaComputed(),
			},
		},
	}
}

func appLogBlobStorageSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"level": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{ // webapps.LoglevelOff is the implied value when this block is removed.
						string(webapps.LogLevelError),
						string(webapps.LogLevelInformation),
						string(webapps.LogLevelVerbose),
						string(webapps.LogLevelWarning),
					}, false),
				},
				"sas_url": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Required: true,
					// TODO: Validation here?
				},
			},
		},
	}
}

func appLogBlobStorageSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"level": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func httpLogSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system": httpLogFileSystemSchema(),

				"azure_blob_storage": httpLogBlobStorageSchema(),
			},
		},
	}
}

func httpLogSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"file_system": httpLogFileSystemSchemaComputed(),

				"azure_blob_storage": httpLogBlobStorageSchemaComputed(),
			},
		},
	}
}

func httpLogFileSystemSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"logs.0.http_logs.0.azure_blob_storage"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"retention_in_mb": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(25, 100),
				},

				"retention_in_days": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(0),
				},
			},
		},
	}
}

func httpLogFileSystemSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"retention_in_mb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func httpLogBlobStorageSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:          pluginsdk.TypeList,
		Optional:      true,
		MaxItems:      1,
		ConflictsWith: []string{"logs.0.http_logs.0.file_system"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Required:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0), // Variable validation here based on the Service Plan SKU
				},
			},
		},
	}
}

func httpLogBlobStorageSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"sas_url": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
				"retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func ExpandLogsConfig(config []LogsConfig) *webapps.SiteLogsConfig {
	result := &webapps.SiteLogsConfig{}
	if len(config) == 0 {
		return result
	}

	result.Properties = &webapps.SiteLogsConfigProperties{}

	logsConfig := config[0]

	if len(logsConfig.ApplicationLogs) == 1 {
		appLogs := logsConfig.ApplicationLogs[0]
		result.Properties.ApplicationLogs = &webapps.ApplicationLogsConfig{
			FileSystem: &webapps.FileSystemApplicationLogsConfig{
				Level: pointer.To(webapps.LogLevel(appLogs.FileSystemLevel)),
			},
		}
		if len(appLogs.AzureBlobStorage) == 1 {
			appLogsBlobs := appLogs.AzureBlobStorage[0]
			result.Properties.ApplicationLogs.AzureBlobStorage = &webapps.AzureBlobStorageApplicationLogsConfig{
				Level:           pointer.To(webapps.LogLevel(appLogsBlobs.Level)),
				SasURL:          pointer.To(appLogsBlobs.SasURL),
				RetentionInDays: pointer.To(appLogsBlobs.RetentionInDays),
			}
		}
	}

	if len(logsConfig.HttpLogs) == 1 {
		httpLogs := logsConfig.HttpLogs[0]
		result.Properties.HTTPLogs = &webapps.HTTPLogsConfig{}

		if len(httpLogs.FileSystems) == 1 {
			httpLogFileSystem := httpLogs.FileSystems[0]
			result.Properties.HTTPLogs.FileSystem = &webapps.FileSystemHTTPLogsConfig{
				Enabled:         pointer.To(true),
				RetentionInMb:   pointer.To(httpLogFileSystem.RetentionMB),
				RetentionInDays: pointer.To(httpLogFileSystem.RetentionDays),
			}
		}

		if len(httpLogs.AzureBlobStorage) == 1 {
			httpLogsBlobStorage := httpLogs.AzureBlobStorage[0]
			result.Properties.HTTPLogs.AzureBlobStorage = &webapps.AzureBlobStorageHTTPLogsConfig{
				Enabled:         pointer.To(httpLogsBlobStorage.SasURL != ""),
				SasURL:          pointer.To(httpLogsBlobStorage.SasURL),
				RetentionInDays: pointer.To(httpLogsBlobStorage.RetentionInDays),
			}
		}
	}

	result.Properties.DetailedErrorMessages = &webapps.EnabledConfig{
		Enabled: pointer.To(logsConfig.DetailedErrorMessages),
	}

	result.Properties.FailedRequestsTracing = &webapps.EnabledConfig{
		Enabled: pointer.To(logsConfig.FailedRequestTracing),
	}

	return result
}

func ExpandBackupConfig(backupConfigs []Backup) (*webapps.BackupRequest, error) {
	result := &webapps.BackupRequest{}
	if len(backupConfigs) == 0 {
		return result, nil
	}

	backupConfig := backupConfigs[0]
	backupSchedule := backupConfig.Schedule[0]
	result.Properties = &webapps.BackupRequestProperties{
		Enabled:           pointer.To(backupConfig.Enabled),
		BackupName:        pointer.To(backupConfig.Name),
		StorageAccountURL: backupConfig.StorageAccountURL,
		BackupSchedule: &webapps.BackupSchedule{
			FrequencyInterval:     backupSchedule.FrequencyInterval,
			FrequencyUnit:         webapps.FrequencyUnit(backupSchedule.FrequencyUnit),
			KeepAtLeastOneBackup:  backupSchedule.KeepAtLeastOneBackup,
			RetentionPeriodInDays: backupSchedule.RetentionPeriodDays,
		},
	}

	if backupSchedule.StartTime != "" {
		dateTimeToStart, err := time.Parse(time.RFC3339, backupSchedule.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parsing back up start_time: %+v", err)
		}
		result.Properties.BackupSchedule.StartTime = pointer.To(dateTimeToStart.Format("2006-01-02T15:04:05.999999"))
	}

	return result, nil
}

func ExpandStorageConfig(storageConfigs []StorageAccount) *webapps.AzureStoragePropertyDictionaryResource {
	storageAccounts := make(map[string]webapps.AzureStorageInfoValue)
	result := &webapps.AzureStoragePropertyDictionaryResource{}
	if len(storageConfigs) == 0 {
		result.Properties = &storageAccounts
		return result
	}

	for _, v := range storageConfigs {
		storageAccounts[v.Name] = webapps.AzureStorageInfoValue{
			Type:        pointer.To(webapps.AzureStorageType(v.Type)),
			AccountName: pointer.To(v.AccountName),
			ShareName:   pointer.To(v.ShareName),
			AccessKey:   pointer.To(v.AccessKey),
			MountPath:   pointer.To(v.MountPath),
		}
	}

	result.Properties = &storageAccounts

	return result
}

func ExpandConnectionStrings(connectionStringsConfig []ConnectionString) *webapps.ConnectionStringDictionary {
	result := &webapps.ConnectionStringDictionary{}
	if len(connectionStringsConfig) == 0 {
		return result
	}

	connectionStrings := make(map[string]webapps.ConnStringValueTypePair)
	for _, v := range connectionStringsConfig {
		connectionStrings[v.Name] = webapps.ConnStringValueTypePair{
			Value: v.Value,
			Type:  webapps.ConnectionStringType(v.Type),
		}
	}
	result.Properties = &connectionStrings

	return result
}

func expandHandlerMapping(handlerMapping []HandlerMappings) *[]webapps.HandlerMapping {
	if len(handlerMapping) == 0 {
		return nil
	}

	result := make([]webapps.HandlerMapping, 0)

	for _, v := range handlerMapping {
		if v.Arguments != "" {
			result = append(result, webapps.HandlerMapping{
				Extension:       pointer.To(v.Extension),
				ScriptProcessor: pointer.To(v.ScriptProcessorPath),
				Arguments:       pointer.To(v.Arguments),
			})
		} else {
			result = append(result, webapps.HandlerMapping{
				Extension:       pointer.To(v.Extension),
				ScriptProcessor: pointer.To(v.ScriptProcessorPath),
			})
		}
	}
	return &result
}

func expandHandlerMappingForUpdate(handlerMapping []HandlerMappings) *[]webapps.HandlerMapping {
	result := make([]webapps.HandlerMapping, 0)
	if len(handlerMapping) == 0 {
		return &result
	}

	for _, v := range handlerMapping {
		if v.Arguments != "" {
			result = append(result, webapps.HandlerMapping{
				Extension:       pointer.To(v.Extension),
				ScriptProcessor: pointer.To(v.ScriptProcessorPath),
				Arguments:       pointer.To(v.Arguments),
			})
		} else {
			result = append(result, webapps.HandlerMapping{
				Extension:       pointer.To(v.Extension),
				ScriptProcessor: pointer.To(v.ScriptProcessorPath),
			})
		}
	}
	return &result
}

func expandVirtualApplications(virtualApplicationConfig []VirtualApplication) *[]webapps.VirtualApplication {
	if len(virtualApplicationConfig) == 0 {
		return nil
	}

	result := make([]webapps.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := webapps.VirtualApplication{
			VirtualPath:    pointer.To(v.VirtualPath),
			PhysicalPath:   pointer.To(v.PhysicalPath),
			PreloadEnabled: pointer.To(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]webapps.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, webapps.VirtualDirectory{
					VirtualPath:  pointer.To(d.VirtualPath),
					PhysicalPath: pointer.To(d.PhysicalPath),
				})
			}
			virtualApp.VirtualDirectories = &virtualDirs
		}

		result = append(result, virtualApp)
	}
	return &result
}

func expandVirtualApplicationsForUpdate(virtualApplicationConfig []VirtualApplication) *[]webapps.VirtualApplication {
	if len(virtualApplicationConfig) == 0 {
		// to remove this block from the config we need to give the service the original default back, sending an empty struct leaves the previous config in place
		return &[]webapps.VirtualApplication{
			{
				VirtualPath:    pointer.To("/"),
				PhysicalPath:   pointer.To("site\\wwwroot"),
				PreloadEnabled: pointer.To(true),
			},
		}
	}

	result := make([]webapps.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := webapps.VirtualApplication{
			VirtualPath:    pointer.To(v.VirtualPath),
			PhysicalPath:   pointer.To(v.PhysicalPath),
			PreloadEnabled: pointer.To(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]webapps.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, webapps.VirtualDirectory{
					VirtualPath:  pointer.To(d.VirtualPath),
					PhysicalPath: pointer.To(d.PhysicalPath),
				})
			}
			virtualApp.VirtualDirectories = &virtualDirs
		}

		result = append(result, virtualApp)
	}
	return &result
}

func FlattenBackupConfig(backupRequest *webapps.BackupRequest) []Backup {
	if backupRequest == nil || backupRequest.Properties == nil {
		return []Backup{}
	}
	props := *backupRequest.Properties
	backup := Backup{
		StorageAccountURL: props.StorageAccountURL,
	}
	if props.BackupName != nil {
		backup.Name = *props.BackupName
	}

	if props.Enabled != nil {
		backup.Enabled = *props.Enabled
	}

	if schedule := props.BackupSchedule; schedule != nil {
		backupSchedule := BackupSchedule{
			FrequencyUnit:        string(schedule.FrequencyUnit),
			FrequencyInterval:    schedule.FrequencyInterval,
			KeepAtLeastOneBackup: schedule.KeepAtLeastOneBackup,
			RetentionPeriodDays:  schedule.RetentionPeriodInDays,
		}

		startTimeAsTime, err := time.Parse("2006-01-02T15:04:05.999999", *schedule.StartTime)
		if err == nil {
			if schedule.StartTime != nil && !startTimeAsTime.IsZero() {
				backupSchedule.StartTime = startTimeAsTime.Format(time.RFC3339)
			}
		}

		if schedule.LastExecutionTime != nil {
			lastExecutionTimeAsTime, err := time.Parse("2006-01-02T15:04:05.999999", *schedule.LastExecutionTime)
			if err == nil {
				if schedule.LastExecutionTime != nil && !lastExecutionTimeAsTime.IsZero() {
					backupSchedule.LastExecutionTime = lastExecutionTimeAsTime.Format(time.RFC3339)
				}
			}
		}

		backup.Schedule = []BackupSchedule{backupSchedule}
	}

	return []Backup{backup}
}

func FlattenLogsConfig(logsConfig *webapps.SiteLogsConfig) []LogsConfig {
	if logsConfig == nil || logsConfig.Properties == nil {
		return []LogsConfig{}
	}
	props := *logsConfig.Properties
	if onlyDefaultLoggingConfig(props) {
		return nil
	}

	logs := LogsConfig{}

	if props.ApplicationLogs != nil {
		appLogs := *props.ApplicationLogs
		applicationLog := ApplicationLog{}

		if appLogs.FileSystem != nil {
			applicationLog.FileSystemLevel = string(pointer.From(appLogs.FileSystem.Level))
			if appLogs.AzureBlobStorage != nil && appLogs.AzureBlobStorage.SasURL != nil {
				blobStorage := AzureBlobStorage{
					Level: string(pointer.From(appLogs.AzureBlobStorage.Level)),
				}

				blobStorage.SasURL = pointer.From(appLogs.AzureBlobStorage.SasURL)

				blobStorage.RetentionInDays = pointer.From(appLogs.AzureBlobStorage.RetentionInDays)

				applicationLog.AzureBlobStorage = []AzureBlobStorage{blobStorage}
			}

			// Only set ApplicationLogs if it's not the default values
			/*
				"applicationLogs": {
					"fileSystem": {
						"level": "Off"
					},
					"azureBlobStorage": {
						"level": "Off",
						"sasUrl": null,
						"retentionInDays": null
					}
				},
			*/
			if !strings.EqualFold(string(pointer.From(appLogs.FileSystem.Level)), string(webapps.LogLevelOff)) && len(applicationLog.AzureBlobStorage) > 0 {
				logs.ApplicationLogs = []ApplicationLog{applicationLog}
			}
		}
	}

	if props.HTTPLogs != nil {
		httpLogs := *props.HTTPLogs
		httpLog := HttpLog{}

		if httpLogs.FileSystem != nil && (httpLogs.FileSystem.Enabled != nil && *httpLogs.FileSystem.Enabled) {
			fileSystem := LogsFileSystem{}
			if httpLogs.FileSystem.RetentionInMb != nil {
				fileSystem.RetentionMB = pointer.From(httpLogs.FileSystem.RetentionInMb)
			}

			if httpLogs.FileSystem.RetentionInDays != nil {
				fileSystem.RetentionDays = pointer.From(httpLogs.FileSystem.RetentionInDays)
			}

			httpLog.FileSystems = []LogsFileSystem{fileSystem}
		}

		if httpLogs.AzureBlobStorage != nil && (httpLogs.AzureBlobStorage.Enabled != nil && *httpLogs.AzureBlobStorage.Enabled) {
			blobStorage := AzureBlobStorageHttp{}
			if httpLogs.AzureBlobStorage.SasURL != nil {
				blobStorage.SasURL = *httpLogs.AzureBlobStorage.SasURL
			}

			if httpLogs.AzureBlobStorage.RetentionInDays != nil {
				blobStorage.RetentionInDays = pointer.From(httpLogs.AzureBlobStorage.RetentionInDays)
			}

			if blobStorage.RetentionInDays != 0 || blobStorage.SasURL != "" {
				httpLog.AzureBlobStorage = []AzureBlobStorageHttp{blobStorage}
			}
		}

		if httpLog.FileSystems != nil || httpLog.AzureBlobStorage != nil {
			logs.HttpLogs = []HttpLog{httpLog}
		}
	}

	// logs.DetailedErrorMessages = false
	if props.DetailedErrorMessages != nil && props.DetailedErrorMessages.Enabled != nil {
		logs.DetailedErrorMessages = *props.DetailedErrorMessages.Enabled
	}

	// logs.FailedRequestTracing = false
	if props.FailedRequestsTracing != nil && props.FailedRequestsTracing.Enabled != nil {
		logs.FailedRequestTracing = *props.FailedRequestsTracing.Enabled
	}

	return []LogsConfig{logs}
}

func onlyDefaultLoggingConfig(props webapps.SiteLogsConfigProperties) bool {
	if props.ApplicationLogs == nil || props.HTTPLogs == nil || props.FailedRequestsTracing == nil || props.DetailedErrorMessages == nil {
		return false
	}
	if props.ApplicationLogs.FileSystem != nil && pointer.From(props.ApplicationLogs.FileSystem.Level) != webapps.LogLevelOff {
		return false
	}
	if props.ApplicationLogs.AzureBlobStorage != nil && pointer.From(props.ApplicationLogs.AzureBlobStorage.Level) != webapps.LogLevelOff {
		return false
	}
	if props.HTTPLogs.FileSystem != nil && props.HTTPLogs.FileSystem.Enabled != nil && (*props.HTTPLogs.FileSystem.Enabled) {
		return false
	}
	if props.HTTPLogs.AzureBlobStorage != nil && props.HTTPLogs.AzureBlobStorage.Enabled != nil && (*props.HTTPLogs.AzureBlobStorage.Enabled) {
		return false
	}
	if props.FailedRequestsTracing.Enabled == nil || *props.FailedRequestsTracing.Enabled {
		return false
	}
	if props.DetailedErrorMessages.Enabled == nil || *props.DetailedErrorMessages.Enabled {
		return false
	}
	return true
}

func FlattenStorageAccounts(appStorageAccounts *webapps.AzureStoragePropertyDictionaryResource) []StorageAccount {
	if appStorageAccounts == nil || len(*appStorageAccounts.Properties) == 0 {
		return []StorageAccount{}
	}

	storageAccounts := make([]StorageAccount, 0, len(*appStorageAccounts.Properties))
	for k, v := range *appStorageAccounts.Properties {
		storageAccount := StorageAccount{
			Name: k,
			Type: string(pointer.From(v.Type)),
		}
		if v.AccountName != nil {
			storageAccount.AccountName = *v.AccountName
		}

		if v.ShareName != nil {
			storageAccount.ShareName = *v.ShareName
		}

		if v.AccessKey != nil {
			storageAccount.AccessKey = *v.AccessKey
		}

		if v.MountPath != nil {
			storageAccount.MountPath = *v.MountPath
		}

		storageAccounts = append(storageAccounts, storageAccount)
	}

	return storageAccounts
}

func FlattenConnectionStrings(appConnectionStrings *webapps.ConnectionStringDictionary) []ConnectionString {
	if appConnectionStrings.Properties == nil || len(*appConnectionStrings.Properties) == 0 {
		return []ConnectionString{}
	}

	connectionStrings := make([]ConnectionString, 0, len(*appConnectionStrings.Properties))
	for k, v := range *appConnectionStrings.Properties {
		connectionString := ConnectionString{
			Name:  k,
			Type:  string(v.Type),
			Value: v.Value,
		}
		connectionStrings = append(connectionStrings, connectionString)
	}

	return connectionStrings
}

func ExpandAppSettingsForUpdate(siteConfigSettings *[]webapps.NameValuePair) *webapps.StringDictionary {
	appSettings := make(map[string]string)
	if siteConfigSettings == nil {
		return &webapps.StringDictionary{
			Properties: &appSettings,
		}
	}

	for _, v := range *siteConfigSettings {
		if name := pointer.From(v.Name); name != "" {
			appSettings[name] = pointer.From(v.Value)
		}
	}

	return &webapps.StringDictionary{
		Properties: &appSettings,
	}
}

func ExpandAppSettingsForCreate(settings map[string]string) *[]webapps.NameValuePair {
	if len(settings) > 0 {
		result := make([]webapps.NameValuePair, 0)
		for k, v := range settings {
			result = append(result, webapps.NameValuePair{
				Name:  pointer.To(k),
				Value: pointer.To(v),
			})
		}
		return &result
	}
	return nil
}

// FilterManagedAppSettings removes app_settings values from the state that are controlled directly be schema properties.
func FilterManagedAppSettings(input map[string]string) map[string]string {
	unmanagedSettings := []string{
		"DOCKER_REGISTRY_SERVER_URL",
		"DOCKER_REGISTRY_SERVER_USERNAME",
		"DOCKER_REGISTRY_SERVER_PASSWORD",
		"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
		"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
		"WEBSITE_HTTPLOGGING_CONTAINER_URL",
		"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
		"WEBSITE_VNET_ROUTE_ALL",
		"spring.datasource.password",
		"spring.datasource.url",
		"spring.datasource.username",
		"WEBSITE_HEALTHCHECK_MAXPINGFAILURES",
	}

	for _, v := range unmanagedSettings { //nolint:typecheck
		delete(input, v)
	}

	return input
}

// FilterManagedAppSettingsDeprecated removes app_settings values from the state that are controlled directly be
// schema properties when the deprecated docker settings are used. This function should be removed in 4.0
func FilterManagedAppSettingsDeprecated(input map[string]string) map[string]string {
	unmanagedSettings := []string{
		"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
		"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
		"WEBSITE_HTTPLOGGING_CONTAINER_URL",
		"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
		"WEBSITE_VNET_ROUTE_ALL",
		"spring.datasource.password",
		"spring.datasource.url",
		"spring.datasource.username",
		"WEBSITE_HEALTHCHECK_MAXPINGFAILURES",
	}

	for _, v := range unmanagedSettings { //nolint:typecheck
		delete(input, v)
	}

	return input
}

func flattenHandlerMapping(appHandlerMappings *[]webapps.HandlerMapping) []HandlerMappings {
	if appHandlerMappings == nil {
		return []HandlerMappings{}
	}

	handlerMappings := make([]HandlerMappings, 0, len(*appHandlerMappings))
	for _, v := range *appHandlerMappings {
		handlerMapping := HandlerMappings{
			Extension:           pointer.From(v.Extension),
			ScriptProcessorPath: pointer.From(v.ScriptProcessor),
		}
		handlerMapping.Arguments = pointer.From(v.Arguments)
		handlerMappings = append(handlerMappings, handlerMapping)
	}

	return handlerMappings
}

func flattenVirtualApplications(appVirtualApplications *[]webapps.VirtualApplication, alwaysOn bool) []VirtualApplication {
	if appVirtualApplications == nil || onlyDefaultVirtualApplication(*appVirtualApplications, alwaysOn) {
		return []VirtualApplication{}
	}

	virtualApplications := make([]VirtualApplication, 0, len(*appVirtualApplications))
	for _, v := range *appVirtualApplications {
		virtualApp := VirtualApplication{
			VirtualPath:  pointer.From(v.VirtualPath),
			PhysicalPath: pointer.From(v.PhysicalPath),
		}
		if preload := v.PreloadEnabled; preload != nil {
			virtualApp.Preload = *preload
		}
		if v.VirtualDirectories != nil && len(*v.VirtualDirectories) > 0 {
			virtualDirs := make([]VirtualDirectory, 0)
			for _, d := range *v.VirtualDirectories {
				virtualDir := VirtualDirectory{
					VirtualPath:  pointer.From(d.VirtualPath),
					PhysicalPath: pointer.From(d.PhysicalPath),
				}
				virtualDirs = append(virtualDirs, virtualDir)
			}
			virtualApp.VirtualDirectories = virtualDirs
		}
		virtualApplications = append(virtualApplications, virtualApp)
	}

	return virtualApplications
}

func onlyDefaultVirtualApplication(input []webapps.VirtualApplication, alwaysOn bool) bool {
	if len(input) > 1 {
		return false
	}
	app := input[0]
	if app.VirtualPath == nil || app.PhysicalPath == nil {
		return false
	}

	if *app.VirtualPath == "/" && *app.PhysicalPath == "site\\wwwroot" && app.VirtualDirectories == nil {
		// if alwaysOn is true, then the default for PreloadEnabled is true
		// if alwaysOn is false, then the default for PreloadEnabled is false
		if (alwaysOn && *app.PreloadEnabled) || (!alwaysOn && !*app.PreloadEnabled) {
			return true
		}
	}
	return false
}

func DisabledLogsConfig() *webapps.SiteLogsConfig {
	return &webapps.SiteLogsConfig{
		Properties: &webapps.SiteLogsConfigProperties{
			DetailedErrorMessages: &webapps.EnabledConfig{
				Enabled: pointer.To(false),
			},
			FailedRequestsTracing: &webapps.EnabledConfig{
				Enabled: pointer.To(false),
			},
			ApplicationLogs: &webapps.ApplicationLogsConfig{
				FileSystem: &webapps.FileSystemApplicationLogsConfig{
					Level: pointer.To(webapps.LogLevelOff),
				},
				AzureBlobStorage: &webapps.AzureBlobStorageApplicationLogsConfig{
					Level: pointer.To(webapps.LogLevelOff),
				},
			},
			HTTPLogs: &webapps.HTTPLogsConfig{
				FileSystem: &webapps.FileSystemHTTPLogsConfig{
					Enabled: pointer.To(false),
				},
				AzureBlobStorage: &webapps.AzureBlobStorageHTTPLogsConfig{
					Enabled: pointer.To(false),
				},
			},
		},
	}
}

func IsFreeOrSharedServicePlan(inputSKU string) bool {
	result := false
	for _, sku := range freeSkus {
		if inputSKU == sku {
			result = true
		}
	}
	for _, sku := range sharedSkus {
		if inputSKU == sku {
			result = true
		}
	}
	return result
}
