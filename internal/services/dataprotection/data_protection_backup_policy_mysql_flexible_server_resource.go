// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/basebackuppolicyresources"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_policy_mysql_flexible_server -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

type BackupPolicyMySQLFlexibleServerModel struct {
	Name                         string                                                `tfschema:"name"`
	VaultId                      string                                                `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals []string                                              `tfschema:"backup_repeating_time_intervals"`
	DefaultRetentionRule         []BackupPolicyMySQLFlexibleServerDefaultRetentionRule `tfschema:"default_retention_rule"`
	RetentionRules               []BackupPolicyMySQLFlexibleServerRetentionRule        `tfschema:"retention_rule"`
	TimeZone                     string                                                `tfschema:"time_zone"`
}

type BackupPolicyMySQLFlexibleServerDefaultRetentionRule struct {
	LifeCycle []BackupPolicyMySQLFlexibleServerLifeCycle `tfschema:"life_cycle"`
}

type BackupPolicyMySQLFlexibleServerLifeCycle struct {
	DataStoreType string `tfschema:"data_store_type"`
	Duration      string `tfschema:"duration"`
}

type BackupPolicyMySQLFlexibleServerRetentionRule struct {
	Name      string                                     `tfschema:"name"`
	Criteria  []BackupPolicyMySQLFlexibleServerCriteria  `tfschema:"criteria"`
	LifeCycle []BackupPolicyMySQLFlexibleServerLifeCycle `tfschema:"life_cycle"`
	Priority  int64                                      `tfschema:"priority"`
}

type BackupPolicyMySQLFlexibleServerCriteria struct {
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}

type DataProtectionBackupPolicyMySQLFlexibleServerResource struct{}

var (
	_ sdk.Resource             = DataProtectionBackupPolicyMySQLFlexibleServerResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupPolicyMySQLFlexibleServerResource{}
)

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Identity() resourceids.ResourceId {
	return &basebackuppolicyresources.BackupPolicyId{}
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_mysql_flexible_server"
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) ModelObject() interface{} {
	return &BackupPolicyMySQLFlexibleServerModel{}
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return basebackuppolicyresources.ValidateBackupPolicyID
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyMySQLFlexibleServerName,
		},

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(basebackuppolicyresources.BackupVaultId{})),

		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azValidate.ISO8601RepeatingTime,
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
										string(basebackuppolicyresources.DataStoreTypesVaultStore),
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
									ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForAbsoluteMarker(), false),
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForDayOfWeek(), false),
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									ForceNew: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForMonth(), false),
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
										ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForWeekNumber(), false),
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
										string(basebackuppolicyresources.DataStoreTypesVaultStore),
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
			ValidateFunc: validate.BackupPolicyMySQLFlexibleServerTimeZone(),
		},
	}
	return arguments
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupPolicyMySQLFlexibleServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vaultId, _ := basebackuppolicyresources.ParseBackupVaultID(model.VaultId)
			id := basebackuppolicyresources.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			existing, err := client.BackupPoliciesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerAzureBackupRules(model.BackupRepeatingTimeIntervals, model.TimeZone, expandBackupPolicyMySQLFlexibleServerTaggingCriteria(model.RetentionRules))...)
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerDefaultAzureRetentionRule(model.DefaultRetentionRule))
			policyRules = append(policyRules, expandBackupPolicyMySQLFlexibleServerAzureRetentionRules(model.RetentionRules)...)

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: basebackuppolicyresources.BackupPolicy{
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DBforMySQL/flexibleServers"},
				},
			}

			if _, err := client.BackupPoliciesCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := basebackuppolicyresources.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.BackupPoliciesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			vaultId := basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
			state := BackupPolicyMySQLFlexibleServerModel{
				Name:    id.BackupPolicyName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.DefaultRetentionRule = flattenBackupPolicyMySQLFlexibleServerDefaultRetentionRule(properties.PolicyRules)
					state.RetentionRules = flattenBackupPolicyMySQLFlexibleServerRetentionRules(properties.PolicyRules)
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyMySQLFlexibleServerBackupRules(properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyMySQLFlexibleServerBackupTimeZone(properties.PolicyRules)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupPolicyMySQLFlexibleServerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			id, err := basebackuppolicyresources.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.BackupPoliciesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandBackupPolicyMySQLFlexibleServerAzureBackupRules(input []string, timeZone string, taggingCriteria []basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "Full",
		},
		Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
			Schedule: basebackuppolicyresources.BackupSchedule{
				RepeatingTimeIntervals: input,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyMySQLFlexibleServerAzureRetentionRules(input []BackupPolicyMySQLFlexibleServerRetentionRule) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	for _, item := range input {
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyMySQLFlexibleServerLifeCycle(item.LifeCycle),
		})
	}

	return results
}

func expandBackupPolicyMySQLFlexibleServerDefaultAzureRetentionRule(input []BackupPolicyMySQLFlexibleServerDefaultRetentionRule) basebackuppolicyresources.BasePolicyRule {
	result := basebackuppolicyresources.AzureRetentionRule{
		Name:      "Default",
		IsDefault: pointer.To(true),
	}

	if len(input) > 0 {
		result.Lifecycles = expandBackupPolicyMySQLFlexibleServerLifeCycle(input[0].LifeCycle)
	}

	return result
}

func expandBackupPolicyMySQLFlexibleServerLifeCycle(input []BackupPolicyMySQLFlexibleServerLifeCycle) []basebackuppolicyresources.SourceLifeCycle {
	results := make([]basebackuppolicyresources.SourceLifeCycle, 0)

	for _, item := range input {
		sourceLifeCycle := basebackuppolicyresources.SourceLifeCycle{
			DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
				Duration: item.Duration,
			},
			SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypes(item.DataStoreType),
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
		}

		results = append(results, sourceLifeCycle)
	}

	return results
}

func expandBackupPolicyMySQLFlexibleServerTaggingCriteria(input []BackupPolicyMySQLFlexibleServerRetentionRule) []basebackuppolicyresources.TaggingCriteria {
	results := []basebackuppolicyresources.TaggingCriteria{
		{
			Criteria:        nil,
			IsDefault:       true,
			TaggingPriority: 99,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To("Default_"),
				TagName: "Default",
			},
		},
	}

	for _, item := range input {
		result := basebackuppolicyresources.TaggingCriteria{
			IsDefault:       false,
			Criteria:        expandBackupPolicyMySQLFlexibleServerCriteria(item.Criteria),
			TaggingPriority: item.Priority,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(item.Name + "_"),
				TagName: item.Name,
			},
		}

		results = append(results, result)
	}

	return results
}

func expandBackupPolicyMySQLFlexibleServerCriteria(input []BackupPolicyMySQLFlexibleServerCriteria) *[]basebackuppolicyresources.BackupCriteria {
	if len(input) == 0 {
		return nil
	}

	results := make([]basebackuppolicyresources.BackupCriteria, 0)

	for _, item := range input {
		var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
		if absoluteCriteriaRaw := item.AbsoluteCriteria; len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfWeek []basebackuppolicyresources.DayOfWeek
		if len(item.DaysOfWeek) > 0 {
			daysOfWeek = make([]basebackuppolicyresources.DayOfWeek, 0)
			for _, value := range item.DaysOfWeek {
				daysOfWeek = append(daysOfWeek, basebackuppolicyresources.DayOfWeek(value))
			}
		}

		var monthsOfYear []basebackuppolicyresources.Month
		if len(item.MonthsOfYear) > 0 {
			monthsOfYear = make([]basebackuppolicyresources.Month, 0)
			for _, value := range item.MonthsOfYear {
				monthsOfYear = append(monthsOfYear, basebackuppolicyresources.Month(value))
			}
		}

		var weeksOfMonth []basebackuppolicyresources.WeekNumber
		if len(item.WeeksOfMonth) > 0 {
			weeksOfMonth = make([]basebackuppolicyresources.WeekNumber, 0)
			for _, value := range item.WeeksOfMonth {
				weeksOfMonth = append(weeksOfMonth, basebackuppolicyresources.WeekNumber(value))
			}
		}

		var scheduleTimes []string
		if len(item.ScheduledBackupTimes) > 0 {
			scheduleTimes = item.ScheduledBackupTimes
		}

		results = append(results, basebackuppolicyresources.ScheduleBasedBackupCriteria{
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

func flattenBackupPolicyMySQLFlexibleServerBackupRules(input []basebackuppolicyresources.BasePolicyRule) []string {
	backupRules := make([]string, 0)

	for _, item := range input {
		if v, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if v.Trigger != nil {
				if scheduleBasedTrigger, ok := v.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					backupRules = scheduleBasedTrigger.Schedule.RepeatingTimeIntervals
					return backupRules
				}
			}
		}
	}

	return backupRules
}

func flattenBackupPolicyMySQLFlexibleServerBackupTimeZone(input []basebackuppolicyresources.BasePolicyRule) string {
	var timeZone string

	for _, item := range input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
					timeZone = pointer.From(scheduleBasedTrigger.Schedule.TimeZone)
					return timeZone
				}
			}
		}
	}

	return timeZone
}

func flattenBackupPolicyMySQLFlexibleServerDefaultRetentionRule(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyMySQLFlexibleServerDefaultRetentionRule {
	results := make([]BackupPolicyMySQLFlexibleServerDefaultRetentionRule, 0)

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			if pointer.From(retentionRule.IsDefault) {
				var lifeCycle []BackupPolicyMySQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyMySQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyMySQLFlexibleServerDefaultRetentionRule{
					LifeCycle: lifeCycle,
				})
			}
		}
	}

	return results
}

func flattenBackupPolicyMySQLFlexibleServerRetentionRules(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyMySQLFlexibleServerRetentionRule {
	results := make([]BackupPolicyMySQLFlexibleServerRetentionRule, 0)
	var taggingCriterias []basebackuppolicyresources.TaggingCriteria

	for _, item := range input {
		if backupRule, ok := item.(basebackuppolicyresources.AzureBackupRule); ok {
			if trigger, ok := backupRule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriterias = trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok {
			var name string
			var taggingPriority int64
			var taggingCriteria []BackupPolicyMySQLFlexibleServerCriteria

			if !pointer.From(retentionRule.IsDefault) {
				name = retentionRule.Name

				for _, criteria := range taggingCriterias {
					if strings.EqualFold(criteria.TagInfo.TagName, name) {
						taggingPriority = criteria.TaggingPriority
						taggingCriteria = flattenBackupPolicyMySQLFlexibleServerBackupCriteria(criteria.Criteria)
						break
					}
				}

				var lifeCycle []BackupPolicyMySQLFlexibleServerLifeCycle
				if v := retentionRule.Lifecycles; len(v) > 0 {
					lifeCycle = flattenBackupPolicyMySQLFlexibleServerLifeCycles(v)
				}

				results = append(results, BackupPolicyMySQLFlexibleServerRetentionRule{
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

func flattenBackupPolicyMySQLFlexibleServerLifeCycles(input []basebackuppolicyresources.SourceLifeCycle) []BackupPolicyMySQLFlexibleServerLifeCycle {
	results := make([]BackupPolicyMySQLFlexibleServerLifeCycle, 0)

	for _, item := range input {
		var duration string
		var dataStoreType string

		if deleteOption, ok := item.DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
			duration = deleteOption.Duration
		}

		dataStoreType = string(item.SourceDataStore.DataStoreType)

		results = append(results, BackupPolicyMySQLFlexibleServerLifeCycle{
			Duration:      duration,
			DataStoreType: dataStoreType,
		})
	}

	return results
}

func flattenBackupPolicyMySQLFlexibleServerBackupCriteria(input *[]basebackuppolicyresources.BackupCriteria) []BackupPolicyMySQLFlexibleServerCriteria {
	results := make([]BackupPolicyMySQLFlexibleServerCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			var absoluteCriteria string
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				absoluteCriteria = string((pointer.From(criteria.AbsoluteCriteria))[0])
			}

			daysOfWeek := make([]string, 0)
			if criteria.DaysOfTheWeek != nil {
				for _, item := range pointer.From(criteria.DaysOfTheWeek) {
					daysOfWeek = append(daysOfWeek, (string)(item))
				}
			}

			monthsOfYear := make([]string, 0)
			if criteria.MonthsOfYear != nil {
				for _, item := range pointer.From(criteria.MonthsOfYear) {
					monthsOfYear = append(monthsOfYear, (string)(item))
				}
			}

			weeksOfMonth := make([]string, 0)
			if criteria.WeeksOfTheMonth != nil {
				for _, item := range pointer.From(criteria.WeeksOfTheMonth) {
					weeksOfMonth = append(weeksOfMonth, (string)(item))
				}
			}

			scheduleTimes := make([]string, 0)
			if criteria.ScheduleTimes != nil {
				scheduleTimes = append(scheduleTimes, pointer.From(criteria.ScheduleTimes)...)
			}

			results = append(results, BackupPolicyMySQLFlexibleServerCriteria{
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
