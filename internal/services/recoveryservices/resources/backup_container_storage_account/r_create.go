// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_container_storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionContainerStorageAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	opStatusClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountID := d.Get("storage_account_id").(string)
	parsedStorageAccountID, err := commonids.ParseStorageAccountID(storageAccountID)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse storage_account_id '%s': %+v", storageAccountID, err)
	}

	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroupName, parsedStorageAccountID.StorageAccountName)

	id := protectioncontainers.NewProtectionContainerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), "Azure", containerName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_backup_protection_container_storage", id.ID())
		}
	}

	parameters := protectioncontainers.ProtectionContainerResource{
		Properties: &protectioncontainers.AzureStorageContainer{
			SourceResourceId:     &storageAccountID,
			FriendlyName:         &parsedStorageAccountID.StorageAccountName,
			BackupManagementType: pointer.To(protectioncontainers.BackupManagementTypeAzureStorage),
		},
	}

	resp, err := client.Register(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("registering %s: %+v", id, err)
	}

	locationURL, err := resp.HttpResponse.Location() // Operation ID found in the Location header
	if locationURL == nil || err != nil {
		return fmt.Errorf("unable to determine operation URL for %s: Location header missing or empty", id)
	}

	opResourceID := helpers.HandleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}

	operationID := parsedLocation.Path["operationResults"]
	if err = resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx, opStatusClient, id.VaultName, id.ResourceGroupName, operationID, d); err != nil {
		return err
	}

	d.SetId(helpers.HandleAzureSdkForGoBug2824(id.ID()))

	return resourceBackupProtectionContainerStorageAccountRead(d, meta)
}
