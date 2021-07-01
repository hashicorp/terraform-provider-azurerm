package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	helperValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataProtectionBackupPolicyDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupPolicyDiskCreate,
		Read:   resourceDataProtectionBackupPolicyDiskRead,
		Delete: resourceDataProtectionBackupPolicyDiskDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.BackupPolicyID(id)
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
				ValidateFunc: validate.BackupVaultID,
			},

			"backup_repeating_time_intervals": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"default_retention_duration": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helperValidate.ISO8601Duration,
			},

			"retention_rule": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"duration": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: helperValidate.ISO8601Duration,
						},

						"criteria": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"absolute_criteria": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(dataprotection.FirstOfDay),
											string(dataprotection.FirstOfWeek),
										}, false),
									},
								},
							},
						},

						"priority": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}
func resourceDataProtectionBackupPolicyDiskCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := parse.BackupVaultID(d.Get("vault_id").(string))
	id := parse.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroup, vaultId.Name, name)

	existing, err := client.Get(ctx, id.BackupVaultName, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_disk", id.ID())
	}

	taggingCriteria := expandBackupPolicyDiskTaggingCriteriaArray(d.Get("retention_rule").([]interface{}))
	policyRules := make([]dataprotection.BasicBasePolicyRule, 0)
	policyRules = append(policyRules, expandBackupPolicyDiskAzureBackupRuleArray(d.Get("backup_repeating_time_intervals").([]interface{}), taggingCriteria)...)
	policyRules = append(policyRules, expandBackupPolicyDiskDefaultAzureRetentionRule(d.Get("default_retention_duration")))
	policyRules = append(policyRules, expandBackupPolicyDiskAzureRetentionRuleArray(d.Get("retention_rule").([]interface{}))...)
	parameters := dataprotection.BaseBackupPolicyResource{
		Properties: &dataprotection.BackupPolicy{
			PolicyRules:     &policyRules,
			DatasourceTypes: &[]string{"Microsoft.Compute/disks"},
			ObjectType:      dataprotection.ObjectTypeBackupPolicy,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.BackupVaultName, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupPolicyDiskRead(d, meta)
}

func resourceDataProtectionBackupPolicyDiskRead(d *schema.ResourceData, meta interface{}) error {
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
	vaultId := parse.NewBackupVaultID(id.SubscriptionId, id.ResourceGroup, id.BackupVaultName)
	d.Set("name", id.Name)
	d.Set("vault_id", vaultId.ID())
	if resp.Properties != nil {
		if props, ok := resp.Properties.AsBackupPolicy(); ok {
			if err := d.Set("backup_repeating_time_intervals", flattenBackupPolicyDiskBackupRuleArray(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `backup_repeating_time_intervals`: %+v", err)
			}
			if err := d.Set("default_retention_duration", flattenBackupPolicyDiskDefaultRetentionRuleDuration(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `default_retention_duration`: %+v", err)
			}
			if err := d.Set("retention_rule", flattenBackupPolicyDiskRetentionRuleArray(props.PolicyRules)); err != nil {
				return fmt.Errorf("setting `retention_rule`: %+v", err)
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupPolicyDiskDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandBackupPolicyDiskAzureBackupRuleArray(input []interface{}, taggingCriteria *[]dataprotection.TaggingCriteria) []dataprotection.BasicBasePolicyRule {
	results := make([]dataprotection.BasicBasePolicyRule, 0)

	results = append(results, dataprotection.AzureBackupRule{
		Name:       utils.String("BackupIntervals"),
		ObjectType: dataprotection.ObjectTypeAzureBackupRule,
		DataStore: &dataprotection.DataStoreInfoBase{
			DataStoreType: dataprotection.OperationalStore,
			ObjectType:    utils.String("DataStoreInfoBase"),
		},
		BackupParameters: &dataprotection.AzureBackupParams{
			BackupType: utils.String("Incremental"),
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

func expandBackupPolicyDiskAzureRetentionRuleArray(input []interface{}) []dataprotection.BasicBasePolicyRule {
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
						DataStoreType: "OperationalStore",
						ObjectType:    utils.String("DataStoreInfoBase"),
					},
					TargetDataStoreCopySettings: &[]dataprotection.TargetCopySetting{},
				},
			},
		})
	}
	return results
}

func expandBackupPolicyDiskDefaultAzureRetentionRule(input interface{}) dataprotection.BasicBasePolicyRule {
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
					DataStoreType: "OperationalStore",
					ObjectType:    utils.String("DataStoreInfoBase"),
				},
				TargetDataStoreCopySettings: &[]dataprotection.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyDiskTaggingCriteriaArray(input []interface{}) *[]dataprotection.TaggingCriteria {
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
			Criteria:        expandBackupPolicyDiskCriteriaArray(v["criteria"].([]interface{})),
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

func expandBackupPolicyDiskCriteriaArray(input []interface{}) *[]dataprotection.BasicBackupCriteria {
	results := make([]dataprotection.BasicBackupCriteria, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		var absoluteCriteria []dataprotection.AbsoluteMarker
		if absoluteCriteriaRaw := v["absolute_criteria"].(string); len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []dataprotection.AbsoluteMarker{dataprotection.AbsoluteMarker(absoluteCriteriaRaw)}
		}
		results = append(results, dataprotection.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
			ObjectType:       dataprotection.ObjectTypeScheduleBasedBackupCriteria,
		})
	}
	return &results
}

func flattenBackupPolicyDiskBackupRuleArray(input *[]dataprotection.BasicBasePolicyRule) []interface{} {
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

func flattenBackupPolicyDiskDefaultRetentionRuleDuration(input *[]dataprotection.BasicBasePolicyRule) interface{} {
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

func flattenBackupPolicyDiskRetentionRuleArray(input *[]dataprotection.BasicBasePolicyRule) []interface{} {
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
					taggingCriteria = flattenBackupPolicyDiskBackupCriteriaArray(criteria.Criteria)
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

func flattenBackupPolicyDiskBackupCriteriaArray(input *[]dataprotection.BasicBackupCriteria) []interface{} {
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

			results = append(results, map[string]interface{}{
				"absolute_criteria": absoluteCriteria,
			})
		}
	}
	return results
}
