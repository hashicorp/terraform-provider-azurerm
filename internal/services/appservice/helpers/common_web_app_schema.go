// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/tombuildsstuff/kermit/sdk/web/2022-09-01/web"
)

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
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AzureStorageTypeAzureBlob),
						string(web.AzureStorageTypeAzureFiles),
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
						string(web.AzureStorageTypeAzureFiles),
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
	StorageAccountUrl string           `tfschema:"storage_account_url"`
	Enabled           bool             `tfschema:"enabled"`
	Schedule          []BackupSchedule `tfschema:"schedule"`
}

type BackupSchedule struct {
	FrequencyInterval    int    `tfschema:"frequency_interval"`
	FrequencyUnit        string `tfschema:"frequency_unit"`
	KeepAtLeastOneBackup bool   `tfschema:"keep_at_least_one_backup"`
	RetentionPeriodDays  int    `tfschema:"retention_period_days"`
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
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ConnectionStringTypeAPIHub),
						string(web.ConnectionStringTypeCustom),
						string(web.ConnectionStringTypeDocDb),
						string(web.ConnectionStringTypeEventHub),
						string(web.ConnectionStringTypeMySQL),
						string(web.ConnectionStringTypeNotificationHub),
						string(web.ConnectionStringTypePostgreSQL),
						string(web.ConnectionStringTypeRedisCache),
						string(web.ConnectionStringTypeServiceBus),
						string(web.ConnectionStringTypeSQLAzure),
						string(web.ConnectionStringTypeSQLServer),
					}, false),
					Description: "Type of database. Possible values include: `MySQL`, `SQLServer`, `SQLAzure`, `Custom`, `NotificationHub`, `ServiceBus`, `EventHub`, `APIHub`, `DocDb`, `RedisCache`, and `PostgreSQL`.",
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
	SasUrl          string `tfschema:"sas_url"`
	RetentionInDays int    `tfschema:"retention_in_days"`
}

type HttpLog struct {
	FileSystems      []LogsFileSystem       `tfschema:"file_system"`
	AzureBlobStorage []AzureBlobStorageHttp `tfschema:"azure_blob_storage"`
}

type AzureBlobStorageHttp struct {
	SasUrl          string `tfschema:"sas_url"`
	RetentionInDays int    `tfschema:"retention_in_days"`
}

type LogsFileSystem struct {
	RetentionMB   int `tfschema:"retention_in_mb"`
	RetentionDays int `tfschema:"retention_in_days"`
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
					ValidateFunc: validation.StringInSlice([]string{
						string(web.LogLevelError),
						string(web.LogLevelInformation),
						string(web.LogLevelVerbose),
						string(web.LogLevelWarning),
					}, false),
				},

				"azure_blob_storage": appLogBlobStorageSchema(),
			},
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
					ValidateFunc: validation.StringInSlice([]string{
						string(web.LogLevelError),
						string(web.LogLevelInformation),
						string(web.LogLevelVerbose),
						string(web.LogLevelWarning),
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

func ExpandLogsConfig(config []LogsConfig) *web.SiteLogsConfig {
	result := &web.SiteLogsConfig{}
	if len(config) == 0 {
		return result
	}

	result.SiteLogsConfigProperties = &web.SiteLogsConfigProperties{}

	logsConfig := config[0]

	if len(logsConfig.ApplicationLogs) == 1 {
		appLogs := logsConfig.ApplicationLogs[0]
		result.SiteLogsConfigProperties.ApplicationLogs = &web.ApplicationLogsConfig{
			FileSystem: &web.FileSystemApplicationLogsConfig{
				Level: web.LogLevel(appLogs.FileSystemLevel),
			},
		}
		if len(appLogs.AzureBlobStorage) == 1 {
			appLogsBlobs := appLogs.AzureBlobStorage[0]
			result.SiteLogsConfigProperties.ApplicationLogs.AzureBlobStorage = &web.AzureBlobStorageApplicationLogsConfig{
				Level:           web.LogLevel(appLogsBlobs.Level),
				SasURL:          pointer.To(appLogsBlobs.SasUrl),
				RetentionInDays: pointer.To(int32(appLogsBlobs.RetentionInDays)),
			}
		}
	}

	if len(logsConfig.HttpLogs) == 1 {
		httpLogs := logsConfig.HttpLogs[0]
		result.HTTPLogs = &web.HTTPLogsConfig{}

		if len(httpLogs.FileSystems) == 1 {
			httpLogFileSystem := httpLogs.FileSystems[0]
			result.HTTPLogs.FileSystem = &web.FileSystemHTTPLogsConfig{
				Enabled:         pointer.To(true),
				RetentionInMb:   pointer.To(int32(httpLogFileSystem.RetentionMB)),
				RetentionInDays: pointer.To(int32(httpLogFileSystem.RetentionDays)),
			}
		}

		if len(httpLogs.AzureBlobStorage) == 1 {
			httpLogsBlobStorage := httpLogs.AzureBlobStorage[0]
			result.HTTPLogs.AzureBlobStorage = &web.AzureBlobStorageHTTPLogsConfig{
				Enabled:         pointer.To(httpLogsBlobStorage.SasUrl != ""),
				SasURL:          pointer.To(httpLogsBlobStorage.SasUrl),
				RetentionInDays: pointer.To(int32(httpLogsBlobStorage.RetentionInDays)),
			}
		}
	}

	result.DetailedErrorMessages = &web.EnabledConfig{
		Enabled: pointer.To(logsConfig.DetailedErrorMessages),
	}

	result.FailedRequestsTracing = &web.EnabledConfig{
		Enabled: pointer.To(logsConfig.FailedRequestTracing),
	}

	return result
}

func ExpandBackupConfig(backupConfigs []Backup) (*web.BackupRequest, error) {
	result := &web.BackupRequest{}
	if len(backupConfigs) == 0 {
		return result, nil
	}

	backupConfig := backupConfigs[0]
	backupSchedule := backupConfig.Schedule[0]
	result.BackupRequestProperties = &web.BackupRequestProperties{
		Enabled:           pointer.To(backupConfig.Enabled),
		BackupName:        pointer.To(backupConfig.Name),
		StorageAccountURL: pointer.To(backupConfig.StorageAccountUrl),
		BackupSchedule: &web.BackupSchedule{
			FrequencyInterval:     pointer.To(int32(backupSchedule.FrequencyInterval)),
			FrequencyUnit:         web.FrequencyUnit(backupSchedule.FrequencyUnit),
			KeepAtLeastOneBackup:  pointer.To(backupSchedule.KeepAtLeastOneBackup),
			RetentionPeriodInDays: pointer.To(int32(backupSchedule.RetentionPeriodDays)),
		},
	}

	if backupSchedule.StartTime != "" {
		dateTimeToStart, err := time.Parse(time.RFC3339, backupSchedule.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parsing back up start_time: %+v", err)
		}
		result.BackupRequestProperties.BackupSchedule.StartTime = &date.Time{Time: dateTimeToStart}
	}

	return result, nil
}

func ExpandStorageConfig(storageConfigs []StorageAccount) *web.AzureStoragePropertyDictionaryResource {
	storageAccounts := make(map[string]*web.AzureStorageInfoValue)
	result := &web.AzureStoragePropertyDictionaryResource{}
	if len(storageConfigs) == 0 {
		result.Properties = storageAccounts
		return result
	}

	for _, v := range storageConfigs {
		storageAccounts[v.Name] = &web.AzureStorageInfoValue{
			Type:        web.AzureStorageType(v.Type),
			AccountName: pointer.To(v.AccountName),
			ShareName:   pointer.To(v.ShareName),
			AccessKey:   pointer.To(v.AccessKey),
			MountPath:   pointer.To(v.MountPath),
		}
	}

	result.Properties = storageAccounts

	return result
}

func ExpandConnectionStrings(connectionStringsConfig []ConnectionString) *web.ConnectionStringDictionary {
	result := &web.ConnectionStringDictionary{}
	if len(connectionStringsConfig) == 0 {
		return result
	}

	connectionStrings := make(map[string]*web.ConnStringValueTypePair)
	for _, v := range connectionStringsConfig {
		connectionStrings[v.Name] = &web.ConnStringValueTypePair{
			Value: pointer.To(v.Value),
			Type:  web.ConnectionStringType(v.Type),
		}
	}
	result.Properties = connectionStrings

	return result
}

func expandVirtualApplications(virtualApplicationConfig []VirtualApplication) *[]web.VirtualApplication {
	if len(virtualApplicationConfig) == 0 {
		return nil
	}

	result := make([]web.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := web.VirtualApplication{
			VirtualPath:    pointer.To(v.VirtualPath),
			PhysicalPath:   pointer.To(v.PhysicalPath),
			PreloadEnabled: pointer.To(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]web.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, web.VirtualDirectory{
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

func expandVirtualApplicationsForUpdate(virtualApplicationConfig []VirtualApplication) *[]web.VirtualApplication {
	if len(virtualApplicationConfig) == 0 {
		// to remove this block from the config we need to give the service the original default back, sending an empty struct leaves the previous config in place
		return &[]web.VirtualApplication{
			{
				VirtualPath:    pointer.To("/"),
				PhysicalPath:   pointer.To("site\\wwwroot"),
				PreloadEnabled: pointer.To(true),
			},
		}
	}

	result := make([]web.VirtualApplication, 0)

	for _, v := range virtualApplicationConfig {
		virtualApp := web.VirtualApplication{
			VirtualPath:    pointer.To(v.VirtualPath),
			PhysicalPath:   pointer.To(v.PhysicalPath),
			PreloadEnabled: pointer.To(v.Preload),
		}
		if len(v.VirtualDirectories) > 0 {
			virtualDirs := make([]web.VirtualDirectory, 0)
			for _, d := range v.VirtualDirectories {
				virtualDirs = append(virtualDirs, web.VirtualDirectory{
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

func FlattenBackupConfig(backupRequest web.BackupRequest) []Backup {
	if backupRequest.BackupRequestProperties == nil {
		return []Backup{}
	}
	props := *backupRequest.BackupRequestProperties
	backup := Backup{}
	if props.BackupName != nil {
		backup.Name = *props.BackupName
	}

	if props.StorageAccountURL != nil {
		backup.StorageAccountUrl = *props.StorageAccountURL
	}

	if props.Enabled != nil {
		backup.Enabled = *props.Enabled
	}

	if schedule := props.BackupSchedule; schedule != nil {
		backupSchedule := BackupSchedule{
			FrequencyUnit: string(schedule.FrequencyUnit),
		}
		if schedule.FrequencyInterval != nil {
			backupSchedule.FrequencyInterval = int(*schedule.FrequencyInterval)
		}

		if schedule.KeepAtLeastOneBackup != nil {
			backupSchedule.KeepAtLeastOneBackup = *schedule.KeepAtLeastOneBackup
		}

		if schedule.RetentionPeriodInDays != nil {
			backupSchedule.RetentionPeriodDays = int(*schedule.RetentionPeriodInDays)
		}

		if schedule.StartTime != nil && !schedule.StartTime.IsZero() {
			backupSchedule.StartTime = schedule.StartTime.Format(time.RFC3339)
		}

		if schedule.LastExecutionTime != nil && !schedule.LastExecutionTime.IsZero() {
			backupSchedule.LastExecutionTime = schedule.LastExecutionTime.Format(time.RFC3339)
		}

		backup.Schedule = []BackupSchedule{backupSchedule}
	}

	return []Backup{backup}
}

func FlattenLogsConfig(logsConfig web.SiteLogsConfig) []LogsConfig {
	if logsConfig.SiteLogsConfigProperties == nil {
		return []LogsConfig{}
	}
	props := *logsConfig.SiteLogsConfigProperties
	if onlyDefaultLoggingConfig(props) {
		return nil
	}

	logs := LogsConfig{}

	if props.ApplicationLogs != nil {
		appLogs := *props.ApplicationLogs
		applicationLog := ApplicationLog{}

		if appLogs.FileSystem != nil && appLogs.FileSystem.Level != web.LogLevelOff {
			applicationLog.FileSystemLevel = string(appLogs.FileSystem.Level)
			if appLogs.AzureBlobStorage != nil && appLogs.AzureBlobStorage.Level != web.LogLevelOff {
				blobStorage := AzureBlobStorage{
					Level: string(appLogs.AzureBlobStorage.Level),
				}

				blobStorage.SasUrl = pointer.From(appLogs.AzureBlobStorage.SasURL)

				blobStorage.RetentionInDays = int(pointer.From(appLogs.AzureBlobStorage.RetentionInDays))

				applicationLog.AzureBlobStorage = []AzureBlobStorage{blobStorage}
			}
			logs.ApplicationLogs = []ApplicationLog{applicationLog}
		}
	}

	if props.HTTPLogs != nil {
		httpLogs := *props.HTTPLogs
		httpLog := HttpLog{}

		if httpLogs.FileSystem != nil && (httpLogs.FileSystem.Enabled != nil && *httpLogs.FileSystem.Enabled) {
			fileSystem := LogsFileSystem{}
			if httpLogs.FileSystem.RetentionInMb != nil {
				fileSystem.RetentionMB = int(*httpLogs.FileSystem.RetentionInMb)
			}

			if httpLogs.FileSystem.RetentionInDays != nil {
				fileSystem.RetentionDays = int(*httpLogs.FileSystem.RetentionInDays)
			}

			httpLog.FileSystems = []LogsFileSystem{fileSystem}
		}

		if httpLogs.AzureBlobStorage != nil && (httpLogs.AzureBlobStorage.Enabled != nil && *httpLogs.AzureBlobStorage.Enabled) {
			blobStorage := AzureBlobStorageHttp{}
			if httpLogs.AzureBlobStorage.SasURL != nil {
				blobStorage.SasUrl = *httpLogs.AzureBlobStorage.SasURL
			}

			if httpLogs.AzureBlobStorage.RetentionInDays != nil {
				blobStorage.RetentionInDays = int(*httpLogs.AzureBlobStorage.RetentionInDays)
			}

			if blobStorage.RetentionInDays != 0 || blobStorage.SasUrl != "" {
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

func onlyDefaultLoggingConfig(props web.SiteLogsConfigProperties) bool {
	if props.ApplicationLogs == nil || props.HTTPLogs == nil || props.FailedRequestsTracing == nil || props.DetailedErrorMessages == nil {
		return false
	}
	if props.ApplicationLogs.FileSystem != nil && props.ApplicationLogs.FileSystem.Level != web.LogLevelOff {
		return false
	}
	if props.ApplicationLogs.AzureBlobStorage != nil && props.ApplicationLogs.AzureBlobStorage.Level != web.LogLevelOff {
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

func FlattenStorageAccounts(appStorageAccounts web.AzureStoragePropertyDictionaryResource) []StorageAccount {
	if len(appStorageAccounts.Properties) == 0 {
		return []StorageAccount{}
	}
	var storageAccounts []StorageAccount
	for k, v := range appStorageAccounts.Properties {
		storageAccount := StorageAccount{
			Name: k,
			Type: string(v.Type),
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

func FlattenConnectionStrings(appConnectionStrings web.ConnectionStringDictionary) []ConnectionString {
	if len(appConnectionStrings.Properties) == 0 {
		return []ConnectionString{}
	}
	var connectionStrings []ConnectionString
	for k, v := range appConnectionStrings.Properties {
		connectionString := ConnectionString{
			Name: k,
			Type: string(v.Type),
		}
		if v.Value != nil {
			connectionString.Value = *v.Value
		}
		connectionStrings = append(connectionStrings, connectionString)
	}

	return connectionStrings
}

func ExpandAppSettingsForUpdate(settings map[string]string) *web.StringDictionary {
	appSettings := make(map[string]*string)
	for k, v := range settings {
		appSettings[k] = pointer.To(v)
	}

	return &web.StringDictionary{
		Properties: appSettings,
	}
}

func ExpandAppSettingsForCreate(settings map[string]string) *[]web.NameValuePair {
	if len(settings) > 0 {
		result := make([]web.NameValuePair, 0)
		for k, v := range settings {
			result = append(result, web.NameValuePair{
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

func flattenVirtualApplications(appVirtualApplications *[]web.VirtualApplication) []VirtualApplication {
	if appVirtualApplications == nil || onlyDefaultVirtualApplication(*appVirtualApplications) {
		return []VirtualApplication{}
	}

	var virtualApplications []VirtualApplication
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

func onlyDefaultVirtualApplication(input []web.VirtualApplication) bool {
	if len(input) > 1 {
		return false
	}
	app := input[0]
	if app.VirtualPath == nil || app.PhysicalPath == nil {
		return false
	}
	if *app.VirtualPath == "/" && *app.PhysicalPath == "site\\wwwroot" && *app.PreloadEnabled && app.VirtualDirectories == nil {
		return true
	}
	return false
}

func DisabledLogsConfig() *web.SiteLogsConfig {
	return &web.SiteLogsConfig{
		SiteLogsConfigProperties: &web.SiteLogsConfigProperties{
			DetailedErrorMessages: &web.EnabledConfig{
				Enabled: pointer.To(false),
			},
			FailedRequestsTracing: &web.EnabledConfig{
				Enabled: pointer.To(false),
			},
			ApplicationLogs: &web.ApplicationLogsConfig{
				FileSystem: &web.FileSystemApplicationLogsConfig{
					Level: web.LogLevelOff,
				},
				AzureBlobStorage: &web.AzureBlobStorageApplicationLogsConfig{
					Level: web.LogLevelOff,
				},
			},
			HTTPLogs: &web.HTTPLogsConfig{
				FileSystem: &web.FileSystemHTTPLogsConfig{
					Enabled: pointer.To(false),
				},
				AzureBlobStorage: &web.AzureBlobStorageHTTPLogsConfig{
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
