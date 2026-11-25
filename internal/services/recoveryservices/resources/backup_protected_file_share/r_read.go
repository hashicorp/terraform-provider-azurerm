// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_file_share

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectedFileShareRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			if item, ok := properties.(protecteditems.AzureFileshareProtectedItem); ok {
				if item.SourceResourceId != nil {
					d.Set("source_storage_account_id", strings.Replace(*item.SourceResourceId, "Microsoft.storage", "Microsoft.Storage", 1))
				}
				d.Set("source_file_share_name", item.FriendlyName)

				if v := item.PolicyId; v != nil {
					d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
				}
			}
		}
	}

	return nil
}

func protectionContainerOperationResultsRefreshFunc(ctx context.Context, client *backup.ProtectionContainerOperationResultsClient, vaultName, resourceGroup, containerName string, operationID string) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, operationID)
		if err != nil {
			return nil, "Error", fmt.Errorf("making Read request on Recovery Service Protection Container operation %q (Vault %q in Resource Group %q): %+v", operationID, vaultName, resourceGroup, err)
		}

		return resp, strconv.Itoa(resp.StatusCode), err
	}
}
