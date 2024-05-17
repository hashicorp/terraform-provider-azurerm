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

func resourceDataProtectionBackupPolicyDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupPolicyDiskCreate,
		Read:   resourceDataProtectionBackupPolicyDiskRead,
		Delete: resourceDataProtectionBackupPolicyDiskDelete,

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
											string(backuppolicies.AbsoluteMarkerFirstOfDay),
											string(backuppolicies.AbsoluteMarkerFirstOfWeek),
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

			"time_zone": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
	vaultId, _ := backuppolicies.ParseBackupVaultID(d.Get("vault_id").(string))
	id := backuppolicies.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_disk", id.ID())
	}

	taggingCriteria := expandBackupPolicyDiskTaggingCriteriaArray(d.Get("retention_rule").([]interface{}))
	policyRules := make([]backuppolicies.BasePolicyRule, 0)
	policyRules = append(policyRules, expandBackupPolicyDiskAzureBackupRuleArray(d.Get("backup_repeating_time_intervals").([]interface{}), d.Get("time_zone").(string), taggingCriteria)...)
	policyRules = append(policyRules, expandBackupPolicyDiskDefaultAzureRetentionRule(d.Get("default_retention_duration")))
	policyRules = append(policyRules, expandBackupPolicyDiskAzureRetentionRuleArray(d.Get("retention_rule").([]interface{}))...)
	parameters := backuppolicies.BaseBackupPolicyResource{
		Properties: &backuppolicies.BackupPolicy{
			PolicyRules:     policyRules,
			DatasourceTypes: []string{"Microsoft.Compute/disks"},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating DataProtection BackupPolicy (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupPolicyDiskRead(d, meta)
}

func resourceDataProtectionBackupPolicyDiskRead(d *schema.ResourceData, meta interface{}) error {
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
				if err := d.Set("backup_repeating_time_intervals", flattenBackupPolicyDiskBackupRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `backup_repeating_time_intervals`: %+v", err)
				}
				if err := d.Set("default_retention_duration", flattenBackupPolicyDiskDefaultRetentionRuleDuration(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `default_retention_duration`: %+v", err)
				}
				if err := d.Set("retention_rule", flattenBackupPolicyDiskRetentionRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `retention_rule`: %+v", err)
				}
				d.Set("time_zone", flattenBackupPolicyDiskBackupTimeZone(&props.PolicyRules))
			}
		}
	}
	return nil
}

func resourceDataProtectionBackupPolicyDiskDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandBackupPolicyDiskAzureBackupRuleArray(input []interface{}, timeZone string, taggingCriteria *[]backuppolicies.TaggingCriteria) []backuppolicies.BasePolicyRule {
	results := make([]backuppolicies.BasePolicyRule, 0)

	results = append(results, backuppolicies.AzureBackupRule{
		Name: "BackupIntervals",
		DataStore: backuppolicies.DataStoreInfoBase{
			DataStoreType: backuppolicies.DataStoreTypesOperationalStore,
			ObjectType:    "DataStoreInfoBase",
		},
		BackupParameters: &backuppolicies.AzureBackupParams{
			BackupType: "Incremental",
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

func expandBackupPolicyDiskAzureRetentionRuleArray(input []interface{}) []backuppolicies.BasePolicyRule {
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
						DataStoreType: "OperationalStore",
						ObjectType:    "DataStoreInfoBase",
					},
					TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
				},
			},
		})
	}
	return results
}

func expandBackupPolicyDiskDefaultAzureRetentionRule(input interface{}) backuppolicies.BasePolicyRule {
	return backuppolicies.AzureRetentionRule{
		Name:      "Default",
		IsDefault: utils.Bool(true),
		Lifecycles: []backuppolicies.SourceLifeCycle{
			{
				DeleteAfter: backuppolicies.AbsoluteDeleteOption{
					Duration: input.(string),
				},
				SourceDataStore: backuppolicies.DataStoreInfoBase{
					DataStoreType: "OperationalStore",
					ObjectType:    "DataStoreInfoBase",
				},
				TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
			},
		},
	}
}

func expandBackupPolicyDiskTaggingCriteriaArray(input []interface{}) *[]backuppolicies.TaggingCriteria {
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
		results = append(results, backuppolicies.TaggingCriteria{
			Criteria:        expandBackupPolicyDiskCriteriaArray(v["criteria"].([]interface{})),
			IsDefault:       false,
			TaggingPriority: int64(v["priority"].(int)),
			TagInfo: backuppolicies.RetentionTag{
				Id:      utils.String(v["name"].(string) + "_"),
				TagName: v["name"].(string),
			},
		})
	}
	return &results
}

func expandBackupPolicyDiskCriteriaArray(input []interface{}) *[]backuppolicies.BackupCriteria {
	results := make([]backuppolicies.BackupCriteria, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		var absoluteCriteria []backuppolicies.AbsoluteMarker
		if absoluteCriteriaRaw := v["absolute_criteria"].(string); len(absoluteCriteriaRaw) > 0 {
			absoluteCriteria = []backuppolicies.AbsoluteMarker{backuppolicies.AbsoluteMarker(absoluteCriteriaRaw)}
		}
		results = append(results, backuppolicies.ScheduleBasedBackupCriteria{
			AbsoluteCriteria: &absoluteCriteria,
		})
	}
	return &results
}

func flattenBackupPolicyDiskBackupRuleArray(input *[]backuppolicies.BasePolicyRule) []interface{} {
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

func flattenBackupPolicyDiskBackupTimeZone(input *[]backuppolicies.BasePolicyRule) string {
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

func flattenBackupPolicyDiskDefaultRetentionRuleDuration(input *[]backuppolicies.BasePolicyRule) interface{} {
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

func flattenBackupPolicyDiskRetentionRuleArray(input *[]backuppolicies.BasePolicyRule) []interface{} {
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
					taggingCriteria = flattenBackupPolicyDiskBackupCriteriaArray(criteria.Criteria)
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

func flattenBackupPolicyDiskBackupCriteriaArray(input *[]backuppolicies.BackupCriteria) []interface{} {
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

			results = append(results, map[string]interface{}{
				"absolute_criteria": absoluteCriteria,
			})
		}
	}
	return results
}
