// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type BackupPolicyPostgreSQLFlexibleServerModel struct {
	Name                         string                                                     `tfschema:"name"`
	VaultId                      string                                                     `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                                   `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRules               []BackupPolicyPostgreSQLFlexibleServerRetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                                                     `tfschema:"time_zone"`
}

type BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule struct {
	LifeCycle []BackupPolicyPostgreSQLFlexibleServerLifeCycle `tfschema:"life_cycle"`
}

type BackupPolicyPostgreSQLFlexibleServerLifeCycle struct {
	DataStoreType string `tfschema:"data_store_type"`
	Duration      string `tfschema:"duration"`
}

type BackupPolicyPostgreSQLFlexibleServerRetentionRule struct {
	Name      string                                          `tfschema:"name"`
	Criteria  []BackupPolicyPostgreSQLFlexibleServerCriteria  `tfschema:"criteria"`
	LifeCycle []BackupPolicyPostgreSQLFlexibleServerLifeCycle `tfschema:"life_cycle"`
	Priority  int64                                           `tfschema:"priority"`
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyPostgreSQLFlexibleServerName,
		},

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(backuppolicies.BackupVaultId{})),

		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"default_retention_rule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										// Confirmed with the service team that current possible value only support `VaultStore`.
										// However, considering that `ArchiveStore` will be supported in the future, it would be exposed for user specification.
										string(backuppolicies.DataStoreTypesVaultStore),
									}, false),
								},

								"duration": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: azValidate.ISO8601Duration,
								},
							},
						},
					},
				},
			},
		},

		"retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(backuppolicies.PossibleValuesForAbsoluteMarker(), false),
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(backuppolicies.PossibleValuesForDayOfWeek(), false),
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(backuppolicies.PossibleValuesForMonth(), false),
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
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(backuppolicies.PossibleValuesForWeekNumber(), false),
									},
								},
							},
						},
					},

					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ForceNew: true,
									ValidateFunc: validation.StringInSlice([]string{
										// Confirmed with the service team that currently only `VaultStore` is supported.
										// However, considering that `ArchiveStore` will be supported in the future, it would be exposed for user specification.
										string(backuppolicies.DataStoreTypesVaultStore),
									}, false),
								},

								"duration": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: azValidate.ISO8601Duration,
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

			vaultId, _ := backuppolicies.ParseBackupVaultID(model.VaultId)
			id := backuppolicies.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			policyRules := make([]backuppolicies.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyPostgreSQLFlexibleServerAzureBackupRules(model.BackupRepeatingTimeIntervals, model.TimeZone, expandBackupPolicyPostgreSQLFlexibleServerTaggingCriteria(model.RetentionRules))...)
			policyRules = append(policyRules, expandBackupPolicyPostgreSQLFlexibleServerDefaultAzureRetentionRule(model.DefaultRetentionRule))
			policyRules = append(policyRules, expandBackupPolicyPostgreSQLFlexibleServerAzureRetentionRules(model.RetentionRules)...)

			parameters := backuppolicies.BaseBackupPolicyResource{
				Properties: backuppolicies.BackupPolicy{
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DBforPostgreSQL/flexibleServers"},
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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

			vaultId := backuppolicies.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
			state := BackupPolicyPostgreSQLFlexibleServerModel{
				Name:    id.BackupPolicyName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(backuppolicies.BackupPolicy); ok {
					state.DefaultRetentionRule = flattenBackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule(properties.PolicyRules)
					state.RetentionRules = flattenBackupPolicyPostgreSQLFlexibleServerRetentionRules(properties.PolicyRules)
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyPostgreSQLFlexibleServerBackupRules(properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyPostgreSQLFlexibleServerBackupTimeZone(properties.PolicyRules)
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

func expandBackupPolicyPostgreSQLFlexibleServerAzureBackupRules(input []string, timeZone string, taggingCriteria []backuppolicies.TaggingCriteria) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)

	results = append(results, backuppolicies.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: backuppolicies.DataStoreInfoBase{
			DataStoreType: backuppolicies.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: backuppolicies.AzureBackupParams{
			BackupType: "Full",
		},
		Trigger: backuppolicies.ScheduleBasedTriggerContext{
			Schedule: backuppolicies.BackupSchedule{
				RepeatingTimeIntervals: input,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyPostgreSQLFlexibleServerAzureRetentionRules(input []BackupPolicyPostgreSQLFlexibleServerRetentionRule) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)

	for _, item := range input {
		results = append(results, backuppolicies.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  utils.Bool(false),
			Lifecycles: expandBackupPolicyPostgreSQLFlexibleServerLifeCycle(item.LifeCycle),
		})
	}

	return results
}

func expandBackupPolicyPostgreSQLFlexibleServerDefaultAzureRetentionRule(input []BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule) backuppolicies.BasePolicyRule {
	result := backuppolicies.AzureRetentionRule{
		Name:      "Default",
		IsDefault: utils.Bool(true),
	}

	if len(input) > 0 {
		result.Lifecycles = expandBackupPolicyPostgreSQLFlexibleServerLifeCycle(input[0].LifeCycle)
	}

	return result
}

func expandBackupPolicyPostgreSQLFlexibleServerLifeCycle(input []BackupPolicyPostgreSQLFlexibleServerLifeCycle) []backuppolicies.SourceLifeCycle {
	results := make([]backuppolicies.SourceLifeCycle, 0)

	for _, item := range input {
		sourceLifeCycle := backuppolicies.SourceLifeCycle{
			DeleteAfter: backuppolicies.AbsoluteDeleteOption{
				Duration: item.Duration,
			},
			SourceDataStore: backuppolicies.DataStoreInfoBase{
				DataStoreType: backuppolicies.DataStoreTypes(item.DataStoreType),
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
		}

		results = append(results, sourceLifeCycle)
	}

	return results
}

func expandBackupPolicyPostgreSQLFlexibleServerTaggingCriteria(input []BackupPolicyPostgreSQLFlexibleServerRetentionRule) []backuppolicies.TaggingCriteria {
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
			Criteria:        expandBackupPolicyPostgreSQLFlexibleServerCriteria(item.Criteria),
			TaggingPriority: item.Priority,
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String(item.Name + "_"),
				TagName: item.Name,
			},
		}

		results = append(results, result)
	}

	return results
}

func expandBackupPolicyPostgreSQLFlexibleServerCriteria(input []BackupPolicyPostgreSQLFlexibleServerCriteria) *[]backuppolicies.BackupCriteria {
	if len(input) == 0 {
		return nil
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

		var scheduleTimes []string
		if len(item.ScheduledBackupTimes) > 0 {
			scheduleTimes = item.ScheduledBackupTimes
		}

		results = append(results, backuppolicies.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To(absoluteCriteria),
			DaysOfMonth:      nil,
			DaysOfTheWeek:    pointer.To(daysOfWeek),
			MonthsOfYear:     pointer.To(monthsOfYear),
			ScheduleTimes:    pointer.To(scheduleTimes),
			WeeksOfTheMonth:  pointer.To(weeksOfMonth),
		})
	}

	return &results
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupRules(input []backuppolicies.BasePolicyRule) []string {
	backupRules := make([]string, 0)

	for _, item := range input {
		if v, ok := item.(backuppolicies.AzureBackupRule); ok {
			if v.Trigger != nil {
				if scheduleBasedTrigger, ok := v.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					backupRules = scheduleBasedTrigger.Schedule.RepeatingTimeIntervals
					return backupRules
				}
			}
		}
	}

	return backupRules
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupTimeZone(input []backuppolicies.BasePolicyRule) string {
	var timeZone string

	for _, item := range input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					timeZone = pointer.From(scheduleBasedTrigger.Schedule.TimeZone)
					return timeZone
				}
			}
		}
	}

	return timeZone
}

func flattenBackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule(input []backuppolicies.BasePolicyRule) []BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule {
	results := make([]BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule, 0)

	for _, item := range input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok {
			if pointer.From(retentionRule.IsDefault) {
				var lifeCycle []BackupPolicyPostgreSQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyPostgreSQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyPostgreSQLFlexibleServerDefaultRetentionRule{
					LifeCycle: lifeCycle,
				})
			}
		}
	}

	return results
}

func flattenBackupPolicyPostgreSQLFlexibleServerRetentionRules(input []backuppolicies.BasePolicyRule) []BackupPolicyPostgreSQLFlexibleServerRetentionRule {
	results := make([]BackupPolicyPostgreSQLFlexibleServerRetentionRule, 0)
	var taggingCriterias []backuppolicies.TaggingCriteria

	for _, item := range input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriterias = trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok {
			var name string
			var taggingPriority int64
			var taggingCriteria []BackupPolicyPostgreSQLFlexibleServerCriteria

			if !pointer.From(retentionRule.IsDefault) {
				name = retentionRule.Name

				for _, criteria := range taggingCriterias {
					if strings.EqualFold(criteria.TagInfo.TagName, name) {
						taggingPriority = criteria.TaggingPriority
						taggingCriteria = flattenBackupPolicyPostgreSQLFlexibleServerBackupCriteria(criteria.Criteria)
						break
					}
				}

				var lifeCycle []BackupPolicyPostgreSQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyPostgreSQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyPostgreSQLFlexibleServerRetentionRule{
					Name:      name,
					Priority:  taggingPriority,
					Criteria:  taggingCriteria,
					LifeCycle: lifeCycle,
				})
			}
		}
	}

	return results
}

func flattenBackupPolicyPostgreSQLFlexibleServerLifeCycles(input []backuppolicies.SourceLifeCycle) []BackupPolicyPostgreSQLFlexibleServerLifeCycle {
	results := make([]BackupPolicyPostgreSQLFlexibleServerLifeCycle, 0)

	for _, item := range input {
		var duration string
		var dataStoreType string

		if deleteOption, ok := item.DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
			duration = deleteOption.Duration
		}

		dataStoreType = string(item.SourceDataStore.DataStoreType)

		results = append(results, BackupPolicyPostgreSQLFlexibleServerLifeCycle{
			Duration:      duration,
			DataStoreType: dataStoreType,
		})
	}

	return results
}

func flattenBackupPolicyPostgreSQLFlexibleServerBackupCriteria(input *[]backuppolicies.BackupCriteria) []BackupPolicyPostgreSQLFlexibleServerCriteria {
	results := make([]BackupPolicyPostgreSQLFlexibleServerCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(backuppolicies.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				absoluteCriteria = string((pointer.From(criteria.AbsoluteCriteria))[0])
			}

			var daysOfWeek []string
			if criteria.DaysOfTheWeek != nil {
				daysOfWeek = make([]string, 0)

				for _, item := range pointer.From(criteria.DaysOfTheWeek) {
					daysOfWeek = append(daysOfWeek, (string)(item))
				}
			}

			var monthsOfYear []string
			if criteria.MonthsOfYear != nil {
				monthsOfYear = make([]string, 0)

				for _, item := range pointer.From(criteria.MonthsOfYear) {
					monthsOfYear = append(monthsOfYear, (string)(item))
				}
			}

			var weeksOfMonth []string
			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth = make([]string, 0)

				for _, item := range pointer.From(criteria.WeeksOfTheMonth) {
					weeksOfMonth = append(weeksOfMonth, (string)(item))
				}
			}

			var scheduleTimes []string
			if criteria.ScheduleTimes != nil {
				scheduleTimes = make([]string, 0)
				scheduleTimes = append(scheduleTimes, pointer.From(criteria.ScheduleTimes)...)
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
