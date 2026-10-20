// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"errors"
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

const (
	backupPolicyCosmosdbAccountFrequencyDaily  = "Daily"
	backupPolicyCosmosdbAccountFrequencyWeekly = "Weekly"
)

type BackupPolicyCosmosdbAccountModel struct {
	Name                        string                                            `tfschema:"name"`
	DataProtectionBackupVaultId string                                            `tfschema:"data_protection_backup_vault_id"`
	BackupFrequency             string                                            `tfschema:"backup_frequency"`
	IncrementalBackupSchedules  []string                                          `tfschema:"incremental_backup_schedules"`
	DefaultRetentionDuration    string                                            `tfschema:"default_retention_duration"`
	BackupSchedule              string                                            `tfschema:"backup_schedule"`
	WeeklyRetentionRule         []BackupPolicyCosmosdbAccountWeeklyRetentionRule  `tfschema:"weekly_retention_rule"`
	MonthlyRetentionRule        []BackupPolicyCosmosdbAccountMonthlyRetentionRule `tfschema:"monthly_retention_rule"`
	YearlyRetentionRule         []BackupPolicyCosmosdbAccountYearlyRetentionRule  `tfschema:"yearly_retention_rule"`
}

type BackupPolicyCosmosdbAccountWeeklyRetentionRule struct {
	Name       string   `tfschema:"name"`
	Duration   string   `tfschema:"duration"`
	DaysOfWeek []string `tfschema:"days_of_week"`
}

type BackupPolicyCosmosdbAccountMonthlyRetentionRule struct {
	Name         string   `tfschema:"name"`
	Duration     string   `tfschema:"duration"`
	DaysOfWeek   []string `tfschema:"days_of_week"`
	WeeksOfMonth []string `tfschema:"weeks_of_month"`
}

type BackupPolicyCosmosdbAccountYearlyRetentionRule struct {
	Name         string   `tfschema:"name"`
	Duration     string   `tfschema:"duration"`
	DaysOfWeek   []string `tfschema:"days_of_week"`
	WeeksOfMonth []string `tfschema:"weeks_of_month"`
	MonthsOfYear []string `tfschema:"months_of_year"`
}

type BackupPolicyCosmosdbAccountRetentionRule struct {
	Name             string
	Duration         string
	BackupOccurrence string
	DaysOfWeek       []string
	MonthsOfYear     []string
	WeeksOfMonth     []string
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
				validate.BackupPolicyCosmosdbAccountBackupSchedule(),
				validation.ISO8601RepeatingTime,
			),
		},

		"default_retention_duration": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.ISO8601Duration,
		},

		"backup_frequency": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  backupPolicyCosmosdbAccountFrequencyDaily,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				backupPolicyCosmosdbAccountFrequencyDaily,
				backupPolicyCosmosdbAccountFrequencyWeekly,
			}, false),
		},

		"monthly_retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
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

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						RequiredWith: []string{
							"monthly_retention_rule.0.days_of_week",
							"monthly_retention_rule.0.weeks_of_month",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForDayOfWeek(), false),
						},
					},

					"weeks_of_month": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 5,
						RequiredWith: []string{
							"monthly_retention_rule.0.days_of_week",
							"monthly_retention_rule.0.weeks_of_month",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForWeekNumber(), false),
						},
					},
				},
			},
		},

		"weekly_retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
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
				},
			},
		},

		"yearly_retention_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
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

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						ForceNew: true,
						MinItems: 1,
						MaxItems: 1,
						RequiredWith: []string{
							"yearly_retention_rule.0.days_of_week",
							"yearly_retention_rule.0.weeks_of_month",
							"yearly_retention_rule.0.months_of_year",
						},
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
						RequiredWith: []string{
							"yearly_retention_rule.0.days_of_week",
							"yearly_retention_rule.0.weeks_of_month",
							"yearly_retention_rule.0.months_of_year",
						},
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
						RequiredWith: []string{
							"yearly_retention_rule.0.days_of_week",
							"yearly_retention_rule.0.weeks_of_month",
							"yearly_retention_rule.0.months_of_year",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(basebackuppolicyresources.PossibleValuesForWeekNumber(), false),
						},
					},
				},
			},
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

			retentionRules := expandBackupPolicyCosmosdbAccountRetentionRuleModels(model)
			policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountRetentionRules(retentionRules)...)
			policyRules = append(policyRules, expandBackupPolicyCosmosdbAccountDefaultRetentionRule(model.DefaultRetentionDuration))
			backupRules, err := expandBackupPolicyCosmosdbAccountBackupRules(model.BackupSchedule, model.BackupFrequency == backupPolicyCosmosdbAccountFrequencyDaily, expandBackupPolicyCosmosdbAccountTaggingCriteria(retentionRules))
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

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			properties, ok := resp.Model.Properties.(basebackuppolicyresources.BackupPolicy)
			if !ok || !backupPolicySupportsCosmosdbAccounts(properties) {
				return fmt.Errorf("%s is not a backup policy for `Microsoft.DocumentDB/databaseAccounts`", *id)
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
			var retentionRules []BackupPolicyCosmosdbAccountRetentionRule
			state.DefaultRetentionDuration, retentionRules, state.BackupSchedule, state.IncrementalBackupSchedules, state.BackupFrequency = flattenBackupPolicyCosmosdbAccountPolicyRules(properties.PolicyRules)

			if err := flattenBackupPolicyCosmosdbAccountRetentionRules(retentionRules, &state); err != nil {
				return fmt.Errorf("flattening retention rules for %s: %+v", *id, err)
			}
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

			resp, err := client.BackupPoliciesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			properties, ok := resp.Model.Properties.(basebackuppolicyresources.BackupPolicy)
			if !ok || !backupPolicySupportsCosmosdbAccounts(properties) {
				return fmt.Errorf("%s is not a backup policy for `Microsoft.DocumentDB/databaseAccounts`", *id)
			}

			if _, err := client.BackupPoliciesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func expandBackupPolicyCosmosdbAccountRetentionRuleModels(input BackupPolicyCosmosdbAccountModel) []BackupPolicyCosmosdbAccountRetentionRule {
	results := make([]BackupPolicyCosmosdbAccountRetentionRule, 0)
	for _, rule := range input.YearlyRetentionRule {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name:         rule.Name,
			Duration:     rule.Duration,
			DaysOfWeek:   rule.DaysOfWeek,
			WeeksOfMonth: rule.WeeksOfMonth,
			MonthsOfYear: rule.MonthsOfYear,
		}
		if len(rule.DaysOfWeek) == 0 {
			result.BackupOccurrence = string(basebackuppolicyresources.AbsoluteMarkerFirstOfYear)
		}
		results = append(results, result)
	}
	for _, rule := range input.MonthlyRetentionRule {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name:         rule.Name,
			Duration:     rule.Duration,
			DaysOfWeek:   rule.DaysOfWeek,
			WeeksOfMonth: rule.WeeksOfMonth,
		}
		if len(rule.DaysOfWeek) == 0 {
			result.BackupOccurrence = string(basebackuppolicyresources.AbsoluteMarkerFirstOfMonth)
		}
		results = append(results, result)
	}
	for _, rule := range input.WeeklyRetentionRule {
		result := BackupPolicyCosmosdbAccountRetentionRule{
			Name:       rule.Name,
			Duration:   rule.Duration,
			DaysOfWeek: rule.DaysOfWeek,
		}
		if len(rule.DaysOfWeek) == 0 {
			result.BackupOccurrence = string(basebackuppolicyresources.AbsoluteMarkerFirstOfWeek)
		}
		results = append(results, result)
	}
	return results
}

func flattenBackupPolicyCosmosdbAccountRetentionRules(input []BackupPolicyCosmosdbAccountRetentionRule, state *BackupPolicyCosmosdbAccountModel) error {
	for _, rule := range input {
		switch {
		case rule.BackupOccurrence == string(basebackuppolicyresources.AbsoluteMarkerFirstOfYear) || len(rule.MonthsOfYear) > 0:
			if len(state.YearlyRetentionRule) > 0 {
				return errors.New("multiple yearly retention rules cannot be represented by `yearly_retention_rule`")
			}
			state.YearlyRetentionRule = []BackupPolicyCosmosdbAccountYearlyRetentionRule{{
				Name:         rule.Name,
				Duration:     rule.Duration,
				DaysOfWeek:   rule.DaysOfWeek,
				WeeksOfMonth: rule.WeeksOfMonth,
				MonthsOfYear: rule.MonthsOfYear,
			}}
		case rule.BackupOccurrence == string(basebackuppolicyresources.AbsoluteMarkerFirstOfMonth) || len(rule.WeeksOfMonth) > 0:
			if len(state.MonthlyRetentionRule) > 0 {
				return errors.New("multiple monthly retention rules cannot be represented by `monthly_retention_rule`")
			}
			state.MonthlyRetentionRule = []BackupPolicyCosmosdbAccountMonthlyRetentionRule{{
				Name:         rule.Name,
				Duration:     rule.Duration,
				DaysOfWeek:   rule.DaysOfWeek,
				WeeksOfMonth: rule.WeeksOfMonth,
			}}
		default:
			if len(state.WeeklyRetentionRule) > 0 {
				return errors.New("multiple weekly retention rules cannot be represented by `weekly_retention_rule`")
			}
			state.WeeklyRetentionRule = []BackupPolicyCosmosdbAccountWeeklyRetentionRule{{
				Name:       rule.Name,
				Duration:   rule.Duration,
				DaysOfWeek: rule.DaysOfWeek,
			}}
		}
	}
	return nil
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

func expandBackupPolicyCosmosdbAccountBackupRules(fullBackupSchedule string, incrementalBackupEnabled bool, taggingCriteria []basebackuppolicyresources.TaggingCriteria) ([]basebackuppolicyresources.BasePolicyRule, error) {
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
	fullBackupTime, err := time.Parse(time.RFC3339, timestamp)
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

func flattenBackupPolicyCosmosdbAccountPolicyRules(input []basebackuppolicyresources.BasePolicyRule) (string, []BackupPolicyCosmosdbAccountRetentionRule, string, []string, string) {
	var taggingCriteria []basebackuppolicyresources.TaggingCriteria
	var nonDefaultRetentionRules []basebackuppolicyresources.AzureRetentionRule
	backupFrequency := backupPolicyCosmosdbAccountFrequencyWeekly
	var fullBackupSchedule string
	var defaultRetentionDuration string

	incrementalBackupSchedule := make([]string, 0)
	retentionRules := make([]BackupPolicyCosmosdbAccountRetentionRule, 0)

	for _, item := range input {
		switch rule := item.(type) {
		case basebackuppolicyresources.AzureBackupRule:
			if trigger, ok := rule.Trigger.(basebackuppolicyresources.ScheduleBasedTriggerContext); ok {
				if parameters, ok := rule.BackupParameters.(basebackuppolicyresources.AzureBackupParams); ok {
					switch {
					case strings.EqualFold(parameters.BackupType, "Full"):
						if len(trigger.Schedule.RepeatingTimeIntervals) > 0 {
							fullBackupSchedule = trigger.Schedule.RepeatingTimeIntervals[0]
						}
						taggingCriteria = trigger.TaggingCriteria
					case strings.EqualFold(parameters.BackupType, "Incremental"):
						backupFrequency = backupPolicyCosmosdbAccountFrequencyDaily
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

	return defaultRetentionDuration, retentionRules, fullBackupSchedule, incrementalBackupSchedule, backupFrequency
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
