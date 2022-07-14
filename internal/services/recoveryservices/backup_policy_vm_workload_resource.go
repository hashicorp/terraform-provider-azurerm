package recoveryservices

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupProtectionPolicyVMWorkloadModel struct {
	Name               string             `tfschema:"name"`
	ResourceGroupName  string             `tfschema:"resource_group_name"`
	RecoveryVaultName  string             `tfschema:"recovery_vault_name"`
	ProtectionPolicies []ProtectionPolicy `tfschema:"protection_policy"`
	Settings           []Settings         `tfschema:"settings"`
	WorkloadType       string             `tfschema:"workload_type"`
}

type ProtectionPolicy struct {
	Backup          []Backup          `tfschema:"backup"`
	PolicyType      string            `tfschema:"policy_type"`
	RetentionDaily  []RetentionDaily  `tfschema:"retention_daily"`
	SimpleRetention []SimpleRetention `tfschema:"simple_retention"`
}

type Backup struct {
	Frequency          string    `tfschema:"frequency"`
	Time               string    `tfschema:"time"`
	FrequencyInMinutes *int32    `tfschema:"frequency_in_minutes"`
	Weekdays           *[]string `tfschema:"weekdays"`
}

type RetentionDaily struct {
	Count int32 `tfschema:"count"`
}

type SimpleRetention struct {
	Count int32 `tfschema:"count"`
}

type Settings struct {
	CompressionEnabled    *bool   `tfschema:"compression_enabled"`
	SqlCompressionEnabled *bool   `tfschema:"sql_compression_enabled"`
	TimeZone              *string `tfschema:"time_zone"`
}

type BackupProtectionPolicyVMWorkloadResource struct{}

var _ sdk.ResourceWithUpdate = BackupProtectionPolicyVMWorkloadResource{}

func (r BackupProtectionPolicyVMWorkloadResource) ResourceType() string {
	return "azurerm_backup_policy_vm_workload"
}

func (r BackupProtectionPolicyVMWorkloadResource) ModelObject() interface{} {
	return &BackupProtectionPolicyVMWorkloadModel{}
}

func (r BackupProtectionPolicyVMWorkloadResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BackupPolicyID
}

func (r BackupProtectionPolicyVMWorkloadResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"protection_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(backup.PolicyTypeDifferential),
							string(backup.PolicyTypeFull),
							string(backup.PolicyTypeIncremental),
							string(backup.PolicyTypeLog),
						}, false),
					},

					"backup": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"frequency": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backup.ScheduleRunTypeDaily),
										string(backup.ScheduleRunTypeWeekly),
									}, false),
								},

								"time": { // applies to all backup schedules & retention times (they all must be the same)
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
										"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
									),
								},

								"frequency_in_minutes": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ValidateFunc: validation.IntInSlice([]int{
										15,
										30,
										60,
										120,
										240,
										480,
										720,
										1440,
									}),
								},

								"weekdays": { // only for weekly
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Set:      set.HashStringIgnoreCase,
									Elem: &pluginsdk.Schema{
										Type:             pluginsdk.TypeString,
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc:     validation.IsDayOfTheWeek(true),
									},
								},
							},
						},
					},

					"retention_daily": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(7, 9999),
								},
							},
						},
					},

					"simple_retention": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(7, 35),
								},
							},
						},
					},
				},
			},
		},

		"settings": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"sql_compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"time_zone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"workload_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(backup.WorkloadTypeSQLDataBase),
				string(backup.WorkloadTypeSAPHanaDatabase),
			}, false),
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BackupProtectionPolicyVMWorkloadResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupProtectionPolicyVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.RecoveryVaultName, model.Name)

			existing, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return tf.ImportAsExistsError("azurerm_backup_policy_vm_workload", id.ID())
			}

			protectionPolicy, err := expandBackupProtectionPolicyVMWorkloadProtectionPolicies(model.ProtectionPolicies)
			if err != nil {
				return err
			}

			properties := &backup.ProtectionPolicyResource{
				Properties: &backup.AzureVMWorkloadProtectionPolicy{
					BackupManagementType: backup.ManagementTypeBasicProtectionPolicyBackupManagementTypeAzureWorkload,
					Settings:             expandBackupProtectionPolicyVMWorkloadSettings(model.Settings),
					SubProtectionPolicy:  protectionPolicy,
					WorkLoadType:         backup.WorkloadType(model.WorkloadType),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id.VaultName, id.ResourceGroup, id.Name, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := parse.BackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupProtectionPolicyVMWorkloadModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if props := existing.Properties; props != nil {
				vmWorkload, _ := props.AsAzureVMWorkloadProtectionPolicy()

				if metadata.ResourceData.HasChange("settings") {
					vmWorkload.Settings = expandBackupProtectionPolicyVMWorkloadSettings(model.Settings)
				}

				if metadata.ResourceData.HasChange("protection_policy") {
					protectionPolicy, err := expandBackupProtectionPolicyVMWorkloadProtectionPolicies(model.ProtectionPolicies)
					if err != nil {
						return err
					}

					vmWorkload.SubProtectionPolicy = protectionPolicy
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id.VaultName, id.ResourceGroup, id.Name, existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := parse.BackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BackupProtectionPolicyVMWorkloadModel{
				Name:              id.Name,
				ResourceGroupName: id.ResourceGroup,
				RecoveryVaultName: id.VaultName,
			}

			if props := resp.Properties; props != nil {
				vmWorkload, _ := props.AsAzureVMWorkloadProtectionPolicy()
				state.WorkloadType = string(vmWorkload.WorkLoadType)
				state.Settings = flattenBackupProtectionPolicyVMWorkloadSettings(vmWorkload.Settings)
				state.ProtectionPolicies = flattenBackupProtectionPolicyVMWorkloadProtectionPolicies(vmWorkload.SubProtectionPolicy)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BackupProtectionPolicyVMWorkloadResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ProtectionPoliciesClient

			id, err := parse.BackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.VaultName, id.ResourceGroup, id.Name)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandBackupProtectionPolicyVMWorkloadSettings(input []Settings) *backup.Settings {
	if len(input) == 0 {
		return &backup.Settings{}
	}

	result := &backup.Settings{}
	settings := input[0]

	if settings.CompressionEnabled != nil {
		result.IsCompression = settings.CompressionEnabled
	}

	if settings.SqlCompressionEnabled != nil {
		result.Issqlcompression = settings.SqlCompressionEnabled
	}

	if settings.TimeZone != nil {
		result.TimeZone = settings.TimeZone
	}

	return result
}

func flattenBackupProtectionPolicyVMWorkloadSettings(input *backup.Settings) []Settings {
	if input == nil {
		return make([]Settings, 0)
	}

	result := make([]Settings, 0)

	result = append(result, Settings{
		CompressionEnabled:    input.IsCompression,
		SqlCompressionEnabled: input.Issqlcompression,
		TimeZone:              input.TimeZone,
	})

	return result
}

func expandBackupProtectionPolicyVMWorkloadProtectionPolicies(input []ProtectionPolicy) (*[]backup.SubProtectionPolicy, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]backup.SubProtectionPolicy, 0)

	for _, item := range input {
		// getting this ready now because its shared between *everything*, time is... complicated for this resource
		timeOfDay := item.Backup[0].Time
		dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
		if err != nil {
			return nil, fmt.Errorf("generating time from %q for policy): %+v", timeOfDay, err)
		}
		times := append(make([]date.Time, 0), date.Time{Time: dateOfDay})

		results = append(results, backup.SubProtectionPolicy{
			PolicyType:      backup.PolicyType(item.PolicyType),
			RetentionPolicy: expandBackupProtectionPolicyVMWorkloadRetentionPolicy(item, times),
			SchedulePolicy:  expandBackupProtectionPolicyVMWorkloadSchedulePolicy(item, times),
		})
	}

	return &results, nil
}

func flattenBackupProtectionPolicyVMWorkloadProtectionPolicies(input *[]backup.SubProtectionPolicy) []ProtectionPolicy {
	results := make([]ProtectionPolicy, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var policyType backup.PolicyType
		if item.PolicyType != "" {
			policyType = item.PolicyType
		}

		results = append(results, ProtectionPolicy{
			PolicyType:      string(policyType),
			Backup:          flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(item.SchedulePolicy, policyType),
			RetentionDaily:  flattenBackupProtectionPolicyVMWorkloadRetentionDaily(item.RetentionPolicy),
			SimpleRetention: flattenBackupProtectionPolicyVMWorkloadSimpleRetention(item.RetentionPolicy),
		})
	}

	return results
}

func expandBackupProtectionPolicyVMWorkloadSchedulePolicy(input ProtectionPolicy, times []date.Time) backup.BasicSchedulePolicy {
	if input.PolicyType == string(backup.PolicyTypeLog) {
		schedule := backup.LogSchedulePolicy{
			SchedulePolicyType: backup.SchedulePolicyTypeLogSchedulePolicy,
		}

		backupBlock := input.Backup[0]

		if v := backupBlock.FrequencyInMinutes; v != nil {
			schedule.ScheduleFrequencyInMins = v
		}

		result, _ := schedule.AsBasicSchedulePolicy()
		return result
	} else {
		schedule := backup.SimpleSchedulePolicy{
			SchedulePolicyType: backup.SchedulePolicyTypeSimpleSchedulePolicy,
			ScheduleRunTimes:   &times,
		}

		backupBlock := input.Backup[0]
		schedule.ScheduleRunFrequency = backup.ScheduleRunType(backupBlock.Frequency)

		if v := backupBlock.Weekdays; v != nil {
			days := make([]backup.DayOfWeek, 0)
			for _, day := range *v {
				days = append(days, backup.DayOfWeek(day))
			}
			schedule.ScheduleRunDays = &days
		}

		result, _ := schedule.AsBasicSchedulePolicy()
		return result
	}

	return nil
}

func flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(input backup.BasicSchedulePolicy, policyType backup.PolicyType) []Backup {
	if input == nil {
		return nil
	}

	backupBlock := Backup{}

	if policyType == backup.PolicyTypeLog {
		logSchedulePolicy, _ := input.AsLogSchedulePolicy()

		if v := logSchedulePolicy.ScheduleFrequencyInMins; v != nil {
			backupBlock.FrequencyInMinutes = v
		}
	} else {
		simpleSchedulePolicy, _ := input.AsSimpleSchedulePolicy()

		backupBlock.Frequency = string(simpleSchedulePolicy.ScheduleRunFrequency)

		if times := simpleSchedulePolicy.ScheduleRunTimes; times != nil && len(*times) > 0 {
			backupBlock.Time = (*times)[0].Format("15:04")
		}

		if days := simpleSchedulePolicy.ScheduleRunDays; days != nil {
			weekdays := make([]string, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			backupBlock.Weekdays = &weekdays
		}
	}

	return []Backup{backupBlock}
}

func expandBackupProtectionPolicyVMWorkloadRetentionPolicy(input ProtectionPolicy, times []date.Time) backup.BasicRetentionPolicy {
	if input.PolicyType == string(backup.PolicyTypeFull) {
		if input.RetentionDaily != nil {
			retentionPolicy := backup.LongTermRetentionPolicy{
				RetentionPolicyType: backup.RetentionPolicyTypeLongTermRetentionPolicy,
			}

			if len(input.RetentionDaily) > 0 {
				retentionDaily := input.RetentionDaily[0]
				retentionPolicy.DailySchedule = &backup.DailyRetentionSchedule{
					RetentionTimes: &times,
					RetentionDuration: &backup.RetentionDuration{
						Count:        utils.Int32(retentionDaily.Count),
						DurationType: backup.RetentionDurationTypeDays,
					},
				}
			}

			return retentionPolicy
		}
	} else {
		if input.SimpleRetention != nil {
			retentionPolicy := backup.SimpleRetentionPolicy{
				RetentionPolicyType: backup.RetentionPolicyTypeSimpleRetentionPolicy,
			}

			if len(input.SimpleRetention) > 0 {
				simpleRetention := input.SimpleRetention[0]
				retentionPolicy.RetentionDuration = &backup.RetentionDuration{
					Count:        utils.Int32(simpleRetention.Count),
					DurationType: backup.RetentionDurationTypeDays,
				}
			}

			return retentionPolicy
		}
	}

	return nil
}

func flattenBackupProtectionPolicyVMWorkloadRetentionDaily(input backup.BasicRetentionPolicy) []RetentionDaily {
	if input == nil {
		return nil
	}

	longTermRetentionPolicy, _ := input.AsLongTermRetentionPolicy()
	retentionDailyBlock := RetentionDaily{}

	if dailySchedule := longTermRetentionPolicy.DailySchedule; dailySchedule != nil {
		if duration := dailySchedule.RetentionDuration; duration != nil {
			if v := duration.Count; v != nil {
				retentionDailyBlock.Count = *v
			}
		}
	}

	return []RetentionDaily{retentionDailyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadSimpleRetention(input backup.BasicRetentionPolicy) []SimpleRetention {
	if input == nil {
		return nil
	}

	simpleRetentionPolicy, _ := input.AsSimpleRetentionPolicy()
	simpleRetentionBlock := SimpleRetention{}

	if duration := simpleRetentionPolicy.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			simpleRetentionBlock.Count = *v
		}
	}

	return []SimpleRetention{simpleRetentionBlock}
}
