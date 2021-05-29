package web

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func schemaAppServiceBackup() *pluginsdk.Schema {
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
				},

				"storage_account_url": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					Sensitive:    true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},

				"enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
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
							},

							"frequency_unit": {
								Type:     pluginsdk.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									"Day",
									"Hour",
								}, false),
							},

							"keep_at_least_one_backup": {
								Type:     pluginsdk.TypeBool,
								Optional: true,
								Default:  false,
							},

							"retention_period_in_days": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								Default:      30,
								ValidateFunc: validation.IntBetween(0, 9999999),
							},

							"start_time": {
								Type:             pluginsdk.TypeString,
								Optional:         true,
								DiffSuppressFunc: suppress.RFC3339Time,
								ValidateFunc:     validation.IsRFC3339Time,
							},
						},
					},
				},
			},
		},
	}
}

func expandAppServiceBackup(input []interface{}) *web.BackupRequest {
	if len(input) == 0 {
		return nil
	}

	vals := input[0].(map[string]interface{})

	name := vals["name"].(string)
	storageAccountUrl := vals["storage_account_url"].(string)
	enabled := vals["enabled"].(bool)

	request := &web.BackupRequest{
		BackupRequestProperties: &web.BackupRequestProperties{
			BackupName:        utils.String(name),
			StorageAccountURL: utils.String(storageAccountUrl),
			Enabled:           utils.Bool(enabled),
		},
	}

	scheduleRaw := vals["schedule"].([]interface{})
	if len(scheduleRaw) > 0 {
		schedule := scheduleRaw[0].(map[string]interface{})
		backupSchedule := web.BackupSchedule{}

		if v, ok := schedule["frequency_interval"].(int); ok {
			backupSchedule.FrequencyInterval = utils.Int32(int32(v))
		}

		if v, ok := schedule["frequency_unit"]; ok {
			backupSchedule.FrequencyUnit = web.FrequencyUnit(v.(string))
		}

		if v, ok := schedule["keep_at_least_one_backup"]; ok {
			backupSchedule.KeepAtLeastOneBackup = utils.Bool(v.(bool))
		}

		if v, ok := schedule["retention_period_in_days"].(int); ok {
			backupSchedule.RetentionPeriodInDays = utils.Int32(int32(v))
		}

		if v, ok := schedule["start_time"].(string); ok {
			dateTimeToStart, _ := time.Parse(time.RFC3339, v) // validated by schema
			backupSchedule.StartTime = &date.Time{Time: dateTimeToStart}
		}

		request.BackupRequestProperties.BackupSchedule = &backupSchedule
	}

	return request
}

func flattenAppServiceBackup(input *web.BackupRequestProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if input.BackupName != nil {
		output["name"] = *input.BackupName
	}
	if input.Enabled != nil {
		output["enabled"] = *input.Enabled
	}
	if input.StorageAccountURL != nil {
		output["storage_account_url"] = *input.StorageAccountURL
	}

	schedules := make([]interface{}, 0)
	if input.BackupSchedule != nil {
		v := *input.BackupSchedule

		schedule := make(map[string]interface{})

		if v.FrequencyInterval != nil {
			schedule["frequency_interval"] = int(*v.FrequencyInterval)
		}

		schedule["frequency_unit"] = string(v.FrequencyUnit)

		if v.KeepAtLeastOneBackup != nil {
			schedule["keep_at_least_one_backup"] = *v.KeepAtLeastOneBackup
		}
		if v.RetentionPeriodInDays != nil {
			schedule["retention_period_in_days"] = int(*v.RetentionPeriodInDays)
		}
		if v.StartTime != nil && !v.StartTime.IsZero() {
			schedule["start_time"] = v.StartTime.Format(time.RFC3339)
		}

		schedules = append(schedules, schedule)
	}
	output["schedule"] = schedules

	return []interface{}{
		output,
	}
}
