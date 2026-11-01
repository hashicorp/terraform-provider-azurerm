// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_disk

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupInstanceDiskRead(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataProtection.BackupInstanceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := backupinstanceresources.ParseBackupInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.BackupInstancesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] dataprotection %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupInstance (%q): %+v", id, err)
	}
	vaultId := backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	d.Set("name", id.BackupInstanceName)
	d.Set("vault_id", vaultId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("disk_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)

			d.Set("protection_state", pointer.FromEnum(props.CurrentProtectionState))
			d.Set("backup_policy_id", props.PolicyInfo.PolicyId)
			if props.PolicyInfo.PolicyParameters != nil && props.PolicyInfo.PolicyParameters.DataStoreParametersList != nil && len(*props.PolicyInfo.PolicyParameters.DataStoreParametersList) > 0 {
				parameter := (*props.PolicyInfo.PolicyParameters.DataStoreParametersList)[0].(backupinstanceresources.AzureOperationalStoreParameters)

				if parameter.ResourceGroupId != nil {
					resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(*parameter.ResourceGroupId)
					if err != nil {
						return err
					}
					d.Set("snapshot_resource_group_name", resourceGroupId.ResourceGroupName)
					d.Set("snapshot_subscription_id", resourceGroupId.SubscriptionId)
				}
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
