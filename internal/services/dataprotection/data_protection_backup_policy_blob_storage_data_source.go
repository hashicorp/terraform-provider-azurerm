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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-09-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupPolicyBlobStorageDataSourceModel struct {
	Name                                string                                           `tfschema:"name"`
	VaultId                             string                                           `tfschema:"vault_id"`
	BackupRepeatingTimeIntervals        []string                                         `tfschema:"backup_repeating_time_intervals"`
	OperationalDefaultRetentionDuration string                                           `tfschema:"operational_default_retention_duration"`
	RetentionRule                       []BackupPolicyBlobStorageDataSourceRetentionRule `tfschema:"retention_rule"`
	TimeZone                            string                                           `tfschema:"time_zone"`
	VaultDefaultRetentionDuration       string                                           `tfschema:"vault_default_retention_duration"`
}

type BackupPolicyBlobStorageDataSourceRetentionRule struct {
	Name      string                                       `tfschema:"name"`
	Criteria  []BackupPolicyBlobStorageDataSourceCriteria  `tfschema:"criteria"`
	LifeCycle []BackupPolicyBlobStorageDataSourceLifeCycle `tfschema:"life_cycle"`
	Priority  int64                                        `tfschema:"priority"`
}

type BackupPolicyBlobStorageDataSourceCriteria struct {
	AbsoluteCriteria     string   `tfschema:"absolute_criteria"`
	DaysOfMonth          []int64  `tfschema:"days_of_month"`
	DaysOfWeek           []string `tfschema:"days_of_week"`
	MonthsOfYear         []string `tfschema:"months_of_year"`
	ScheduledBackupTimes []string `tfschema:"scheduled_backup_times"`
	WeeksOfMonth         []string `tfschema:"weeks_of_month"`
}

type BackupPolicyBlobStorageDataSourceLifeCycle struct {
	DataStoreType string `tfschema:"data_store_type"`
	Duration      string `tfschema:"duration"`
}

type DataProtectionBackupPolicyBlobStorageDataSource struct{}

var _ sdk.DataSource = DataProtectionBackupPolicyBlobStorageDataSource{}

func (r DataProtectionBackupPolicyBlobStorageDataSource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_blob_storage"
}

func (r DataProtectionBackupPolicyBlobStorageDataSource) ModelObject() interface{} {
	return &BackupPolicyBlobStorageDataSourceModel{}
}

func (r DataProtectionBackupPolicyBlobStorageDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupVaultID,
		},
	}
}

func (r DataProtectionBackupPolicyBlobStorageDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"backup_repeating_time_intervals": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"operational_default_retention_duration": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"retention_rule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"criteria": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"absolute_criteria": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"days_of_month": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeInt,
									},
								},

								"days_of_week": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"months_of_year": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"scheduled_backup_times": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"weeks_of_month": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"life_cycle": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_store_type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"duration": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"priority": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vault_default_retention_duration": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupPolicyBlobStorageDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupPolicyClient

			var model BackupPolicyBlobStorageDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			vaultId, err := basebackuppolicyresources.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}

			id := basebackuppolicyresources.NewBackupPolicyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			resp, err := client.BackupPoliciesGet(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := BackupPolicyBlobStorageDataSourceModel{
				Name:    id.BackupPolicyName,
				VaultId: vaultId.ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if properties, ok := respModel.Properties.(basebackuppolicyresources.BackupPolicy); ok {
					state.BackupRepeatingTimeIntervals = flattenBackupPolicyBackupRepeatingTimeIntervals(properties.PolicyRules)
					state.TimeZone = flattenBackupPolicyBackupTimeZone(properties.PolicyRules)
					state.OperationalDefaultRetentionDuration = flattenBackupPolicyBlobStorageDataSourceDefaultRetentionDuration(properties.PolicyRules, basebackuppolicyresources.DataStoreTypesOperationalStore)
					state.VaultDefaultRetentionDuration = flattenBackupPolicyBlobStorageDataSourceDefaultRetentionDuration(properties.PolicyRules, basebackuppolicyresources.DataStoreTypesVaultStore)
					state.RetentionRule = flattenBackupPolicyBlobStorageDataSourceRetentionRules(properties.PolicyRules)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenBackupPolicyBlobStorageDataSourceDefaultRetentionDuration(input []basebackuppolicyresources.BasePolicyRule, dsType basebackuppolicyresources.DataStoreTypes) string {
	if input == nil {
		return ""
	}

	for _, item := range input {
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
					if (retentionRule.Lifecycles)[0].SourceDataStore.DataStoreType == dsType {
						return deleteOption.Duration
					}
				}
			}
		}
	}
	return ""
}

func flattenBackupPolicyBlobStorageDataSourceRetentionRules(input []basebackuppolicyresources.BasePolicyRule) []BackupPolicyBlobStorageDataSourceRetentionRule {
	results := make([]BackupPolicyBlobStorageDataSourceRetentionRule, 0)
	if len(input) == 0 {
		return results
	}

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
		if retentionRule, ok := item.(basebackuppolicyresources.AzureRetentionRule); ok && (retentionRule.IsDefault == nil || !*retentionRule.IsDefault) {
			name := retentionRule.Name
			var taggingPriority int64
			var taggingCriteria []BackupPolicyBlobStorageDataSourceCriteria
			for _, criteria := range taggingCriterias {
				if strings.EqualFold(criteria.TagInfo.TagName, name) {
					taggingPriority = criteria.TaggingPriority
					taggingCriteria = flattenBackupPolicyBlobStorageDataSourceCriteria(criteria.Criteria)
					break
				}
			}

			var lifeCycle []BackupPolicyBlobStorageDataSourceLifeCycle
			if v := retentionRule.Lifecycles; len(v) > 0 {
				lifeCycle = flattenBackupPolicyBlobStorageDataSourceLifeCycles(v, basebackuppolicyresources.DataStoreTypesVaultStore)
			}
			results = append(results, BackupPolicyBlobStorageDataSourceRetentionRule{
				Name:      name,
				Priority:  taggingPriority,
				Criteria:  taggingCriteria,
				LifeCycle: lifeCycle,
			})
		}
	}
	return results
}

func flattenBackupPolicyBlobStorageDataSourceLifeCycles(input []basebackuppolicyresources.SourceLifeCycle, dsType basebackuppolicyresources.DataStoreTypes) []BackupPolicyBlobStorageDataSourceLifeCycle {
	results := make([]BackupPolicyBlobStorageDataSourceLifeCycle, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		var duration string
		dataStoreType := item.SourceDataStore.DataStoreType
		if deleteOption, ok := item.DeleteAfter.(basebackuppolicyresources.AbsoluteDeleteOption); ok {
			if dataStoreType == dsType {
				duration = deleteOption.Duration
			} else {
				continue
			}
		}

		results = append(results, BackupPolicyBlobStorageDataSourceLifeCycle{
			Duration:      duration,
			DataStoreType: string(dataStoreType),
		})
	}
	return results
}

func flattenBackupPolicyBlobStorageDataSourceCriteria(input *[]basebackuppolicyresources.BackupCriteria) []BackupPolicyBlobStorageDataSourceCriteria {
	results := make([]BackupPolicyBlobStorageDataSourceCriteria, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if criteria, ok := item.(basebackuppolicyresources.ScheduleBasedBackupCriteria); ok {
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
			var daysOfMonth []int64
			if criteria.DaysOfMonth != nil {
				daysOfMonth = make([]int64, 0)
				for _, item := range *criteria.DaysOfMonth {
					daysOfMonth = append(daysOfMonth, pointer.From(item.Date))
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

			results = append(results, BackupPolicyBlobStorageDataSourceCriteria{
				AbsoluteCriteria:     absoluteCriteria,
				DaysOfWeek:           daysOfWeek,
				DaysOfMonth:          daysOfMonth,
				MonthsOfYear:         monthsOfYear,
				WeeksOfMonth:         weeksOfMonth,
				ScheduledBackupTimes: scheduleTimes,
			})
		}
	}
	return results
}
