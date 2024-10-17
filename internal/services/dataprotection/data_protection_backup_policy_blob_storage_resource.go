// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupPolicyBlobStorage() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceDataProtectionBackupPolicyBlobStorageCreate,
		Read:   resourceDataProtectionBackupPolicyBlobStorageRead,
		Delete: resourceDataProtectionBackupPolicyBlobStorageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := backuppolicies.ParseBackupPolicyID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"),
					"DataProtection BackupPolicy name must be 3 - 150 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: backuppolicies.ValidateBackupVaultID,
			},

			"backup_repeating_time_intervals": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"operational_default_retention_duration": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
				ValidateFunc: helperValidate.ISO8601Duration,
			},

			"time_zone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vault_default_retention_duration": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"operational_default_retention_duration", "vault_default_retention_duration"},
				RequiredWith: []string{"backup_repeating_time_intervals"},
				ValidateFunc: helperValidate.ISO8601Duration,
			},

			"retention_rule": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"vault_default_retention_duration"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
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
										ValidateFunc: validation.StringInSlice(
											backuppolicies.PossibleValuesForAbsoluteMarker(), false),
									},

									"days_of_month": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										ForceNew: true,
										MinItems: 1,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeInt,
											ValidateFunc: validation.Any(
												validation.IntBetween(0, 28),
											),
										},
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
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"data_store_type": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											// confirmed with the service team that currently only `VaultStore` is supported.
											// However, since `ArchiveStore` may be supported in the future, it is open to user specification.
											string(backuppolicies.DataStoreTypesVaultStore),
										}, false),
									},

									"duration": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: helperValidate.ISO8601Duration,
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

	return resource
}

func resourceDataProtectionBackupPolicyBlobStorageCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := backuppolicies.ParseBackupVaultID(d.Get("vault_id").(string))
	id := backuppolicies.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_blob_storage", id.ID())
	}

	policyRules := make([]backuppolicies.BasePolicyRule, 0)
	// expand the default operational retention rule when the operational default duration is specified
	operationalDefaultDuration := d.Get("operational_default_retention_duration").(string)
	if operationalDefaultDuration != "" {
		policyRules = append(policyRules, expandBackupPolicyBlobStorageDefaultRetentionRuleArray(operationalDefaultDuration, backuppolicies.DataStoreTypesOperationalStore))
	}

	// expand the default vault retention rule when the vault default duration is specified
	if v, ok := d.GetOk("vault_default_retention_duration"); ok {
		taggingCriteria, err := expandBackupPolicyBlobStorageTaggingCriteriaArray(d.Get("retention_rule").([]interface{}))
		if err != nil {
			return err
		}
		policyRules = append(policyRules, expandBackupPolicyBlobStorageAzureBackupRuleArray(d.Get("backup_repeating_time_intervals").([]interface{}), d.Get("time_zone").(string), taggingCriteria)...)
		policyRules = append(policyRules, expandBackupPolicyBlobStorageDefaultRetentionRuleArray(v.(string), backuppolicies.DataStoreTypesVaultStore))
	}

	// expand the vault retention rule when the vault retention rules are specified, the operational backup cannot specify retention rules.
	if _, ok := d.GetOk("retention_rule"); ok {
		policyRules = append(policyRules, expandBackupPolicyBlobStorageAzureRetentionRuleArray(d.Get("retention_rule").([]interface{}))...)
	}

	parameters := backuppolicies.BaseBackupPolicyResource{
		Properties: &backuppolicies.BackupPolicy{
			PolicyRules:     policyRules,
			DatasourceTypes: []string{"Microsoft.Storage/storageAccounts/blobServices"},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupPolicyBlobStorageRead(d, meta)
}

func resourceDataProtectionBackupPolicyBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
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
	vaultId := backuppolicies.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupPolicyName)
	d.Set("vault_id", vaultId.ID())
	if resp.Model != nil {
		if resp.Model.Properties != nil {
			if props, ok := resp.Model.Properties.(backuppolicies.BackupPolicy); ok {
				if err := d.Set("backup_repeating_time_intervals", flattenBackupPolicyBlobStorageVaultBackupRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `backup_repeating_time_intervals`: %+v", err)
				}
				if err := d.Set("operational_default_retention_duration", flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(props.PolicyRules, backuppolicies.DataStoreTypesOperationalStore)); err != nil {
					return fmt.Errorf("setting `operational_default_retention_duration`: %+v", err)
				}
				d.Set("time_zone", flattenBackupPolicyBlobStorageVaultBackupTimeZone(&props.PolicyRules))
				if err := d.Set("vault_default_retention_duration", flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(props.PolicyRules, backuppolicies.DataStoreTypesVaultStore)); err != nil {
					return fmt.Errorf("setting `vault_default_retention_duration`: %+v", err)
				}
				if err := d.Set("retention_rule", flattenBackupPolicyBlobStorageRetentionRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `retention_rule`: %+v", err)
				}
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupPolicyBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandBackupPolicyBlobStorageTaggingCriteriaArray(input []interface{}) (*[]backuppolicies.TaggingCriteria, error) {
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
				Id:      pointer.To(v["name"].(string) + "_"),
				TagName: v["name"].(string),
			},
		}

		criteria, err := expandBackupPolicyBlobStorageCriteriaArray(v["criteria"].([]interface{}))
		if err != nil {
			return nil, err
		}
		result.Criteria = criteria
		results = append(results, result)
	}
	return &results, nil
}

func expandBackupPolicyBlobStorageCriteriaArray(input []interface{}) (*[]backuppolicies.BackupCriteria, error) {
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

		var daysOfMonth []backuppolicies.Day
		if v["days_of_month"].(*pluginsdk.Set).Len() > 0 {
			daysOfMonth = make([]backuppolicies.Day, 0)
			for _, value := range v["days_of_month"].(*pluginsdk.Set).List() {
				isLast := false
				if value == 0 {
					isLast = true
				}
				daysOfMonth = append(daysOfMonth, backuppolicies.Day{
					Date: pointer.To(int64(value.(int))), IsLast: pointer.To(isLast),
				})
			}
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

		var scheduleTimes *[]string
		if v["scheduled_backup_times"].(*pluginsdk.Set).Len() > 0 {
			scheduleTimes = utils.ExpandStringSlice(v["scheduled_backup_times"].(*pluginsdk.Set).List())
		}

		var weeksOfMonth []backuppolicies.WeekNumber
		if v["weeks_of_month"].(*pluginsdk.Set).Len() > 0 {
			weeksOfMonth = make([]backuppolicies.WeekNumber, 0)
			for _, value := range v["weeks_of_month"].(*pluginsdk.Set).List() {
				weeksOfMonth = append(weeksOfMonth, backuppolicies.WeekNumber(value.(string)))
			}
		}

		results = append(results, backuppolicies.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			DaysOfMonth:      &daysOfMonth,
			DaysOfTheWeek:    &daysOfWeek,
			MonthsOfYear:     &monthsOfYear,
			ScheduleTimes:    scheduleTimes,
			WeeksOfTheMonth:  &weeksOfMonth,
		})
	}
	return &results, nil
}

func expandBackupPolicyBlobStorageAzureBackupRuleArray(input []interface{}, timeZone string, taggingCriteria *[]backuppolicies.TaggingCriteria) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	results = append(results, backuppolicies.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: backuppolicies.DataStoreInfoBase{
			DataStoreType: backuppolicies.DataStoreTypesVaultStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: &backuppolicies.AzureBackupParams{
			BackupType: "Discrete",
		},
		Trigger: backuppolicies.ScheduleBasedTriggerContext{
			Schedule: backuppolicies.BackupSchedule{
				RepeatingTimeIntervals: *utils.ExpandStringSlice(input),
				TimeZone:               pointer.To(timeZone),
			},
			TaggingCriteria: *taggingCriteria,
		},
	})

	return results
}

func expandBackupPolicyBlobStorageDefaultRetentionRuleArray(input interface{}, dataStoreType backuppolicies.DataStoreTypes) backuppolicies.BasePolicyRule {
	return backuppolicies.AzureRetentionRule{
		Name:      "Default",
		IsDefault: pointer.To(true),
		Lifecycles: []backuppolicies.SourceLifeCycle{
			{
				DeleteAfter: backuppolicies.AbsoluteDeleteOption{
					Duration: input.(string),
				},
				SourceDataStore: backuppolicies.DataStoreInfoBase{
					DataStoreType: dataStoreType,
					ObjectType:    "DataStoreInfoBase",
				},
				TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyBlobStorageAzureRetentionRuleArray(input []interface{}) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, backuppolicies.AzureRetentionRule{
			Name:       v["name"].(string),
			IsDefault:  pointer.To(false),
			Lifecycles: expandBackupPolicyBlobStorageLifeCycle(v["life_cycle"].([]interface{})),
		})
	}
	return results
}

func expandBackupPolicyBlobStorageLifeCycle(input []interface{}) []backuppolicies.SourceLifeCycle {
	results := make([]backuppolicies.SourceLifeCycle, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		sourceLifeCycle := backuppolicies.SourceLifeCycle{
			DeleteAfter: backuppolicies.AbsoluteDeleteOption{
				Duration: v["duration"].(string),
			},
			SourceDataStore: backuppolicies.DataStoreInfoBase{
				DataStoreType: backuppolicies.DataStoreTypes(v["data_store_type"].(string)),
				ObjectType:    "DataStoreInfoBase",
			},
			TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
		}
		results = append(results, sourceLifeCycle)
	}

	return results
}

func flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(input []backuppolicies.BasePolicyRule, dsType backuppolicies.DataStoreTypes) interface{} {
	if input == nil {
		return nil
	}

	for _, item := range input {
		if retentionRule, ok := item.(backuppolicies.AzureRetentionRule); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if retentionRule.Lifecycles != nil && len(retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (retentionRule.Lifecycles)[0].DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
					if (retentionRule.Lifecycles)[0].SourceDataStore.DataStoreType == dsType {
						return deleteOption.Duration
					}
				}
			}
		}
	}
	return nil
}

func flattenBackupPolicyBlobStorageVaultBackupRuleArray(input *[]backuppolicies.BasePolicyRule) []interface{} {
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

func flattenBackupPolicyBlobStorageVaultBackupTimeZone(input *[]backuppolicies.BasePolicyRule) string {
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

func flattenBackupPolicyBlobStorageRetentionRuleArray(input *[]backuppolicies.BasePolicyRule) []interface{} {
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
					taggingCriteria = flattenBackupPolicyBlobStorageBackupCriteriaArray(criteria.Criteria)
					break
				}
			}

			var lifeCycle []interface{}
			if v := retentionRule.Lifecycles; len(v) > 0 {
				lifeCycle = flattenBackupPolicyBlobStorageBackupLifeCycleArray(v, backuppolicies.DataStoreTypesVaultStore)
			}
			results = append(results, map[string]interface{}{
				"name":       name,
				"priority":   taggingPriority,
				"criteria":   taggingCriteria,
				"life_cycle": lifeCycle,
			})
		}
	}
	return results
}

func flattenBackupPolicyBlobStorageBackupLifeCycleArray(input []backuppolicies.SourceLifeCycle, dsType backuppolicies.DataStoreTypes) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		var duration string
		dataStoreType := item.SourceDataStore.DataStoreType
		if deleteOption, ok := item.DeleteAfter.(backuppolicies.AbsoluteDeleteOption); ok {
			if dataStoreType == dsType {
				duration = deleteOption.Duration
			} else {
				continue
			}
		}

		results = append(results, map[string]interface{}{
			"duration":        duration,
			"data_store_type": string(dataStoreType),
		})
	}
	return results
}

func flattenBackupPolicyBlobStorageBackupCriteriaArray(input *[]backuppolicies.BackupCriteria) []interface{} {
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
			var daysOfMonth []int
			if criteria.DaysOfMonth != nil {
				daysOfMonth = make([]int, 0)
				for _, item := range *criteria.DaysOfMonth {
					daysOfMonth = append(daysOfMonth, (int)(pointer.From(item.Date)))
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
				"days_of_month":          daysOfMonth,
				"months_of_year":         monthsOfYear,
				"weeks_of_month":         weeksOfMonth,
				"scheduled_backup_times": scheduleTimes,
			})
		}
	}
	return results
}
