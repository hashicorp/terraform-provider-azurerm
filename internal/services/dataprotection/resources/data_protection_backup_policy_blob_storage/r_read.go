// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_blob_storage

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupPolicyBlobStorageRead(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupPolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := basebackuppolicyresources.ParseBackupPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.BackupPoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	vaultId := basebackuppolicyresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupPolicyName)
	d.Set("vault_id", vaultId.ID())
	if resp.Model != nil {
		if resp.Model.Properties != nil {
			if props, ok := resp.Model.Properties.(basebackuppolicyresources.BackupPolicy); ok {
				if err := d.Set("backup_repeating_time_intervals", flattenBackupPolicyBlobStorageVaultBackupRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `backup_repeating_time_intervals`: %+v", err)
				}
				if err := d.Set("operational_default_retention_duration", flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(props.PolicyRules, basebackuppolicyresources.DataStoreTypesOperationalStore)); err != nil {
					return fmt.Errorf("setting `operational_default_retention_duration`: %+v", err)
				}
				d.Set("time_zone", flattenBackupPolicyBlobStorageVaultBackupTimeZone(&props.PolicyRules))
				if err := d.Set("vault_default_retention_duration", flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(props.PolicyRules, basebackuppolicyresources.DataStoreTypesVaultStore)); err != nil {
					return fmt.Errorf("setting `vault_default_retention_duration`: %+v", err)
				}
				if err := d.Set("retention_rule", flattenBackupPolicyBlobStorageRetentionRuleArray(&props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `retention_rule`: %+v", err)
				}
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

const (
	blobStorageOperationalDefaultRetentionRuleName = "Default_OperationalStore"
	blobStorageVaultDefaultRetentionRuleName       = "Default"
)
