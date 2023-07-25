// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupPolicyPostgreSQL() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataProtectionBackupPolicyPostgreSQLCreate,
		Read:   resourceDataProtectionBackupPolicyPostgreSQLRead,
		Delete: resourceDataProtectionBackupPolicyPostgreSQLDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := backuppolicies.ParseBackupPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
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
		},
	}
}

func resourceDataProtectionBackupPolicyPostgreSQLCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("vault_name").(string)

	id := backuppolicies.NewBackupPolicyID(subscriptionId, resourceGroup, vaultName, name)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_postgresql", id.ID())
	}

	taggingCriteria, err := expandBackupPolicyPostgreSQLTaggingCriteriaArray(d.Get("retention_rule").([]interface{}))
	if err != nil {
		return err
	}

	policyRules := make([]backuppolicies.BasePolicyRule, 0)
	policyRules = append(policyRules, expandBackupPolicyPostgreSQLAzureBackupRuleArray(d.Get("backup_repeating_time_intervals").([]interface{}), taggingCriteria)...)
	policyRules = append(policyRules, expandBackupPolicyPostgreSQLDefaultAzureRetentionRule(d.Get("default_retention_duration")))
	policyRules = append(policyRules, expandBackupPolicyPostgreSQLAzureRetentionRuleArray(d.Get("retention_rule").([]interface{}))...)
	parameters := backuppolicies.BaseBackupPolicyResource{
		Properties: &backuppolicies.BackupPolicy{
			PolicyRules:     policyRules,
			DatasourceTypes: []string{"Microsoft.DBforPostgreSQL/servers/databases"},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupPolicyPostgreSQLRead(d, meta)
}

func resourceDataProtectionBackupPolicyPostgreSQLRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backuppolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	d.Set("name", id.BackupPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("vault_name", id.BackupVaultName)

	if resp.Model != nil {
		if resp.Model.Properties != nil {
			if props, ok := resp.Model.Properties.(backuppolicies.BackupPolicy); ok {
				if err := d.Set("backup_repeating_time_intervals", flattenBackupPolicyPostgreSQLBackupRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `backup_rule`: %+v", err)
				}
				if err := d.Set("default_retention_duration", flattenBackupPolicyPostgreSQLDefaultRetentionRuleDuration(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `default_retention_duration`: %+v", err)
				}
				if err := d.Set("retention_rule", flattenBackupPolicyPostgreSQLRetentionRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `retention_rule`: %+v", err)
				}
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupPolicyPostgreSQLDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backuppolicies.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return nil
}

func expandBackupPolicyPostgreSQLAzureBackupRuleArray(input []interface{}, taggingCriteria *[]backuppolicies.TaggingCriteria) []backuppolicies.BasePolicyRule {
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
				RepeatingTimeIntervals: *utils.ExpandStringSlice(input),
			},
			TaggingCriteria: *taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyPostgreSQLAzureRetentionRuleArray(input []interface{}) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, backuppolicies.AzureRetentionRule{
			Name:      v["name"].(string),
			IsDefault: utils.Bool(false),
			Lifecycles: []backuppolicies.SourceLifeCycle{
				{
					DeleteAfter: backuppolicies.AbsoluteDeleteOption{
						Duration: v["duration"].(string),
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

func expandBackupPolicyPostgreSQLDefaultAzureRetentionRule(input interface{}) backuppolicies.BasePolicyRule {
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

func expandBackupPolicyPostgreSQLTaggingCriteriaArray(input []interface{}) (*[]backuppolicies.TaggingCriteria, error) {
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
		v := item.(map[string]interface{})
		result := backuppolicies.TaggingCriteria{
			IsDefault:       false,
			TaggingPriority: int64(v["priority"].(int)),
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String(v["name"].(string) + "_"),
				TagName: v["name"].(string),
			},
		}

		criteria, err := expandBackupPolicyPostgreSQLCriteriaArray(v["criteria"].([]interface{}))
		if err != nil {
			return nil, err
		}
		result.Criteria = criteria

		results = append(results, result)
	}
	return &results, nil
}

func expandBackupPolicyPostgreSQLCriteriaArray(input []interface{}) (*[]backuppolicies.BackupCriteria, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, fmt.Errorf("criteria is a required field, cannot leave blank")
	}
	results := make([]backuppolicies.BackupCriteria, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		var absoluteCriteria []backuppolicies.AbsoluteMarker
		if absoluteCriteriaRaw := v["absolute_criteria"].(string); len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []backuppolicies.AbsoluteMarker{backuppolicies.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfWeek []backuppolicies.DayOfWeek
		if v["days_of_week"].(*pluginsdk.Set).Len() > 0 {
			daysOfWeek = make([]backuppolicies.DayOfWeek, 0)
			for _, value := range v["days_of_week"].(*pluginsdk.Set).List() {
				daysOfWeek = append(daysOfWeek, backuppolicies.DayOfWeek(value.(string)))
			}
		}

		var monthsOfYear []backuppolicies.Month
		if v["months_of_year"].(*pluginsdk.Set).Len() > 0 {
			monthsOfYear = make([]backuppolicies.Month, 0)
			for _, value := range v["months_of_year"].(*pluginsdk.Set).List() {
				monthsOfYear = append(monthsOfYear, backuppolicies.Month(value.(string)))
			}
		}

		var weeksOfMonth []backuppolicies.WeekNumber
		if v["weeks_of_month"].(*pluginsdk.Set).Len() > 0 {
			weeksOfMonth = make([]backuppolicies.WeekNumber, 0)
			for _, value := range v["weeks_of_month"].(*pluginsdk.Set).List() {
				weeksOfMonth = append(weeksOfMonth, backuppolicies.WeekNumber(value.(string)))
			}
		}

		var scheduleTimes *[]string
		if v["scheduled_backup_times"].(*pluginsdk.Set).Len() > 0 {
			scheduleTimes = utils.ExpandStringSlice(v["scheduled_backup_times"].(*pluginsdk.Set).List())
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

func flattenBackupPolicyPostgreSQLBackupRuleArray(input *[]backuppolicies.BasePolicyRule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	for _, item := range *input {
		if backupRule, ok := item.(backuppolicies.AzureBackupRule); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.(backuppolicies.ScheduleBasedTriggerContext); ok {
					return utils.FlattenStringSlice(&scheduleBasedTrigger.Schedule.RepeatingTimeIntervals)
				}
			}
		}
	}
	return make([]interface{}, 0)
}

func flattenBackupPolicyPostgreSQLDefaultRetentionRuleDuration(input *[]backuppolicies.BasePolicyRule) interface{} {
	if input == nil {
		return nil
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
	return nil
}

func flattenBackupPolicyPostgreSQLRetentionRuleArray(input *[]backuppolicies.BasePolicyRule) []interface{} {
	results := make([]interface{}, 0)
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
			var taggingPriority int64
			var taggingCriteria []interface{}
			for _, criteria := range taggingCriterias {
				if strings.EqualFold(criteria.TagInfo.TagName, name) {
					taggingPriority = criteria.TaggingPriority
					taggingCriteria = flattenBackupPolicyPostgreSQLBackupCriteriaArray(criteria.Criteria)
				}
			}
			var duration string
			if retentionRule.Lifecycles != nil && len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
					duration = deleteOption.Duration
				}
			}
			results = append(results, map[string]interface{}{
				"name":     name,
				"priority": taggingPriority,
				"criteria": taggingCriteria,
				"duration": duration,
			})
		}
	}
	return results
}

func flattenBackupPolicyPostgreSQLBackupCriteriaArray(input *[]backuppolicies.BackupCriteria) []interface{} {
	results := make([]interface{}, 0)
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

			results = append(results, map[string]interface{}{
				"absolute_criteria":      absoluteCriteria,
				"days_of_week":           daysOfWeek,
				"months_of_year":         monthsOfYear,
				"weeks_of_month":         weeksOfMonth,
				"scheduled_backup_times": scheduleTimes,
			})
		}
	}
	return results
}
