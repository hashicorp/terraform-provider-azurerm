// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupPolicyPostgreSQLFlexibleServerModel struct {
	Name                         string                                              `tfschema:"name"`
	ResourceGroupName            string                                              `tfschema:"resource_group_name"`
	VaultName                    string                                              `tfschema:"vault_name"`
	BackupRepeatingTimeIntervals []string                                            `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionDuration     string                                              `tfschema:"default_retention_duration"`
	RetentionRule                []BackupPolicyPostgreSQLFlexibleServerRetentionRule `tfschema:"retention_rule"`
	TimeZone                     string                                              `tfschema:"time_zone"`
}

type BackupPolicyPostgreSQLFlexibleServerRetentionRule struct {
	Name     string                                         `tfschema:"name"`
	Criteria []BackupPolicyPostgreSQLFlexibleServerCriteria `tfschema:"criteria"`
	Duration string                                         `tfschema:"duration"`
	Priority int                                            `tfschema:"priority"`
}

type BackupPolicyPostgreSQLFlexibleServerCriteria struct {
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}

type DataProtectionBackupPolicyPostgreSQLFlexibleServerResource struct{}

var _ sdk.Resource = DataProtectionBackupPolicyPostgreSQLFlexibleServerResource{}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_postgresql_flexible_server"
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) ModelObject() interface{} {
	return &BackupPolicyPostgreSQLFlexibleServerModel{}
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backuppolicies.ValidateBackupPolicyID
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"),
				"DataProtection BackupPolicy name must be 3 - 150 characters long, contain only letters, numbers and hyphens.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"vault_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"default_retention_duration": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ISO8601Duration,
		},

		"retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"duration": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.ISO8601Duration,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(backuppolicies.AbsoluteMarkerAllBackup),
										string(backuppolicies.AbsoluteMarkerFirstOfDay),
										string(backuppolicies.AbsoluteMarkerFirstOfMonth),
										string(backuppolicies.AbsoluteMarkerFirstOfWeek),
										string(backuppolicies.AbsoluteMarkerFirstOfYear),
									}, false),
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsDayOfTheWeek(false),
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsMonth(false),
									},
								},

								"scheduled_backup_times": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},

								"weeks_of_month": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											string(backuppolicies.WeekNumberFirst),
											string(backuppolicies.WeekNumberSecond),
											string(backuppolicies.WeekNumberThird),
											string(backuppolicies.WeekNumberFourth),
											string(backuppolicies.WeekNumberLast),
										}, false),
									},
								},
							},
						},
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	return arguments
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupPolicyPostgreSQLFlexibleServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := backuppolicies.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.VaultName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			taggingCriteria, err := expandBackupPolicyPostgreSQLFlexibleServerTaggingCriteriaArray(model.RetentionRule)
			if err != nil {
				return err
			}

			policyRules := make([]backuppolicies.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyPostgreSQLFlexibleServerAzureBackupRuleArray(model.BackupRepeatingTimeIntervals, model.TimeZone, taggingCriteria)...)
			policyRules = append(policyRules, expandBackupPolicyPostgreSQLFlexibleServerDefaultAzureRetentionRule(model.DefaultRetentionDuration))
			policyRules = append(policyRules, expandBackupPolicyPostgreSQLFlexibleServerAzureRetentionRuleArray(model.RetentionRule)...)

			parameters := backuppolicies.BaseBackupPolicyResource{
				Properties: &backuppolicies.BackupPolicy{
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DBforPostgreSQL/flexibleServers/databases"},
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := backuppolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := BackupPolicyPostgreSQLFlexibleServerModel{
				Name:              id.BackupPolicyName,
				ResourceGroupName: id.ResourceGroupName,
				VaultName:         id.BackupVaultName,
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(backuppolicies.BackupPolicy); ok {
					state.DefaultRetentionDuration = flattenBackupPolicyPostgreSQLFlexibleServerDefaultRetentionRuleDuration(&properties.PolicyRules)
					state.RetentionRule = flattenBackupPolicyPostgreSQLFlexibleServerRetentionRuleArray(&properties.PolicyRules)
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyPostgreSQLFlexibleServerBackupRuleArray(&properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyPostgreSQLFlexibleServerBackupTimeZone(&properties.PolicyRules)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupPolicyPostgreSQLFlexibleServerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := backuppolicies.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandBackupPolicyPostgreSQLFlexibleServerAzureBackupRuleArray(input []string, timeZone string, taggingCriteria *[]backuppolicies.TaggingCriteria) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	results = append(results, backuppolicies.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: backuppolicies.DataStoreInfoBase{
			DataStoreType: backuppolicies.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: &backuppolicies.AzureBackupParams{
			BackupType: "Full",
		},
		Trigger: backuppolicies.ScheduleBasedTriggerContext{
			Schedule: backuppolicies.BackupSchedule{
				RepeatingTimeIntervals: input,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: *taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyPostgreSQLFlexibleServerAzureRetentionRuleArray(input []BackupPolicyPostgreSQLFlexibleServerRetentionRule) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	for _, item := range input {
		results = append(results, backuppolicies.AzureRetentionRule{
			Name:      item.Name,
			IsDefault: utils.Bool(false),
			Lifecycles: []backuppolicies.SourceLifeCycle{
				{
					DeleteAfter: backuppolicies.AbsoluteDeleteOption{
						Duration: item.Duration,
					},
					SourceDataStore: backuppolicies.DataStoreInfoBase{
						DataStoreType: "VaultStore",
						ObjectType:    "DataStoreInfoBase",
					},
					TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
				},
			},
		})
	}
	return results
}

func expandBackupPolicyPostgreSQLFlexibleServerDefaultAzureRetentionRule(input interface{}) backuppolicies.BasePolicyRule {
	return backuppolicies.AzureRetentionRule{
		Name:      "Default",
		IsDefault: utils.Bool(true),
		Lifecycles: []backuppolicies.SourceLifeCycle{
			{
				DeleteAfter: backuppolicies.AbsoluteDeleteOption{
					Duration: input.(string),
				},
				SourceDataStore: backuppolicies.DataStoreInfoBase{
					DataStoreType: "VaultStore",
					ObjectType:    "DataStoreInfoBase",
				},
				TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyPostgreSQLFlexibleServerTaggingCriteriaArray(input []BackupPolicyPostgreSQLFlexibleServerRetentionRule) (*[]backuppolicies.TaggingCriteria, error) {
	results := []backuppolicies.TaggingCriteria{
		{
			Criteria:        nil,
			IsDefault:       true,
			TaggingPriority: 99,
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String("Default_"),
				TagName: "Default",
			},
		},
	}
	for _, item := range input {
		result := backuppolicies.TaggingCriteria{
			IsDefault:       false,
			TaggingPriority: int64(item.Priority),
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String(item.Name + "_"),
				TagName: item.Name,
			},
		}

		criteria, err := expandBackupPolicyPostgreSQLFlexibleServerCriteriaArray(item.Criteria)
		if err != nil {
			return nil, err
		}
		result.Criteria = criteria

		results = append(results, result)
	}
	return &results, nil
}

func expandBackupPolicyPostgreSQLFlexibleServerCriteriaArray(input []BackupPolicyPostgreSQLFlexibleServerCriteria) (*[]backuppolicies.BackupCriteria, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("criteria is a required field, cannot leave blank")
	}
	results := make([]backuppolicies.BackupCriteria, 0)

	for _, item := range input {
		var absoluteCriteria []backuppolicies.AbsoluteMarker
		if absoluteCriteriaRaw := item.AbsoluteCriteria; len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []backuppolicies.AbsoluteMarker{backuppolicies.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfWeek []backuppolicies.DayOfWeek
		if len(item.DaysOfWeek) > 0 {
			daysOfWeek = make([]backuppolicies.DayOfWeek, 0)
			for _, value := range item.DaysOfWeek {
				daysOfWeek = append(daysOfWeek, backuppolicies.DayOfWeek(value))
			}
		}

		var monthsOfYear []backuppolicies.Month
		if len(item.MonthsOfYear) > 0 {
			monthsOfYear = make([]backuppolicies.Month, 0)
			for _, value := range item.MonthsOfYear {
				monthsOfYear = append(monthsOfYear, backuppolicies.Month(value))
			}
		}

		var weeksOfMonth []backuppolicies.WeekNumber
		if len(item.WeeksOfMonth) > 0 {
			weeksOfMonth = make([]backuppolicies.WeekNumber, 0)
			for _, value := range item.WeeksOfMonth {
				weeksOfMonth = append(weeksOfMonth, backuppolicies.WeekNumber(value))
			}
		}

		var scheduleTimes *[]string
		if len(item.ScheduledBackupTimes) > 0 {
			scheduleTimes = pointer.To(item.ScheduledBackupTimes)
		}
		results = append(results, backuppolicies.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			DaysOfMonth:      nil,
			DaysOfTheWeek:    &daysOfWeek,
			MonthsOfYear:     &monthsOfYear,
			ScheduleTimes:    scheduleTimes,
			WeeksOfTheMonth:  &weeksOfMonth,
		})
	}
	return &results, nil
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupRuleArray(input *[]backuppolicies.BasePolicyRule) []string {
	if input == nil {
		return make([]string, 0)
	}
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					return scheduleBasedTrigger.Schedule.RepeatingTimeIntervals
				}
			}
		}
	}
	return make([]string, 0)
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupTimeZone(input *[]backuppolicies.BasePolicyRule) string {
	if input == nil {
		return ""
	}
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					return pointer.From(scheduleBasedTrigger.Schedule.TimeZone)
				}
			}
		}
	}
	return ""
}

func flattenBackupPolicyPostgreSQLFlexibleServerDefaultRetentionRuleDuration(input *[]backuppolicies.BasePolicyRule) string {
	if input == nil {
		return ""
	}

	for _, item := range *input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if retentionRule.Lifecycles != nil && len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
					return deleteOption.Duration
				}
			}
		}
	}
	return ""
}

func flattenBackupPolicyPostgreSQLFlexibleServerRetentionRuleArray(input *[]backuppolicies.BasePolicyRule) []BackupPolicyPostgreSQLFlexibleServerRetentionRule {
	results := make([]BackupPolicyPostgreSQLFlexibleServerRetentionRule, 0)
	if input == nil {
		return results
	}

	var taggingCriterias []backuppolicies.TaggingCriteria
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriterias = trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range *input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok && (retentionRule.IsDefault == nil || !*retentionRule.IsDefault) {
			name := retentionRule.Name
			var taggingPriority int
			var taggingCriteria []BackupPolicyPostgreSQLFlexibleServerCriteria
			for _, criteria := range taggingCriterias {
				if strings.EqualFold(criteria.TagInfo.TagName, name) {
					taggingPriority = int(criteria.TaggingPriority)
					taggingCriteria = flattenBackupPolicyPostgreSQLFlexibleServerBackupCriteriaArray(criteria.Criteria)
				}
			}
			var duration string
			if retentionRule.Lifecycles != nil && len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
					duration = deleteOption.Duration
				}
			}
			results = append(results, BackupPolicyPostgreSQLFlexibleServerRetentionRule{
				Name:     name,
				Priority: taggingPriority,
				Criteria: taggingCriteria,
				Duration: duration,
			})
		}
	}
	return results
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupCriteriaArray(input *[]backuppolicies.BackupCriteria) []BackupPolicyPostgreSQLFlexibleServerCriteria {
	results := make([]BackupPolicyPostgreSQLFlexibleServerCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if criteria, ok := item.(backuppolicies.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(*criteria.AbsoluteCriteria) > 0 {
				absoluteCriteria = string((*criteria.AbsoluteCriteria)[0])
			}
			var daysOfWeek []string
			if criteria.DaysOfTheWeek != nil {
				daysOfWeek = make([]string, 0)
				for _, item := range *criteria.DaysOfTheWeek {
					daysOfWeek = append(daysOfWeek, (string)(item))
				}
			}
			var monthsOfYear []string
			if criteria.MonthsOfYear != nil {
				monthsOfYear = make([]string, 0)
				for _, item := range *criteria.MonthsOfYear {
					monthsOfYear = append(monthsOfYear, (string)(item))
				}
			}
			var weeksOfMonth []string
			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth = make([]string, 0)
				for _, item := range *criteria.WeeksOfTheMonth {
					weeksOfMonth = append(weeksOfMonth, (string)(item))
				}
			}
			var scheduleTimes []string
			if criteria.ScheduleTimes != nil {
				scheduleTimes = make([]string, 0)
				scheduleTimes = append(scheduleTimes, *criteria.ScheduleTimes...)
			}

			results = append(results, BackupPolicyPostgreSQLFlexibleServerCriteria{
				AbsoluteCriteria:     absoluteCriteria,
				DaysOfWeek:           daysOfWeek,
				MonthsOfYear:         monthsOfYear,
				WeeksOfMonth:         weeksOfMonth,
				ScheduledBackupTimes: scheduleTimes,
			})
		}
	}
	return results
}
