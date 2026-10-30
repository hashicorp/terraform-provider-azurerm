// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_container_storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionContainerStorageAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectioncontainers.ParseProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Unregister(ctx, *id)
	if err != nil {
		return fmt.Errorf("deregistering %s: %+v", id, err)
	}

	locationURL, err := resp.HttpResponse.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("unregistering backup protection container %s : Location header missing or empty", id.String())
	}

	opResourceID := helpers.HandleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["backupOperationResults"]

	if err = resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx, opClient, id.VaultName, id.ResourceGroupName, operationID, d); err != nil {
		return err
	}

	return nil
}
