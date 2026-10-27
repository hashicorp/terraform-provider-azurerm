// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_disk

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataProtectionBackupInstanceDiskSchema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"vault_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backupvaultresources.ValidateBackupVaultID,
		},

		"disk_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateManagedDiskID,
		},

		"snapshot_resource_group_name": commonschema.ResourceGroupName(),

		"backup_policy_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: basebackuppolicyresources.ValidateBackupPolicyID,
		},

		"snapshot_subscription_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// vault_id: ID of the parent resource; must share the same subscription ID as this backup instance.
				// Suppress diff if snapshot_subscription_id matches this backup instance's subscription.
				planVaultId := d.Get("vault_id")
				vaultId, err := backupvaultresources.ParseBackupVaultID(planVaultId.(string))
				if err != nil {
					return false
				}

				// oldValue represents the value in the state.
				// newValue represents the value in the config, in case it is "", it means the user doesn't specify it in the config.
				if strings.EqualFold(oldValue, vaultId.SubscriptionId) && newValue == "" {
					return true
				}

				return false
			},
		},

		"protection_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
