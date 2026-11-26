// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_disk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupPolicyDiskCreate(d *schema.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	vaultId, _ := basebackuppolicyresources.ParseBackupVaultID(d.Get("vault_id").(string))
	id := basebackuppolicyresources.NewBackupPolicyID(subscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.BackupPoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing DataProtection BackupPolicy (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_protection_backup_policy_disk", id.ID())
		}
	}

	taggingCriteria := expandBackupPolicyDiskTaggingCriteriaArray(d.Get("retention_rule").([]any))
	policyRules := make([]basebackuppolicyresources.BasePolicyRule, 0)
	policyRules = append(policyRules, expandBackupPolicyDiskAzureBackupRuleArray(d.Get("backup_repeating_time_intervals").([]any), d.Get("time_zone").(string), taggingCriteria)...)
	policyRules = append(policyRules, expandBackupPolicyDiskDefaultAzureRetentionRule(d.Get("default_retention_duration")))
	policyRules = append(policyRules, expandBackupPolicyDiskAzureRetentionRuleArray(d.Get("retention_rule").([]any))...)
	parameters := basebackuppolicyresources.BaseBackupPolicyResource{
		Properties: &basebackuppolicyresources.BackupPolicy{
			PolicyRules:     policyRules,
			DatasourceTypes: []string{"Microsoft.Compute/disks"},
		},
	}

	if _, err := client.BackupPoliciesCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}
	return resourceDataProtectionBackupPolicyDiskRead(d, meta)
}
