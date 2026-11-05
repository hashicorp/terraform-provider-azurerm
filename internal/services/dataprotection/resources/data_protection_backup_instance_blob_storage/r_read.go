// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_blob_storage

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataProtectionBackupInstanceBlobStorageRead(d *schema.ResourceData, meta any) error {
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
			d.Set("storage_account_id", props.DataSourceInfo.ResourceID)
			d.Set("location", props.DataSourceInfo.ResourceLocation)
			d.Set("backup_policy_id", props.PolicyInfo.PolicyId)
			d.Set("protection_state", pointer.FromEnum(props.CurrentProtectionState))
			if policyParas := props.PolicyInfo.PolicyParameters; policyParas != nil {
				if dataStoreParas := policyParas.BackupDatasourceParametersList; dataStoreParas != nil {
					if dsp := pointer.From(dataStoreParas); len(dsp) > 0 {
						if parameter, ok := dsp[0].(backupinstanceresources.BlobBackupDatasourceParameters); ok {
							if err := d.Set("storage_account_container_names", pluginsdk.FlattenSlice(&parameter.ContainersList)); err != nil {
								return fmt.Errorf("setting `storage_account_container_names`: %+v", err)
							}
						}
					}
				}
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}
