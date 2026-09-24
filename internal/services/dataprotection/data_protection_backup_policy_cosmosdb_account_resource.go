// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_policy_cosmosdb_account -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type BackupPolicyCosmosdbAccountModel struct {
	Name                        string                                     `tfschema:"name"`
	DataProtectionBackupVaultId string                                     `tfschema:"data_protection_backup_vault_id"`
	BackupSchedule              []string                                   `tfschema:"backup_schedule"`
	DefaultRetentionDuration    string                                     `tfschema:"default_retention_duration"`
	RetentionRules              []BackupPolicyCosmosdbAccountRetentionRule `tfschema:"retention_rule"`
	TimeZone                    string                                     `tfschema:"time_zone"`
}

type BackupPolicyCosmosdbAccountRetentionRule struct {
	Name                 string   `tfschema:"name"`
	Duration             string   `tfschema:"duration"`
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}

type DataProtectionBackupPolicyCosmosdbAccountResource struct{}

var (
	_ sdk.Resource             = DataProtectionBackupPolicyCosmosdbAccountResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupPolicyCosmosdbAccountResource{}
)

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Identity() resourceids.ResourceId {
	return &basebackuppolicyresources.BackupPolicyId{}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_cosmosdb_account"
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) ModelObject() interface{} {
	return &BackupPolicyCosmosdbAccountModel{}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return basebackuppolicyresources.ValidateBackupPolicyID
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{2,149}$"),
				"`name` must be 3 - 150 characters long, contain only letters, numbers and hyphens, and cannot start with a number or hyphen",
			),
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(basebackuppolicyresources.BackupVaultId{})),

		// check maxItems
		"backup_schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.ISO8601RepeatingTime,
			},
		},

		"default_retention_duration": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
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

					"duration": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.ISO8601Duration,
					},

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

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.BackupPolicyCosmosdbAccountTimeZone(),
		},
	}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient20260601
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackupPolicyCosmosdbAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			for _, rule := range model.RetentionRules {
				if rule.AbsoluteCriteria == "" && len(rule.DaysOfWeek) == 0 {
					return fmt.Errorf("`retention_rule` %q requires at least one of `absolute_criteria` and `days_of_week` to be specified", rule.Name)
				}
			}

			vaultId, err := basebackuppolicyresources.ParseBackupVaultID(model.DataProtectionBackupVaultId)
			if err != nil {
				return err
			}
			id := basebackuppolicyresources.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupPoliciesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountRetentionRules(model.RetentionRules)...)
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountDefaultRetentionRule(model.DefaultRetentionDuration))
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountBackupRules(model.BackupSchedule, model.TimeZone, expandBackupPolicyCosmosdbAccountTaggingCriteria(model.RetentionRules))...)

			parameters := basebackuppolicyresources.BaseBackupPolicyResource{
				Properties: &basebackuppolicyresources.BackupPolicy{
					ObjectType:      "BackupPolicy",
					PolicyRules:     policyRules,
					DatasourceTypes: []string{"Microsoft.DocumentDB/databaseAccounts"},
				},
			}
			if _, err := client.BackupPoliciesCreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient20260601
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
			state := BackupPolicyCosmosdbAccountModel{
				Name:                        id.BackupPolicyName,
				DataProtectionBackupVaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties, ok := model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.DefaultRetentionDuration, state.RetentionRules, state.BackupSchedule, state.TimeZone = flattenBackupPolicyCosmosdbAccountPolicyRules(properties.PolicyRules)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient20260601
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

func expandBackupPolicyCosmosdbAccountRetentionRules(input []BackupPolicyCosmosdbAccountRetentionRule) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	for _, item := range input {
		results = append(results, basebackuppolicyresources.AzureRetentionRule{
			Name:       item.Name,
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyCosmosdbAccountLifeCycle(item.Duration),
		})
	}

	return results
}

func expandBackupPolicyCosmosdbAccountDefaultRetentionRule(duration string) basebackuppolicyresources.BasePolicyRule {
	return basebackuppolicyresources.AzureRetentionRule{
		Name:       "Default",
		IsDefault:  pointer.To(true),
		Lifecycles: expandBackupPolicyCosmosdbAccountLifeCycle(duration),
	}
}

func expandBackupPolicyCosmosdbAccountBackupRules(input []string, timeZone string, taggingCriteria []basebackuppolicyresources.TaggingCriteria) []basebackuppolicyresources.BasePolicyRule {
	results := make([]basebackuppolicyresources.BasePolicyRule, 0)

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupRule",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "full",
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

func expandBackupPolicyCosmosdbAccountLifeCycle(duration string) []basebackuppolicyresources.SourceLifeCycle {
	// NOTE: currently only `VaultStore` is supported by the service team. When other options are supported
	// in the future, export `data_store_type` as a schema field and use `VaultStore` as the default value.
	return []basebackuppolicyresources.SourceLifeCycle{
		{
			DeleteAfter: basebackuppolicyresources.AbsoluteDeleteOption{
				Duration: duration,
			},
			SourceDataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]basebackuppolicyresources.TargetCopySetting{},
		},
	}
}

func expandBackupPolicyCosmosdbAccountTaggingCriteria(input []BackupPolicyCosmosdbAccountRetentionRule) []basebackuppolicyresources.TaggingCriteria {
	results := []basebackuppolicyresources.TaggingCriteria{
		{
			IsDefault:       true,
			TaggingPriority: 99,
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To("Default_"),
				TagName: "Default",
			},
		},
	}

	for i, item := range input {
		result := basebackuppolicyresources.TaggingCriteria{
			Criteria:        expandBackupPolicyCosmosdbAccountRetentionRuleCriteria(item),
			TaggingPriority: int64(i + 1),
			TagInfo: basebackuppolicyresources.RetentionTag{
				Id:      pointer.To(item.Name + "_"),
				TagName: item.Name,
			},
		}

		results = append(results, result)
	}

	return results
}

func expandBackupPolicyCosmosdbAccountRetentionRuleCriteria(input BackupPolicyCosmosdbAccountRetentionRule) *[]basebackuppolicyresources.BackupCriteria {
	var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
	if len(input.AbsoluteCriteria) > 0 {
		absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(input.AbsoluteCriteria)}
	}

	var daysOfWeek []basebackuppolicyresources.DayOfWeek
	if len(input.DaysOfWeek) > 0 {
		daysOfWeek = make([]basebackuppolicyresources.DayOfWeek, 0)
		for _, value := range input.DaysOfWeek {
			daysOfWeek = append(daysOfWeek, basebackuppolicyresources.DayOfWeek(value))
		}
	}

	var monthsOfYear []basebackuppolicyresources.Month
	if len(input.MonthsOfYear) > 0 {
		monthsOfYear = make([]basebackuppolicyresources.Month, 0)
		for _, value := range input.MonthsOfYear {
			monthsOfYear = append(monthsOfYear, basebackuppolicyresources.Month(value))
		}
	}

	var weeksOfMonth []basebackuppolicyresources.WeekNumber
	if len(input.WeeksOfMonth) > 0 {
		weeksOfMonth = make([]basebackuppolicyresources.WeekNumber, 0)
		for _, value := range input.WeeksOfMonth {
			weeksOfMonth = append(weeksOfMonth, basebackuppolicyresources.WeekNumber(value))
		}
	}

	var scheduleTimes []string
	if len(input.ScheduledBackupTimes) > 0 {
		scheduleTimes = input.ScheduledBackupTimes
	}

	if len(absoluteCriteria) == 0 && len(daysOfWeek) == 0 && len(monthsOfYear) == 0 && len(weeksOfMonth) == 0 && len(scheduleTimes) == 0 {
		return nil
	}

	return &[]basebackuppolicyresources.BackupCriteria{
		basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To(absoluteCriteria),
			DaysOfTheWeek:    pointer.To(daysOfWeek),
			MonthsOfYear:     pointer.To(monthsOfYear),
			ScheduleTimes:    pointer.To(scheduleTimes),
			WeeksOfTheMonth:  pointer.To(weeksOfMonth),
		},
	}
}

func flattenBackupPolicyCosmosdbAccountPolicyRules(input []basebackuppolicyresources.BasePolicyRule) (string, []BackupPolicyCosmosdbAccountRetentionRule, []string, string) {
	var taggingCriteria []basebackuppolicyresources.TaggingCriteria
	var nonDefaultRetentionRules []basebackuppolicyresources.AzureRetentionRule
	var backupSchedule []string
	var timeZone string
	var defaultRetentionDuration string
	retentionRules := make([]BackupPolicyCosmosdbAccountRetentionRule, 0)

	for _, item := range input {
		switch rule := item.(type) {
		case basebackuppolicyresources.AzureBackupRule:
			if trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				backupSchedule = trigger.Schedule.RepeatingTimeIntervals
				timeZone = pointer.From(trigger.Schedule.TimeZone)
				taggingCriteria = trigger.TaggingCriteria
			}
		case basebackuppolicyresources.AzureRetentionRule:
			if pointer.From(rule.IsDefault) {
				if v := rule.Lifecycles; len(v) > 0 {
					if deleteOption, ok := v[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
						defaultRetentionDuration = deleteOption.Duration
					}
				}
			} else {
				nonDefaultRetentionRules = append(nonDefaultRetentionRules, rule)
			}
		}
	}

	for _, rule := range nonDefaultRetentionRules {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name: rule.Name,
		}

		for _, criteria := range taggingCriteria {
			if strings.EqualFold(criteria.TagInfo.TagName, rule.Name) {
				flattenBackupPolicyCosmosdbAccountCriteriaIntoRule(criteria.Criteria, &result)
				break
			}
		}

		if v := rule.Lifecycles; len(v) > 0 {
			if deleteOption, ok := v[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
				result.Duration = deleteOption.Duration
			}
		}

		retentionRules = append(retentionRules, result)
	}

	return defaultRetentionDuration, retentionRules, backupSchedule, timeZone
}

func flattenBackupPolicyCosmosdbAccountCriteriaIntoRule(input *[]basebackuppolicyresources.BackupCriteria, rule *BackupPolicyCosmosdbAccountRetentionRule) {
	if input == nil {
		return
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				rule.AbsoluteCriteria = string(pointer.From(criteria.AbsoluteCriteria)[0])
			}

			if criteria.DaysOfTheWeek != nil {
				daysOfWeek := make([]string, 0)
				for _, item := range pointer.From(criteria.DaysOfTheWeek) {
					daysOfWeek = append(daysOfWeek, string(item))
				}
				rule.DaysOfWeek = daysOfWeek
			}

			if criteria.MonthsOfYear != nil {
				monthsOfYear := make([]string, 0)
				for _, item := range pointer.From(criteria.MonthsOfYear) {
					monthsOfYear = append(monthsOfYear, string(item))
				}
				rule.MonthsOfYear = monthsOfYear
			}

			if criteria.WeeksOfTheMonth != nil {
				weeksOfMonth := make([]string, 0)
				for _, item := range pointer.From(criteria.WeeksOfTheMonth) {
					weeksOfMonth = append(weeksOfMonth, string(item))
				}
				rule.WeeksOfMonth = weeksOfMonth
			}

			if criteria.ScheduleTimes != nil {
				scheduleTimes := make([]string, 0)
				scheduleTimes = append(scheduleTimes, pointer.From(criteria.ScheduleTimes)...)
				rule.ScheduledBackupTimes = scheduleTimes
			}
		}
	}
}
