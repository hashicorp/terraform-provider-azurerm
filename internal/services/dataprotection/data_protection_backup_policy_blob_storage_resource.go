// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/backuppolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataProtectionBackupPolicyBlobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataProtectionBackupPolicyBlobStorageCreate,
		Read:   resourceDataProtectionBackupPolicyBlobStorageRead,
		Delete: resourceDataProtectionBackupPolicyBlobStorageDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
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

			"retention_duration": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: helperValidate.ISO8601Duration,
			},
		},
	}
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

	parameters := backuppolicies.BaseBackupPolicyResource{
		Properties: &backuppolicies.BackupPolicy{
			PolicyRules: []backuppolicies.BasePolicyRule{
				backuppolicies.AzureRetentionRule{
					Name:      "Default",
					IsDefault: utils.Bool(true),
					Lifecycles: []backuppolicies.SourceLifeCycle{
						{
							DeleteAfter: backuppolicies.AbsoluteDeleteOption{
								Duration: d.Get("retention_duration").(string),
							},
							SourceDataStore: backuppolicies.DataStoreInfoBase{
								DataStoreType: "OperationalStore",
								ObjectType:    "DataStoreInfoBase",
							},
							TargetDataStoreCopySettings: &[]backuppolicies.TargetCopySetting{},
						},
					},
				},
			},
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
				if err := d.Set("retention_duration", flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(props.PolicyRules)); err != nil {
					return fmt.Errorf("setting `default_retention_duration`: %+v", err)
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

func flattenBackupPolicyBlobStorageDefaultRetentionRuleDuration(input []backuppolicies.BasePolicyRule) interface{} {
	if input == nil {
		return nil
	}

	for _, item := range input {
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
