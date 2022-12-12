package recoveryservices

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
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
	Frequency          string   `tfschema:"frequency"`
	FrequencyInMinutes int32    `tfschema:"frequency_in_minutes"`
	Time               string   `tfschema:"time"`
	Weekdays           []string `tfschema:"weekdays"`
}

type RetentionDaily struct {
	Count int32 `tfschema:"count"`
}

type RetentionWeekly struct {
	Count    int32    `tfschema:"count"`
	Weekdays []string `tfschema:"weekdays"`
}

type RetentionMonthly struct {
	Count      int32    `tfschema:"count"`
	FormatType string   `tfschema:"format_type"`
	Monthdays  []int    `tfschema:"monthdays"`
	Weeks      []string `tfschema:"weeks"`
	Weekdays   []string `tfschema:"weekdays"`
}

type RetentionYearly struct {
	Count      int32    `tfschema:"count"`
	FormatType string   `tfschema:"format_type"`
	Months     []string `tfschema:"months"`
	Monthdays  []int    `tfschema:"monthdays"`
	Weeks      []string `tfschema:"weeks"`
	Weekdays   []string `tfschema:"weekdays"`
}

type SimpleRetention struct {
	Count int32 `tfschema:"count"`
}

type Settings struct {
	CompressionEnabled bool   `tfschema:"compression_enabled"`
	TimeZone           string `tfschema:"time_zone"`
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
			Type:     pluginsdk.TypeSet,
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
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"frequency": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backup.ScheduleRunTypeDaily),
										string(backup.ScheduleRunTypeWeekly),
									}, false),
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

								"time": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"),
										"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
									),
								},

								"weekdays": {
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
						Optional: true,
						MaxItems: 1,
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
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5163),
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
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 1188),
								},

								"format_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backup.RetentionScheduleFormatDaily),
										string(backup.RetentionScheduleFormatWeekly),
									}, false),
								},

								"monthdays": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 28),
									},
								},

								"weeks": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
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

					"retention_yearly": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"count": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 99),
								},

								"format_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backup.RetentionScheduleFormatDaily),
										string(backup.RetentionScheduleFormatWeekly),
									}, false),
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

								"monthdays": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 28),
									},
								},

								"weeks": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
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

					"simple_retention": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
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
					"time_zone": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"compression_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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

			protectionPolicy, err := expandBackupProtectionPolicyVMWorkloadProtectionPolicies(model.ProtectionPolicies, model.WorkloadType)
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

			props, _ := existing.Properties.AsAzureVMWorkloadProtectionPolicy()

			if metadata.ResourceData.HasChange("settings") {
				props.Settings = expandBackupProtectionPolicyVMWorkloadSettings(model.Settings)
			}

			if metadata.ResourceData.HasChange("protection_policy") {
				protectionPolicy, err := expandBackupProtectionPolicyVMWorkloadProtectionPolicies(model.ProtectionPolicies, model.WorkloadType)
				if err != nil {
					return err
				}

				props.SubProtectionPolicy = protectionPolicy
			}

			existing.Properties = props

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

	settings := input[0]
	result := &backup.Settings{
		IsCompression: utils.Bool(settings.CompressionEnabled),
	}

	if settings.TimeZone != "" {
		result.TimeZone = utils.String(settings.TimeZone)
	}

	return result
}

func flattenBackupProtectionPolicyVMWorkloadSettings(input *backup.Settings) []Settings {
	if input == nil {
		return make([]Settings, 0)
	}

	result := make([]Settings, 0)

	result = append(result, Settings{
		CompressionEnabled: *input.IsCompression,
		TimeZone:           *input.TimeZone,
	})

	return result
}

func expandBackupProtectionPolicyVMWorkloadProtectionPolicies(input []ProtectionPolicy, workloadType string) (*[]backup.SubProtectionPolicy, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]backup.SubProtectionPolicy, 0)

	for _, item := range input {
		if workloadType == string(backup.WorkloadTypeSQLDataBase) && item.PolicyType == string(backup.PolicyTypeIncremental) {
			return nil, fmt.Errorf("the Incremental backup isn't supported when `workload_type` is `SQLDataBase`")
		}

		backupBlock := item.Backup[0]

		// getting this ready now because its shared between *everything*, time is... complicated for this resource
		timeOfDay := backupBlock.Time
		times := make([]date.Time, 0)
		if timeOfDay != "" {
			dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
			if err != nil {
				return nil, fmt.Errorf("generating time from %q for policy): %+v", timeOfDay, err)
			}
			times = append(times, date.Time{Time: dateOfDay})
		}

		switch backupBlock.Frequency {
		case string(backup.ScheduleRunTypeDaily):
			if item.RetentionDaily == nil || len(item.RetentionDaily) == 0 {
				return nil, fmt.Errorf("`retention_daily` must be set when `backup.0.frequency` is `Daily`")
			}

			if weekdays := backupBlock.Weekdays; len(weekdays) > 0 {
				return nil, fmt.Errorf("`backup.0.weekdays` should be not set when `backup.0.frequency` is `Daily`")
			}
		case string(backup.ScheduleRunTypeWeekly):
			if len(item.RetentionDaily) > 0 {
				return nil, fmt.Errorf("`retention_daily` must be not set when `backup.0.frequency` is `Weekly`")
			}

			if item.PolicyType != string(backup.PolicyTypeLog) && len(backupBlock.Weekdays) == 0 {
				return nil, fmt.Errorf("`backup.weekdays` must be set when `policy_type` is not `Log` and `backup.frequency` is `Weekly`")
			}

			if item.PolicyType == string(backup.PolicyTypeFull) && len(item.RetentionWeekly) == 0 {
				return nil, fmt.Errorf("`retention_weekly` must be set when `policy_type` is `Full` and `backup.frequency` is `Weekly`")
			}
		}

		result := backup.SubProtectionPolicy{
			PolicyType:     backup.PolicyType(item.PolicyType),
			SchedulePolicy: expandBackupProtectionPolicyVMWorkloadSchedulePolicy(item, times),
		}

		if v, err := expandBackupProtectionPolicyVMWorkloadRetentionPolicy(item, times); err != nil {
			return nil, err
		} else {
			result.RetentionPolicy = v
		}

		results = append(results, result)
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

		result := ProtectionPolicy{
			PolicyType: string(policyType),
			Backup:     flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(item.SchedulePolicy, policyType),
		}

		if retentionPolicy := item.RetentionPolicy; retentionPolicy != nil {
			if longTermRetentionPolicy, ok := retentionPolicy.AsLongTermRetentionPolicy(); ok {
				result.RetentionDaily = flattenBackupProtectionPolicyVMWorkloadRetentionDaily(longTermRetentionPolicy.DailySchedule)
				result.RetentionWeekly = flattenBackupProtectionPolicyVMWorkloadRetentionWeekly(longTermRetentionPolicy.WeeklySchedule)
				result.RetentionMonthly = flattenBackupProtectionPolicyVMWorkloadRetentionMonthly(longTermRetentionPolicy.MonthlySchedule)
				result.RetentionYearly = flattenBackupProtectionPolicyVMWorkloadRetentionYearly(longTermRetentionPolicy.YearlySchedule)
			} else {
				simpleRetentionPolicy, _ := retentionPolicy.AsSimpleRetentionPolicy()
				result.SimpleRetention = flattenBackupProtectionPolicyVMWorkloadSimpleRetention(simpleRetentionPolicy.RetentionDuration)
			}
		}

		results = append(results, result)
	}

	return results
}

func expandBackupProtectionPolicyVMWorkloadSchedulePolicy(input ProtectionPolicy, times []date.Time) backup.BasicSchedulePolicy {
	if input.PolicyType == string(backup.PolicyTypeLog) {
		schedule := backup.LogSchedulePolicy{
			SchedulePolicyType: backup.SchedulePolicyTypeLogSchedulePolicy,
		}

		if v := input.Backup[0].FrequencyInMinutes; v != 0 {
			schedule.ScheduleFrequencyInMins = utils.Int32(v)
		}

		result, _ := schedule.AsBasicSchedulePolicy()
		return result
	} else {
		schedule := backup.SimpleSchedulePolicy{
			SchedulePolicyType: backup.SchedulePolicyTypeSimpleSchedulePolicy,
		}

		backupBlock := input.Backup[0]
		if backupBlock.Frequency != "" {
			schedule.ScheduleRunFrequency = backup.ScheduleRunType(backupBlock.Frequency)
		}

		if len(times) > 0 {
			schedule.ScheduleRunTimes = &times
		}

		if v := backupBlock.Weekdays; len(v) > 0 {
			days := make([]backup.DayOfWeek, 0)
			for _, day := range v {
				days = append(days, backup.DayOfWeek(day))
			}
			schedule.ScheduleRunDays = &days
		}

		result, _ := schedule.AsBasicSchedulePolicy()
		return result
	}
}

func flattenBackupProtectionPolicyVMWorkloadSchedulePolicy(input backup.BasicSchedulePolicy, policyType backup.PolicyType) []Backup {
	if input == nil {
		return nil
	}

	backupBlock := Backup{}

	if policyType == backup.PolicyTypeLog {
		logSchedulePolicy, _ := input.AsLogSchedulePolicy()

		if v := logSchedulePolicy.ScheduleFrequencyInMins; v != nil {
			backupBlock.FrequencyInMinutes = *v
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
			backupBlock.Weekdays = weekdays
		}
	}

	return []Backup{backupBlock}
}

func expandBackupProtectionPolicyVMWorkloadRetentionPolicy(input ProtectionPolicy, times []date.Time) (backup.BasicRetentionPolicy, error) {
	if input.PolicyType == string(backup.PolicyTypeFull) {
		retentionPolicy := backup.LongTermRetentionPolicy{
			RetentionPolicyType: backup.RetentionPolicyTypeLongTermRetentionPolicy,
		}

		if input.RetentionDaily != nil && len(input.RetentionDaily) > 0 {
			retentionDaily := input.RetentionDaily[0]

			retentionPolicy.DailySchedule = &backup.DailyRetentionSchedule{
				RetentionTimes: &times,
				RetentionDuration: &backup.RetentionDuration{
					Count:        utils.Int32(retentionDaily.Count),
					DurationType: backup.RetentionDurationTypeDays,
				},
			}
		}

		if input.RetentionWeekly != nil && len(input.RetentionWeekly) > 0 {
			retentionWeekly := input.RetentionWeekly[0]

			retentionPolicy.WeeklySchedule = &backup.WeeklyRetentionSchedule{
				RetentionTimes: &times,
				RetentionDuration: &backup.RetentionDuration{
					Count:        utils.Int32(retentionWeekly.Count),
					DurationType: backup.RetentionDurationTypeWeeks,
				},
			}

			if v := retentionWeekly.Weekdays; len(v) > 0 {
				days := make([]backup.DayOfWeek, 0)
				for _, day := range v {
					days = append(days, backup.DayOfWeek(day))
				}
				retentionPolicy.WeeklySchedule.DaysOfTheWeek = &days
			}
		}

		if input.RetentionMonthly != nil && len(input.RetentionMonthly) > 0 {
			retentionMonthly := input.RetentionMonthly[0]

			if input.Backup[0].Frequency == string(backup.ScheduleRunTypeWeekly) && retentionMonthly.FormatType != string(backup.RetentionScheduleFormatWeekly) {
				return nil, fmt.Errorf("`retention_monthly.format_type` must be `Weekly` when `policy_type` is `Full` and `frequency` is `Weekly`")
			}

			if retentionMonthly.FormatType == string(backup.RetentionScheduleFormatDaily) && (retentionMonthly.Monthdays == nil || len(retentionMonthly.Monthdays) == 0) {
				return nil, fmt.Errorf("`retention_monthly.monthdays` must be set when `retention_monthly.format_type` is `Daily`")
			}

			if retentionMonthly.FormatType == string(backup.RetentionScheduleFormatWeekly) && ((retentionMonthly.Weeks == nil || len(retentionMonthly.Weeks) == 0) || (retentionMonthly.Weekdays == nil || len(retentionMonthly.Weekdays) == 0)) {
				return nil, fmt.Errorf("`retention_monthly.weeks` and `retention_monthly.weekdays` must be set when `retention_monthly.format_type` is `Weekly`")
			}

			retentionPolicy.MonthlySchedule = &backup.MonthlyRetentionSchedule{
				RetentionScheduleFormatType: backup.RetentionScheduleFormat(retentionMonthly.FormatType),
				RetentionScheduleDaily:      expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(retentionMonthly.Monthdays),
				RetentionScheduleWeekly:     expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(retentionMonthly.Weekdays, retentionMonthly.Weeks),
				RetentionTimes:              &times,
				RetentionDuration: &backup.RetentionDuration{
					Count:        utils.Int32(retentionMonthly.Count),
					DurationType: backup.RetentionDurationTypeMonths,
				},
			}
		}

		if input.RetentionYearly != nil && len(input.RetentionYearly) > 0 {
			retentionYearly := input.RetentionYearly[0]

			if input.Backup[0].Frequency == string(backup.ScheduleRunTypeWeekly) && retentionYearly.FormatType != string(backup.RetentionScheduleFormatWeekly) {
				return nil, fmt.Errorf("`retention_yearly.format_type` must be `Weekly` when `policy_type` is `Full` and `frequency` is `Weekly`")
			}

			if retentionYearly.FormatType == string(backup.RetentionScheduleFormatDaily) && (retentionYearly.Monthdays == nil || len(retentionYearly.Monthdays) == 0) {
				return nil, fmt.Errorf("`retention_yearly.monthdays` must be set when `retention_yearly.format_type` is `Daily`")
			}

			if retentionYearly.FormatType == string(backup.RetentionScheduleFormatWeekly) && ((retentionYearly.Weeks == nil || len(retentionYearly.Weeks) == 0) || (retentionYearly.Weekdays == nil || len(retentionYearly.Weekdays) == 0)) {
				return nil, fmt.Errorf("`retention_yearly.weeks` and `retention_yearly.weekdays` must be set when `retention_yearly.format_type` is `Weekly`")
			}

			retentionPolicy.YearlySchedule = &backup.YearlyRetentionSchedule{
				RetentionScheduleFormatType: backup.RetentionScheduleFormat(retentionYearly.FormatType),
				RetentionScheduleDaily:      expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(retentionYearly.Monthdays),
				RetentionScheduleWeekly:     expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(retentionYearly.Weekdays, retentionYearly.Weeks),
				RetentionTimes:              &times,
				RetentionDuration: &backup.RetentionDuration{
					Count:        utils.Int32(retentionYearly.Count),
					DurationType: backup.RetentionDurationTypeYears,
				},
			}

			if v := retentionYearly.Months; v != nil {
				months := make([]backup.MonthOfYear, 0)
				for _, month := range v {
					months = append(months, backup.MonthOfYear(month))
				}
				retentionPolicy.YearlySchedule.MonthsOfYear = &months
			}
		}

		return retentionPolicy, nil
	} else {
		retentionPolicy := backup.SimpleRetentionPolicy{
			RetentionPolicyType: backup.RetentionPolicyTypeSimpleRetentionPolicy,
		}

		if input.SimpleRetention != nil && len(input.SimpleRetention) > 0 {
			simpleRetention := input.SimpleRetention[0]

			retentionPolicy.RetentionDuration = &backup.RetentionDuration{
				Count:        utils.Int32(simpleRetention.Count),
				DurationType: backup.RetentionDurationTypeDays,
			}
		}

		return retentionPolicy, nil
	}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionDaily(input *backup.DailyRetentionSchedule) []RetentionDaily {
	if input == nil {
		return nil
	}

	retentionDailyBlock := RetentionDaily{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionDailyBlock.Count = *v
		}
	}

	return []RetentionDaily{retentionDailyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionWeekly(input *backup.WeeklyRetentionSchedule) []RetentionWeekly {
	if input == nil {
		return nil
	}

	retentionWeeklyBlock := RetentionWeekly{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionWeeklyBlock.Count = *v
		}
	}

	if days := input.DaysOfTheWeek; days != nil {
		weekdays := make([]string, 0)
		for _, d := range *days {
			weekdays = append(weekdays, string(d))
		}
		retentionWeeklyBlock.Weekdays = weekdays
	}

	return []RetentionWeekly{retentionWeeklyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionMonthly(input *backup.MonthlyRetentionSchedule) []RetentionMonthly {
	if input == nil {
		return nil
	}

	retentionMonthlyBlock := RetentionMonthly{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionMonthlyBlock.Count = *v
		}
	}

	if formatType := input.RetentionScheduleFormatType; formatType != "" {
		retentionMonthlyBlock.FormatType = string(formatType)
	}

	if weekly := input.RetentionScheduleWeekly; weekly != nil {
		retentionMonthlyBlock.Weekdays, retentionMonthlyBlock.Weeks = flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekly)
	}

	if daily := input.RetentionScheduleDaily; daily != nil {
		retentionMonthlyBlock.Monthdays = flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(daily)
	}

	return []RetentionMonthly{retentionMonthlyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadRetentionYearly(input *backup.YearlyRetentionSchedule) []RetentionYearly {
	if input == nil {
		return nil
	}

	retentionYearlyBlock := RetentionYearly{}

	if duration := input.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			retentionYearlyBlock.Count = *v
		}
	}

	if formatType := input.RetentionScheduleFormatType; formatType != "" {
		retentionYearlyBlock.FormatType = string(formatType)
	}

	if weekly := input.RetentionScheduleWeekly; weekly != nil {
		retentionYearlyBlock.Weekdays, retentionYearlyBlock.Weeks = flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekly)
	}

	if v := input.MonthsOfYear; v != nil {
		months := make([]string, 0)
		for _, d := range *v {
			months = append(months, string(d))
		}
		retentionYearlyBlock.Months = months
	}

	if daily := input.RetentionScheduleDaily; daily != nil {
		retentionYearlyBlock.Monthdays = flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(daily)
	}

	return []RetentionYearly{retentionYearlyBlock}
}

func flattenBackupProtectionPolicyVMWorkloadSimpleRetention(input *backup.RetentionDuration) []SimpleRetention {
	if input == nil {
		return nil
	}

	simpleRetentionBlock := SimpleRetention{}

	if v := input.Count; v != nil {
		simpleRetentionBlock.Count = *v
	}

	return []SimpleRetention{simpleRetentionBlock}
}

func expandBackupProtectionPolicyVMWorkloadRetentionDailyFormat(input []int) *backup.DailyRetentionFormat {
	if len(input) == 0 {
		return nil
	}

	daily := backup.DailyRetentionFormat{}

	days := make([]backup.Day, 0)
	for _, item := range input {
		day := backup.Day{
			Date: utils.Int32(int32((item))),
		}

		if item == 0 {
			day.IsLast = utils.Bool(true)
		} else {
			day.IsLast = utils.Bool(false)
		}

		days = append(days, day)
	}
	daily.DaysOfTheMonth = &days

	return &daily
}

func expandBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(weekdays, weeks []string) *backup.WeeklyRetentionFormat {
	if len(weekdays) == 0 && len(weeks) == 0 {
		return nil
	}

	weekly := backup.WeeklyRetentionFormat{}

	if len(weekdays) > 0 {
		weekdaysBlock := make([]backup.DayOfWeek, 0)
		for _, day := range weekdays {
			weekdaysBlock = append(weekdaysBlock, backup.DayOfWeek(day))
		}
		weekly.DaysOfTheWeek = &weekdaysBlock
	}

	if len(weeks) > 0 {
		weeksBlock := make([]backup.WeekOfMonth, 0)
		for _, week := range weeks {
			weeksBlock = append(weeksBlock, backup.WeekOfMonth(week))
		}
		weekly.WeeksOfTheMonth = &weeksBlock
	}

	return &weekly
}

func flattenBackupProtectionPolicyVMWorkloadRetentionWeeklyFormat(input *backup.WeeklyRetentionFormat) (weekdays, weeks []string) {
	if v := input.DaysOfTheWeek; v != nil {
		days := make([]string, 0)
		for _, d := range *v {
			days = append(days, string(d))
		}
		weekdays = days
	}

	if v := input.WeeksOfTheMonth; v != nil {
		days := make([]string, 0)
		for _, d := range *v {
			days = append(days, string(d))
		}
		weeks = days
	}

	return weekdays, weeks
}

func flattenBackupProtectionPolicyVMWorkloadRetentionDailyFormat(input *backup.DailyRetentionFormat) []int {
	result := make([]int, 0)

	if days := input.DaysOfTheMonth; days != nil {
		for _, d := range *days {
			result = append(result, int(*d.Date))
		}
	}

	return result
}
