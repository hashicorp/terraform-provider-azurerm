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
	Backup           []Backup           `tfschema:"backup"`
	PolicyType       string             `tfschema:"policy_type"`
	RetentionDaily   []RetentionDaily   `tfschema:"retention_daily"`
	RetentionWeekly  []RetentionWeekly  `tfschema:"retention_weekly"`
	RetentionMonthly []RetentionMonthly `tfschema:"retention_monthly"`
	RetentionYearly  []RetentionYearly  `tfschema:"retention_yearly"`
	SimpleRetention  []SimpleRetention  `tfschema:"simple_retention"`
}

type Backup struct {
	Frequency          string         `tfschema:"frequency"`
	Time               string         `tfschema:"time"`
	FrequencyInMinutes *int32         `tfschema:"frequency_in_minutes"`
	Weekdays           *pluginsdk.Set `tfschema:"weekdays"`
}

type RetentionDaily struct {
	Count int32 `tfschema:"count"`
}

type RetentionWeekly struct {
	Count    int32          `tfschema:"count"`
	Weekdays *pluginsdk.Set `tfschema:"weekdays"`
}

type RetentionMonthly struct {
	Count      int32          `tfschema:"count"`
	FormatType string         `tfschema:"format_type"`
	MonthDays  []MonthDay     `tfschema:"month_day"`
	Weeks      *pluginsdk.Set `tfschema:"weeks"`
	Weekdays   *pluginsdk.Set `tfschema:"weekdays"`
}

type MonthDay struct {
	Date   int32 `tfschema:"date"`
	IsLast bool  `tfschema:"is_last"`
}

type RetentionYearly struct {
	Count      int32          `tfschema:"count"`
	FormatType string         `tfschema:"format_type"`
	MonthDays  []MonthDay     `tfschema:"month_day"`
	Months     *pluginsdk.Set `tfschema:"months"`
	Weeks      *pluginsdk.Set `tfschema:"weeks"`
	Weekdays   *pluginsdk.Set `tfschema:"weekdays"`
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

					"retention_weekly": {
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

								"weekdays": {
									Type:     pluginsdk.TypeSet,
									Required: true,
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

					"retention_monthly": {
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

								"format_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backup.RetentionScheduleFormatDaily),
										string(backup.RetentionScheduleFormatWeekly),
									}, false),
								},

								"month_day": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"date": {
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntBetween(0, 28),
											},

											"is_last": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
										},
									},
								},

								"weeks": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Set:      set.HashStringIgnoreCase,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											string(backup.WeekOfMonthFirst),
											string(backup.WeekOfMonthSecond),
											string(backup.WeekOfMonthThird),
											string(backup.WeekOfMonthFourth),
											string(backup.WeekOfMonthLast),
										}, false),
									},
								},

								"weekdays": {
									Type:     pluginsdk.TypeSet,
									Required: true,
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

					"retention_yearly": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 9999),
								},

								"format_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backup.RetentionScheduleFormatDaily),
										string(backup.RetentionScheduleFormatWeekly),
									}, false),
								},

								"month_day": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"date": {
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntBetween(0, 28),
											},

											"is_last": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
										},
									},
								},

								"months": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Set:      set.HashStringIgnoreCase,
									Elem: &pluginsdk.Schema{
										Type:             pluginsdk.TypeString,
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc:     validation.IsMonth(true),
									},
								},

								"weeks": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									Set:      set.HashStringIgnoreCase,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											string(backup.WeekOfMonthFirst),
											string(backup.WeekOfMonthSecond),
											string(backup.WeekOfMonthThird),
											string(backup.WeekOfMonthFourth),
											string(backup.WeekOfMonthLast),
										}, false),
									},
								},

								"weekdays": {
									Type:     pluginsdk.TypeSet,
									Required: true,
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
			PolicyType:       string(policyType),
			Backup:           flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(item.SchedulePolicy, policyType),
			RetentionDaily:   flattenBackupProtectionPolicyVMWorkloadRetentionDaily(item.RetentionPolicy),
			RetentionWeekly:  flattenBackupProtectionPolicyVMWorkloadRetentionWeekly(item.RetentionPolicy),
			RetentionMonthly: flattenBackupProtectionPolicyVMWorkloadRetentionMonthly(item.RetentionPolicy),
			RetentionYearly:  flattenBackupProtectionPolicyVMWorkloadRetentionYearly(item.RetentionPolicy),
			SimpleRetention:  flattenBackupProtectionPolicyVMWorkloadSimpleRetention(item.RetentionPolicy),
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
			for _, day := range v.List() {
				days = append(days, backup.DayOfWeek(day.(string)))
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
			weekdays := make([]interface{}, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			backupBlock.Weekdays = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
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

			if len(input.RetentionWeekly) > 0 {
				retentionWeekly := input.RetentionWeekly[0]

				retentionPolicy.WeeklySchedule = &backup.WeeklyRetentionSchedule{
					RetentionTimes: &times,
					RetentionDuration: &backup.RetentionDuration{
						Count:        utils.Int32(retentionWeekly.Count),
						DurationType: backup.RetentionDurationTypeWeeks,
					},
				}

				if v := retentionWeekly.Weekdays; v != nil {
					days := make([]backup.DayOfWeek, 0)
					for _, day := range v.List() {
						days = append(days, backup.DayOfWeek(day.(string)))
					}
					retentionPolicy.WeeklySchedule.DaysOfTheWeek = &days
				}
			}

			if len(input.RetentionMonthly) > 0 {
				retentionMonthly := input.RetentionMonthly[0]

				retentionPolicy.MonthlySchedule = &backup.MonthlyRetentionSchedule{
					RetentionScheduleFormatType: backup.RetentionScheduleFormat(retentionMonthly.FormatType),
					RetentionScheduleDaily:      expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(retentionMonthly.MonthDays),
					RetentionScheduleWeekly:     expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(retentionMonthly.Weekdays, retentionMonthly.Weeks),
					RetentionTimes:              &times,
					RetentionDuration: &backup.RetentionDuration{
						Count:        utils.Int32(retentionMonthly.Count),
						DurationType: backup.RetentionDurationTypeMonths,
					},
				}
			}

			if len(input.RetentionYearly) > 0 {
				retentionYearly := input.RetentionYearly[0]

				retentionPolicy.YearlySchedule = &backup.YearlyRetentionSchedule{
					RetentionScheduleFormatType: backup.RetentionScheduleFormat(retentionYearly.FormatType),
					RetentionScheduleDaily:      expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(retentionYearly.MonthDays),
					RetentionScheduleWeekly:     expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(retentionYearly.Weekdays, retentionYearly.Weeks),
					RetentionTimes:              &times,
					RetentionDuration: &backup.RetentionDuration{
						Count:        utils.Int32(retentionYearly.Count),
						DurationType: backup.RetentionDurationTypeYears,
					},
				}

				if v := retentionYearly.Months; v != nil {
					months := make([]backup.MonthOfYear, 0)
					for _, month := range v.List() {
						months = append(months, backup.MonthOfYear(month.(string)))
					}
					retentionPolicy.YearlySchedule.MonthsOfYear = &months
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

func flattenBackupProtectionPolicyVMWorkloadRetentionWeekly(input backup.BasicRetentionPolicy) []RetentionWeekly {
	if input == nil {
		return nil
	}

	longTermRetentionPolicy, _ := input.AsLongTermRetentionPolicy()
	retentionWeeklyBlock := RetentionWeekly{}

	if weeklySchedule := longTermRetentionPolicy.WeeklySchedule; weeklySchedule != nil {
		if duration := weeklySchedule.RetentionDuration; duration != nil {
			if v := duration.Count; v != nil {
				retentionWeeklyBlock.Count = *v
			}
		}

		if days := weeklySchedule.DaysOfTheWeek; days != nil {
			weekdays := make([]interface{}, 0)
			for _, d := range *days {
				weekdays = append(weekdays, string(d))
			}
			retentionWeeklyBlock.Weekdays = pluginsdk.NewSet(pluginsdk.HashString, weekdays)
		}
	}

	return []RetentionWeekly{retentionWeeklyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionMonthly(input backup.BasicRetentionPolicy) []RetentionMonthly {
	if input == nil {
		return nil
	}

	longTermRetentionPolicy, _ := input.AsLongTermRetentionPolicy()
	retentionMonthlyBlock := RetentionMonthly{}

	if monthlySchedule := longTermRetentionPolicy.MonthlySchedule; monthlySchedule != nil {
		if duration := monthlySchedule.RetentionDuration; duration != nil {
			if v := duration.Count; v != nil {
				retentionMonthlyBlock.Count = *v
			}
		}

		if formatType := monthlySchedule.RetentionScheduleFormatType; formatType != "" {
			retentionMonthlyBlock.FormatType = string(formatType)
		}

		if weekly := monthlySchedule.RetentionScheduleWeekly; weekly != nil {
			retentionMonthlyBlock.Weekdays, retentionMonthlyBlock.Weeks = flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekly)
		}

		if daily := monthlySchedule.RetentionScheduleDaily; daily != nil {
			retentionMonthlyBlock.MonthDays = flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(daily)
		}
	}

	return []RetentionMonthly{retentionMonthlyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionYearly(input backup.BasicRetentionPolicy) []RetentionYearly {
	if input == nil {
		return nil
	}

	longTermRetentionPolicy, _ := input.AsLongTermRetentionPolicy()
	retentionYearlyBlock := RetentionYearly{}

	if yearlySchedule := longTermRetentionPolicy.YearlySchedule; yearlySchedule != nil {
		if duration := yearlySchedule.RetentionDuration; duration != nil {
			if v := duration.Count; v != nil {
				retentionYearlyBlock.Count = *v
			}
		}

		if formatType := yearlySchedule.RetentionScheduleFormatType; formatType != "" {
			retentionYearlyBlock.FormatType = string(formatType)
		}

		if weekly := yearlySchedule.RetentionScheduleWeekly; weekly != nil {
			retentionYearlyBlock.Weekdays, retentionYearlyBlock.Weeks = flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekly)
		}

		if months := yearlySchedule.MonthsOfYear; months != nil {
			slice := make([]interface{}, 0)
			for _, d := range *months {
				slice = append(slice, string(d))
			}
			retentionYearlyBlock.Months = pluginsdk.NewSet(pluginsdk.HashString, slice)
		}

		if daily := yearlySchedule.RetentionScheduleDaily; daily != nil {
			retentionYearlyBlock.MonthDays = flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(daily)
		}
	}

	return []RetentionYearly{retentionYearlyBlock}
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

func expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(monthdays []MonthDay) *backup.DailyRetentionFormat {
	daily := backup.DailyRetentionFormat{}

	if v := monthdays; v != nil {
		days := make([]backup.Day, 0)
		for _, day := range v {
			days = append(days, backup.Day{
				Date:   &day.Date,
				IsLast: &day.IsLast,
			})
		}
		daily.DaysOfTheMonth = &days
	}

	return &daily
}

func expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekdays, weeks *pluginsdk.Set) *backup.WeeklyRetentionFormat {
	weekly := backup.WeeklyRetentionFormat{}

	if v := weekdays; v != nil {
		days := make([]backup.DayOfWeek, 0)
		for _, day := range v.List() {
			days = append(days, backup.DayOfWeek(day.(string)))
		}
		weekly.DaysOfTheWeek = &days
	}

	if v := weeks; v != nil {
		weeks := make([]backup.WeekOfMonth, 0)
		for _, week := range v.List() {
			weeks = append(weeks, backup.WeekOfMonth(week.(string)))
		}
		weekly.WeeksOfTheMonth = &weeks
	}

	return &weekly
}

func flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(input *backup.WeeklyRetentionFormat) (weekdays, weeks *pluginsdk.Set) {
	if days := input.DaysOfTheWeek; days != nil {
		slice := make([]interface{}, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weekdays = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	if days := input.WeeksOfTheMonth; days != nil {
		slice := make([]interface{}, 0)
		for _, d := range *days {
			slice = append(slice, string(d))
		}
		weeks = pluginsdk.NewSet(pluginsdk.HashString, slice)
	}

	return weekdays, weeks
}

func flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(input *backup.DailyRetentionFormat) []MonthDay {
	monthdays := make([]MonthDay, 0)
	if days := input.DaysOfTheMonth; days != nil {
		for _, d := range *days {
			monthdays = append(monthdays, MonthDay{
				Date:   *d.Date,
				IsLast: *d.IsLast,
			})
		}
	}

	return monthdays
}
