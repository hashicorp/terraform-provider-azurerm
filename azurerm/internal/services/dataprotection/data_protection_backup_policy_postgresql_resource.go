package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			_, err := parse.BackupPolicyID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
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
										ValidateFunc: validation.StringInSlice([]string{
											string(dataprotection.AllBackup),
											string(dataprotection.FirstOfDay),
											string(dataprotection.FirstOfMonth),
											string(dataprotection.FirstOfWeek),
											string(dataprotection.FirstOfYear),
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
												string(dataprotection.First),
												string(dataprotection.Second),
												string(dataprotection.Third),
												string(dataprotection.Fourth),
												string(dataprotection.Last),
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

	id := parse.NewBackupPolicyID(subscriptionId, resourceGroup, vaultName, name)

	existing, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_postgresql", id.ID())
	}

	taggingCriteria := expandBackupPolicyPostgreSQLTaggingCriteriaArray(d.Get("retention_rule").([]interface{}))
	policyRules := make([]dataprotection.BasicBasePolicyRule, 0)
	policyRules = append(policyRules, expandBackupPolicyPostgreSQLAzureBackupRuleArray(d.Get("backup_repeating_time_intervals").([]interface{}), taggingCriteria)...)
	policyRules = append(policyRules, expandBackupPolicyPostgreSQLDefaultAzureRetentionRule(d.Get("default_retention_duration")))
	policyRules = append(policyRules, expandBackupPolicyPostgreSQLAzureRetentionRuleArray(d.Get("retention_rule").([]interface{}))...)
	parameters := dataprotection.BaseBackupPolicyResource{
		Properties: &dataprotection.BackupPolicy{
			PolicyRules:     &policyRules,
			DatasourceTypes: &[]string{"Microsoft.DBforPostgreSQL/servers/databases"},
			ObjectType:      dataprotection.ObjectTypeBackupPolicy,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.BackupVaultName, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupPolicyPostgreSQLRead(d, meta)
}

func resourceDataProtectionBackupPolicyPostgreSQLRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("vault_name", id.BackupVaultName)
	if resp.Properties != nil {
		if props, ok := resp.Properties.AsBackupPolicy(); ok {
			if err := d.Set("backup_repeating_time_intervals", flattenBackupPolicyPostgreSQLBackupRuleArray(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `backup_rule`: %+v", err)
			}
			if err := d.Set("default_retention_duration", flattenBackupPolicyPostgreSQLDefaultRetentionRuleDuration(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `default_retention_duration`: %+v", err)
			}
			if err := d.Set("retention_rule", flattenBackupPolicyPostgreSQLRetentionRuleArray(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `retention_rule`: %+v", err)
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupPolicyPostgreSQLDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.BackupVaultName, id.ResourceGroup, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return nil
}

func expandBackupPolicyPostgreSQLAzureBackupRuleArray(input []interface{}, taggingCriteria *[]dataprotection.TaggingCriteria) []dataprotection.BasicBasePolicyRule {
	results := make([]dataprotection.BasicBasePolicyRule, 0)
	results = append(results, dataprotection.AzureBackupRule{
		Name:       utils.String("BackupIntervals"),
		ObjectType: dataprotection.ObjectTypeAzureBackupRule,
		DataStore: &dataprotection.DataStoreInfoBase{
			DataStoreType: dataprotection.VaultStore,
			ObjectType:    utils.String("DataStoreInfoBase"),
		},
		BackupParameters: &dataprotection.AzureBackupParams{
			BackupType: utils.String("Full"),
			ObjectType: dataprotection.ObjectTypeAzureBackupParams,
		},
		Trigger: dataprotection.ScheduleBasedTriggerContext{
			Schedule: &dataprotection.BackupSchedule{
				RepeatingTimeIntervals: utils.ExpandStringSlice(input),
			},
			TaggingCriteria: taggingCriteria,
			ObjectType:      dataprotection.ObjectTypeScheduleBasedTriggerContext,
		},
	})

	return results
}

func expandBackupPolicyPostgreSQLAzureRetentionRuleArray(input []interface{}) []dataprotection.BasicBasePolicyRule {
	results := make([]dataprotection.BasicBasePolicyRule, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, dataprotection.AzureRetentionRule{
			Name:       utils.String(v["name"].(string)),
			ObjectType: dataprotection.ObjectTypeAzureRetentionRule,
			IsDefault:  utils.Bool(false),
			Lifecycles: &[]dataprotection.SourceLifeCycle{
				{
					DeleteAfter: dataprotection.AbsoluteDeleteOption{
						Duration:   utils.String(v["duration"].(string)),
						ObjectType: dataprotection.ObjectTypeAbsoluteDeleteOption,
					},
					SourceDataStore: &dataprotection.DataStoreInfoBase{
						DataStoreType: "VaultStore",
						ObjectType:    utils.String("DataStoreInfoBase"),
					},
					TargetDataStoreCopySettings: &[]dataprotection.TargetCopySetting{},
				},
			},
		})
	}
	return results
}

func expandBackupPolicyPostgreSQLDefaultAzureRetentionRule(input interface{}) dataprotection.BasicBasePolicyRule {
	return dataprotection.AzureRetentionRule{
		Name:       utils.String("Default"),
		ObjectType: dataprotection.ObjectTypeAzureRetentionRule,
		IsDefault:  utils.Bool(true),
		Lifecycles: &[]dataprotection.SourceLifeCycle{
			{
				DeleteAfter: dataprotection.AbsoluteDeleteOption{
					Duration:   utils.String(input.(string)),
					ObjectType: dataprotection.ObjectTypeAbsoluteDeleteOption,
				},
				SourceDataStore: &dataprotection.DataStoreInfoBase{
					DataStoreType: "VaultStore",
					ObjectType:    utils.String("DataStoreInfoBase"),
				},
				TargetDataStoreCopySettings: &[]dataprotection.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyPostgreSQLTaggingCriteriaArray(input []interface{}) *[]dataprotection.TaggingCriteria {
	results := []dataprotection.TaggingCriteria{
		{
			Criteria:        nil,
			IsDefault:       utils.Bool(true),
			TaggingPriority: utils.Int64(99),
			TagInfo: &dataprotection.RetentionTag{
				ID:      utils.String("Default_"),
				TagName: utils.String("Default"),
			},
		},
	}
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, dataprotection.TaggingCriteria{
			Criteria:        expandBackupPolicyPostgreSQLCriteriaArray(v["criteria"].([]interface{})),
			IsDefault:       utils.Bool(false),
			TaggingPriority: utils.Int64(int64(v["priority"].(int))),
			TagInfo: &dataprotection.RetentionTag{
				ID:      utils.String(v["name"].(string) + "_"),
				TagName: utils.String(v["name"].(string)),
			},
		})
	}
	return &results
}

func expandBackupPolicyPostgreSQLCriteriaArray(input []interface{}) *[]dataprotection.BasicBackupCriteria {
	results := make([]dataprotection.BasicBackupCriteria, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		var absoluteCriteria []dataprotection.AbsoluteMarker
		if absoluteCriteriaRaw := v["absolute_criteria"].(string); len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []dataprotection.AbsoluteMarker{dataprotection.AbsoluteMarker(absoluteCriteriaRaw)}
		}

		var daysOfWeek []dataprotection.DayOfWeek
		if v["days_of_week"].(*pluginsdk.Set).Len() > 0 {
			daysOfWeek = make([]dataprotection.DayOfWeek, 0)
			for _, value := range v["days_of_week"].(*pluginsdk.Set).List() {
				daysOfWeek = append(daysOfWeek, dataprotection.DayOfWeek(value.(string)))
			}
		}

		var monthsOfYear []dataprotection.Month
		if v["months_of_year"].(*pluginsdk.Set).Len() > 0 {
			monthsOfYear = make([]dataprotection.Month, 0)
			for _, value := range v["months_of_year"].(*pluginsdk.Set).List() {
				monthsOfYear = append(monthsOfYear, dataprotection.Month(value.(string)))
			}
		}

		var weeksOfMonth []dataprotection.WeekNumber
		if v["weeks_of_month"].(*pluginsdk.Set).Len() > 0 {
			weeksOfMonth = make([]dataprotection.WeekNumber, 0)
			for _, value := range v["weeks_of_month"].(*pluginsdk.Set).List() {
				weeksOfMonth = append(weeksOfMonth, dataprotection.WeekNumber(value.(string)))
			}
		}

		var scheduleTimes []date.Time
		if v["scheduled_backup_times"].(*pluginsdk.Set).Len() > 0 {
			scheduleTimes = make([]date.Time, 0)
			for _, value := range v["scheduled_backup_times"].(*pluginsdk.Set).List() {
				t, _ := time.Parse(time.RFC3339, value.(string))
				scheduleTimes = append(scheduleTimes, date.Time{Time: t})
			}
		}
		results = append(results, dataprotection.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			DaysOfMonth:      nil,
			DaysOfTheWeek:    &daysOfWeek,
			MonthsOfYear:     &monthsOfYear,
			ScheduleTimes:    &scheduleTimes,
			WeeksOfTheMonth:  &weeksOfMonth,
			ObjectType:       dataprotection.ObjectTypeScheduleBasedBackupCriteria,
		})
	}
	return &results
}

func flattenBackupPolicyPostgreSQLBackupRuleArray(input *[]dataprotection.BasicBasePolicyRule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	for _, item := range *input {
		if backupRule, ok := item.AsAzureBackupRule(); ok {
			if backupRule.Trigger != nil {
				if scheduleBasedTrigger, ok := backupRule.Trigger.AsScheduleBasedTriggerContext(); ok {
					if scheduleBasedTrigger.Schedule != nil {
						return utils.FlattenStringSlice(scheduleBasedTrigger.Schedule.RepeatingTimeIntervals)
					}
				}
			}
		}
	}
	return make([]interface{}, 0)
}

func flattenBackupPolicyPostgreSQLDefaultRetentionRuleDuration(input *[]dataprotection.BasicBasePolicyRule) interface{} {
	if input == nil {
		return nil
	}

	for _, item := range *input {
		if retentionRule, ok := item.AsAzureRetentionRule(); ok && retentionRule.IsDefault != nil && *retentionRule.IsDefault {
			if retentionRule.Lifecycles != nil && len(*retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (*retentionRule.Lifecycles)[0].DeleteAfter.AsAbsoluteDeleteOption(); ok {
					return *deleteOption.Duration
				}
			}
		}
	}
	return nil
}

func flattenBackupPolicyPostgreSQLRetentionRuleArray(input *[]dataprotection.BasicBasePolicyRule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	var taggingCriterias []dataprotection.TaggingCriteria
	for _, item := range *input {
		if backupRule, ok := item.AsAzureBackupRule(); ok {
			if trigger, ok := backupRule.Trigger.AsScheduleBasedTriggerContext(); ok {
				if trigger.TaggingCriteria != nil {
					taggingCriterias = *trigger.TaggingCriteria
				}
			}
		}
	}

	for _, item := range *input {
		if retentionRule, ok := item.AsAzureRetentionRule(); ok && (retentionRule.IsDefault == nil || !*retentionRule.IsDefault) {
			var name string
			if retentionRule.Name != nil {
				name = *retentionRule.Name
			}
			var taggingPriority int64
			var taggingCriteria []interface{}
			for _, criteria := range taggingCriterias {
				if criteria.TagInfo != nil && criteria.TagInfo.TagName != nil && strings.EqualFold(*criteria.TagInfo.TagName, name) {
					taggingPriority = *criteria.TaggingPriority
					taggingCriteria = flattenBackupPolicyPostgreSQLBackupCriteriaArray(criteria.Criteria)
				}
			}
			var duration string
			if retentionRule.Lifecycles != nil && len(*retentionRule.Lifecycles) > 0 {
				if deleteOption, ok := (*retentionRule.Lifecycles)[0].DeleteAfter.AsAbsoluteDeleteOption(); ok {
					duration = *deleteOption.Duration
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

func flattenBackupPolicyPostgreSQLBackupCriteriaArray(input *[]dataprotection.BasicBackupCriteria) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if criteria, ok := item.AsScheduleBasedBackupCriteria(); ok {
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
				for _, item := range *criteria.ScheduleTimes {
					scheduleTimes = append(scheduleTimes, item.String())
				}
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
