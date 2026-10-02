// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/btubbs/datetime"
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
	DailyBackupEnabled          bool                                       `tfschema:"daily_backup_enabled"`
	IncrementalBackupSchedules  []string                                   `tfschema:"incremental_backup_schedules"`
	DefaultRetentionDuration    string                                     `tfschema:"default_retention_duration"`
	BackupSchedule              string                                     `tfschema:"backup_schedule"`
	RetentionRules              []BackupPolicyCosmosdbAccountRetentionRule `tfschema:"retention_rule"`
	TimeZone                    string                                     `tfschema:"time_zone"`
}

type BackupPolicyCosmosdbAccountRetentionRule struct {
	Name             string   `tfschema:"name"`
	Duration         string   `tfschema:"duration"`
	BackupOccurrence string   `tfschema:"backup_occurrence"`
	DaysOfWeek       []string `tfschema:"days_of_week"`
	MonthsOfYear     []string `tfschema:"months_of_year"`
	WeeksOfMonth     []string `tfschema:"weeks_of_month"`
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

func (r DataProtectionBackupPolicyCosmosdbAccountResource) ModelObject() any {
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

		"backup_schedule": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.ISO8601RepeatingTime,
				validate.BackupPolicyCosmosdbAccountBackupSchedule(),
			),
		},

		"daily_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
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

					// Only Week, Month and Year are supported for cosmosdb account backup policies
					"backup_occurrence": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(basebackuppolicyresources.AbsoluteMarkerFirstOfMonth),
							string(basebackuppolicyresources.AbsoluteMarkerFirstOfWeek),
							string(basebackuppolicyresources.AbsoluteMarkerFirstOfYear),
						}, false),
					},

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
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
						MaxItems: 12,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForMonth(), false),
						},
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 5,
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
	return map[string]*pluginsdk.Schema{
		"incremental_backup_schedules": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
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
				hasBackupOccurrence := rule.BackupOccurrence != ""
				hasDaysOfWeek := len(rule.DaysOfWeek) > 0
				if (hasBackupOccurrence && hasDaysOfWeek) || (!hasBackupOccurrence && !hasDaysOfWeek) {
					return fmt.Errorf("`retention_rule` %q requires exactly one of `backup_occurrence` and `days_of_week` to be specified", rule.Name)
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
			backupRules, err := expandBackupPolicyCosmosdbAccountBackupRules(model.BackupSchedule, model.DailyBackupEnabled, model.TimeZone, expandBackupPolicyCosmosdbAccountTaggingCriteria(model.RetentionRules))
			if err != nil {
				return fmt.Errorf("expanding backup schedule: %+v", err)
			}
			policyRules = append(policyRules, backupRules...)

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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r DataProtectionBackupPolicyCosmosdbAccountResource) flatten(metadata sdk.ResourceMetaData, id *basebackuppolicyresources.BackupPolicyId, model *basebackuppolicyresources.BaseBackupPolicyResource) error {
	vaultId := basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	state := BackupPolicyCosmosdbAccountModel{
		Name:                        id.BackupPolicyName,
		DataProtectionBackupVaultId: vaultId.ID(),
	}

	if model != nil {
		if properties, ok := model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
			state.DefaultRetentionDuration, state.RetentionRules, state.BackupSchedule, state.IncrementalBackupSchedules, state.DailyBackupEnabled, state.TimeZone = flattenBackupPolicyCosmosdbAccountPolicyRules(properties.PolicyRules)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
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

func expandBackupPolicyCosmosdbAccountBackupRules(fullBackupSchedule string, incrementalBackupEnabled bool, timeZone string, taggingCriteria []basebackuppolicyresources.TaggingCriteria) ([]basebackuppolicyresources.BasePolicyRule, error) {
	results := []basebackuppolicyresources.BasePolicyRule{
		basebackuppolicyresources.AzureBackupRule{
			Name: "BackupWeekly",
			DataStore: basebackuppolicyresources.DataStoreInfoBase{
				DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
				ObjectType:    "DataStoreInfoBase",
			},
			BackupParameters: basebackuppolicyresources.AzureBackupParams{
				BackupType: "Full",
			},
			Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
				Schedule: basebackuppolicyresources.BackupSchedule{
					RepeatingTimeIntervals: []string{fullBackupSchedule},
					TimeZone:               pointer.To(timeZone),
				},
				TaggingCriteria: taggingCriteria,
			},
		},
	}

	if !incrementalBackupEnabled {
		return results, nil
	}

	incrementalBackupSchedule, err := generateBackupPolicyCosmosdbAccountIncrementalSchedules(fullBackupSchedule)
	if err != nil {
		return nil, err
	}

	results = append(results, basebackuppolicyresources.AzureBackupRule{
		Name: "BackupWeeklyIncremental",
		DataStore: basebackuppolicyresources.DataStoreInfoBase{
			DataStoreType: basebackuppolicyresources.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: basebackuppolicyresources.AzureBackupParams{
			BackupType: "Incremental",
		},
		Trigger: basebackuppolicyresources.ScheduleBasedTriggerContext{
			Schedule: basebackuppolicyresources.BackupSchedule{
				RepeatingTimeIntervals: incrementalBackupSchedule,
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: []basebackuppolicyresources.TaggingCriteria{
				expandBackupPolicyCosmosdbAccountDefaultTaggingCriteria(),
			},
		},
	})

	return results, nil
}

func generateBackupPolicyCosmosdbAccountIncrementalSchedules(fullBackupSchedule string) ([]string, error) {
	const (
		recurrencePrefix = "R/"
		recurrenceSuffix = "/P1W"
		timeFormat       = "2006-01-02T15:04:05-07:00"
	)

	timestamp := strings.TrimSuffix(strings.TrimPrefix(fullBackupSchedule, recurrencePrefix), recurrenceSuffix)
	fullBackupTime, err := datetime.Parse(timestamp, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("parsing timestamp in `backup_schedule` value `%s`: %+v", fullBackupSchedule, err)
	}

	results := make([]string, 0, 6)
	for day := 1; day <= 6; day++ {
		incrementalBackupTime := fullBackupTime.Add(time.Duration(day) * 24 * time.Hour)
		results = append(results, fmt.Sprintf("%s%s%s", recurrencePrefix, incrementalBackupTime.Format(timeFormat), recurrenceSuffix))
	}

	return results, nil
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
		expandBackupPolicyCosmosdbAccountDefaultTaggingCriteria(),
	}

	retentionRules := append([]BackupPolicyCosmosdbAccountRetentionRule(nil), input...)
	sort.SliceStable(retentionRules, func(i, j int) bool {
		return backupPolicyCosmosdbAccountRetentionRuleOrder(retentionRules[i]) < backupPolicyCosmosdbAccountRetentionRuleOrder(retentionRules[j])
	})

	for i, item := range retentionRules {
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

func backupPolicyCosmosdbAccountRetentionRuleOrder(input BackupPolicyCosmosdbAccountRetentionRule) int {
	if input.BackupOccurrence == string(basebackuppolicyresources.AbsoluteMarkerFirstOfYear) || len(input.MonthsOfYear) > 0 {
		return 1
	}
	if input.BackupOccurrence == string(basebackuppolicyresources.AbsoluteMarkerFirstOfMonth) || len(input.WeeksOfMonth) > 0 {
		return 2
	}
	return 3
}

func expandBackupPolicyCosmosdbAccountDefaultTaggingCriteria() basebackuppolicyresources.TaggingCriteria {
	return basebackuppolicyresources.TaggingCriteria{
		IsDefault:       true,
		TaggingPriority: 99,
		TagInfo: basebackuppolicyresources.RetentionTag{
			Id:      pointer.To("Default_"),
			TagName: "Default",
		},
	}
}

func expandBackupPolicyCosmosdbAccountRetentionRuleCriteria(input BackupPolicyCosmosdbAccountRetentionRule) *[]basebackuppolicyresources.BackupCriteria {
	var absoluteCriteria []basebackuppolicyresources.AbsoluteMarker
	if len(input.BackupOccurrence) > 0 {
		absoluteCriteria = []basebackuppolicyresources.AbsoluteMarker{basebackuppolicyresources.AbsoluteMarker(input.BackupOccurrence)}
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

	if len(absoluteCriteria) == 0 && len(daysOfWeek) == 0 && len(monthsOfYear) == 0 && len(weeksOfMonth) == 0 {
		return nil
	}

	return &[]basebackuppolicyresources.BackupCriteria{
		basebackuppolicyresources.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: pointer.To(absoluteCriteria),
			DaysOfTheWeek:    pointer.To(daysOfWeek),
			MonthsOfYear:     pointer.To(monthsOfYear),
			WeeksOfTheMonth:  pointer.To(weeksOfMonth),
		},
	}
}

func flattenBackupPolicyCosmosdbAccountPolicyRules(input []basebackuppolicyresources.BasePolicyRule) (string, []BackupPolicyCosmosdbAccountRetentionRule, string, []string, bool, string) {
	var taggingCriteria []basebackuppolicyresources.TaggingCriteria
	var nonDefaultRetentionRules []basebackuppolicyresources.AzureRetentionRule
	var incrementalBackupEnabled bool
	var fullBackupSchedule string
	var timeZone string
	var defaultRetentionDuration string

	incrementalBackupSchedule := make([]string, 0)
	retentionRules := make([]BackupPolicyCosmosdbAccountRetentionRule, 0)

	for _, item := range input {
		switch rule := item.(type) {
		case basebackuppolicyresources.AzureBackupRule:
			if trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if timeZone == "" {
					timeZone = pointer.From(trigger.Schedule.TimeZone)
				}

				if parameters, ok := rule.BackupParameters.(basebackuppolicyresources.AzureBackupParams); ok {
					switch {
					case strings.EqualFold(parameters.BackupType, "Full"):
						if len(trigger.Schedule.RepeatingTimeIntervals) > 0 {
							fullBackupSchedule = trigger.Schedule.RepeatingTimeIntervals[0]
						}
						taggingCriteria = trigger.TaggingCriteria
					case strings.EqualFold(parameters.BackupType, "Incremental"):
						incrementalBackupEnabled = true
						incrementalBackupSchedule = append(incrementalBackupSchedule, trigger.Schedule.RepeatingTimeIntervals...)
					}
				}
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

	return defaultRetentionDuration, retentionRules, fullBackupSchedule, incrementalBackupSchedule, incrementalBackupEnabled, timeZone
}

func flattenBackupPolicyCosmosdbAccountCriteriaIntoRule(input *[]basebackuppolicyresources.BackupCriteria, rule *BackupPolicyCosmosdbAccountRetentionRule) {
	if input == nil {
		return
	}

	for _, item := range pointer.From(input) {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
			if criteria.AbsoluteCriteria != nil && len(pointer.From(criteria.AbsoluteCriteria)) > 0 {
				rule.BackupOccurrence = string(pointer.From(criteria.AbsoluteCriteria)[0])
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
		}
	}
}
