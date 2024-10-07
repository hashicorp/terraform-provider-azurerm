// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotectableitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	recoveryServicesValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBackupProtectedFileShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectedFileShareCreateUpdate,
		Read:   resourceBackupProtectedFileShareRead,
		Update: resourceBackupProtectedFileShareCreateUpdate,
		Delete: resourceBackupProtectedFileShareDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ProtectedItemID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(80 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(80 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(80 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: recoveryServicesValidate.RecoveryServicesVaultName,
			},

			"source_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"source_file_share_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareName,
			},

			"backup_policy_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceBackupProtectedFileShareCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	protectedClient := meta.(*clients.Client).RecoveryServices.ProtectedItemsGroupClient
	protectableClient := meta.(*clients.Client).RecoveryServices.ProtectableItemsClient
	protectionContainerClient := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	opResultClient := meta.(*clients.Client).RecoveryServices.ProtectionContainerOperationResultsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)

	vaultName := d.Get("recovery_vault_name").(string)
	storageAccountID := d.Get("source_storage_account_id").(string)
	fileShareName := d.Get("source_file_share_name").(string)
	policyID := d.Get("backup_policy_id").(string)

	vaultId := backupprotectableitems.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string))

	// get storage account name from id
	parsedStorageAccountID, err := commonids.ParseStorageAccountID(storageAccountID)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_storage_account_id '%s': %+v", storageAccountID, err)
	}

	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroupName, parsedStorageAccountID.StorageAccountName)
	log.Printf("[DEBUG] creating/updating Recovery Service Protected File Share %q (Container Name %q)", fileShareName, containerName)

	// the fileshare has a user defined name, but its system name (fileShareSystemName) is only known to Azure Backup
	fileShareSystemName := ""
	// @aristosvo: preferred filter would be like below but the 'and' expression seems to fail
	//   filter := fmt.Sprintf("backupManagementType eq 'AzureStorage' and friendlyName eq '%s'", fileShareName)
	// this means which means we have to do it client side and loop over backupProtectedItems en backupProtectableItems until share is found
	filter := "backupManagementType eq 'AzureStorage'"

	protectionContainerId := protectioncontainers.NewProtectionContainerID(subscriptionId, resourceGroup, vaultName, "Azure", containerName)

	// There is an issue https://github.com/hashicorp/terraform-provider-azurerm/issues/11184 (When a new file share is added to an existing storage account,
	// it cannot be listed by Backup Protectable Items - List API after the storage account is registered with a RSV).
	// After confirming with the service team, whenever new file shares are added, we need to run an 'inquire' API. but inquiry APIs are long running APIs and hence can't be included in GET API's (Backup Protectable Items - List) response.
	// Therefore, add 'inquire' API to inquire all unprotected files shares under a storage account to fix this usecase.
	respContainer, err := protectionContainerClient.Inquire(ctx, protectionContainerId, protectioncontainers.InquireOperationOptions{Filter: pointer.To(filter)})
	if err != nil {
		return fmt.Errorf("inquire all unprotected files shares for %s: %+v", parsedStorageAccountID, err)
	}

	// TODO: @tombuildsstuff: this manual LRO is not needed and should be removed - the existing Azure SDK has logic to handle this as does hashicorp/go-azure-sdk
	// therefore we should not be invoking the Future by hand, there's already logic to do that for us:
	// When using `Azure/go-autorest`: https://github.com/hashicorp/go-azure-helpers/blob/8045457c83689876d4c63fecebd4753925ea73ab/polling/poller.go#L30
	// When using `hashicorp/go-azure-sdk`: https://github.com/hashicorp/go-azure-sdk/blob/02376e1c45321faa0a561e0c9b43463f1acbc3bb/sdk/client/resourcemanager/poller.go#L16
	locationURL, err := respContainer.HttpResponse.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("inquire all unprotected files shares %q (Vault %q): Location header missing or empty", containerName, vaultName)
	}

	opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["operationResults"]

	// `inquire` API is an async operation and the results should be tracked using location header or Azure-async-url.
	//  The Azure-AsyncOperation is not included in swagger, so call location (https://docs.microsoft.com/en-us/rest/api/backup/protection-container-operation-results/get)
	//  to wait the operation successfully completes.
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 10 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"202"},
		Target:     []string{"200", "204"},
		Refresh:    protectionContainerOperationResultsRefreshFunc(ctx, opResultClient, vaultName, resourceGroup, containerName, operationID),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := state.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Recovery Service Protection Container operation %q (Vault %q in Resource Group %q): %+v", operationID, vaultName, resourceGroup, err)
	}

	backupProtectableItemsListOptions := backupprotectableitems.ListOperationOptions{
		Filter: pointer.To(filter),
	}
	backupProtectableItemsResponse, err := protectableClient.List(ctx, vaultId, backupProtectableItemsListOptions)
	if err != nil {
		return fmt.Errorf("checking for protectable fileshares in Recovery Service Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}

	if backupProtectableItemsResponse.Model != nil {
		for _, protectableItem := range *backupProtectableItemsResponse.Model {
			if *protectableItem.Name == "" || protectableItem.Properties == nil {
				continue
			}
			azureFileShareProtectableItem, check := protectableItem.Properties.(backupprotectableitems.AzureFileShareProtectableItem)

			// check if protected item has the same fileshare name and is from the same storage account
			if check && *azureFileShareProtectableItem.FriendlyName == fileShareName && *azureFileShareProtectableItem.ParentContainerFriendlyName == parsedStorageAccountID.StorageAccountName {
				fileShareSystemName = *protectableItem.Name
				break
			}
		}
	}

	// fileShareSystemName not found? Check if already protected by this vault!
	if fileShareSystemName == "" {
		vaultId := backupprotecteditems.NewVaultID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string))
		backupProtectedItemsResponse, err := protectedClient.List(ctx, vaultId, backupprotecteditems.ListOperationOptions{})
		if err != nil {
			return fmt.Errorf("checking for protected fileshares in Recovery Service Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
		}

		if model := backupProtectedItemsResponse.Model; model != nil {
			for _, protectedItem := range *model {
				if *protectedItem.Name == "" || protectedItem.Properties == nil {
					continue
				}
				azureFileShareProtectedItem, check := protectedItem.Properties.(backupprotecteditems.AzureFileshareProtectedItem)

				// check if protected item has the same fileshare name and is from the same storage account
				if check && *azureFileShareProtectedItem.FriendlyName == fileShareName && strings.EqualFold(*azureFileShareProtectedItem.SourceResourceId, storageAccountID) {
					fileShareSystemName = *protectedItem.Name
					break
				}
			}
		}
	}
	if fileShareSystemName == "" {
		return fmt.Errorf("[ERROR] fileshare '%s' not found in protectable or protected fileshares, make sure Storage Account %q is registered with Recovery Service Vault %q (Resource Group %q)", fileShareName, parsedStorageAccountID.StorageAccountName, vaultName, resourceGroup)
	}

	id := protecteditems.NewProtectedItemID(subscriptionId, d.Get("resource_group_name").(string), d.Get("recovery_vault_name").(string), "Azure", containerName, fileShareSystemName)

	if d.IsNewResource() {
		existing, err2 := client.Get(ctx, id, protecteditems.GetOperationOptions{})
		if err2 != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Recovery Service Protected File Share %q (Resource Group %q): %+v", fileShareName, resourceGroup, err2)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_backup_protected_file_share", id.ID())
		}
	}

	item := protecteditems.ProtectedItemResource{
		Properties: &protecteditems.AzureFileshareProtectedItem{
			PolicyId:         &policyID,
			WorkloadType:     pointer.To(protecteditems.DataSourceTypeAzureFileShare),
			SourceResourceId: utils.String(storageAccountID),
			FriendlyName:     utils.String(fileShareName),
		},
	}

	resp, err := client.CreateOrUpdate(ctx, id, item)
	if err != nil {
		return fmt.Errorf("creating/updating Recovery Service Protected File Share %q (Resource Group %q): %+v", fileShareName, resourceGroup, err)
	}

	locationURL, err = resp.HttpResponse.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("creating/updating Azure File Share backup item %q (Vault %q): Location header missing or empty", containerName, vaultName)
	}

	opResourceID = handleAzureSdkForGoBug2824(locationURL.String())

	parsedLocation, err = azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID = parsedLocation.Path["operationResults"]

	if _, err := resourceBackupProtectedFileShareWaitForOperation(ctx, opClient, vaultName, resourceGroup, operationID, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceBackupProtectedFileShareRead(d, meta)
}

func resourceBackupProtectedFileShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading %s", *id)

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
					sourceResourceID := strings.Replace(*item.SourceResourceId, "Microsoft.storage", "Microsoft.Storage", 1) // The SDK is returning inconsistent capitalization
					d.Set("source_storage_account_id", sourceResourceID)
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

func resourceBackupProtectedFileShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("issuing delete request for %s: %+v", id, err)
		}
	}

	locationURL, err := resp.HttpResponse.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("deleting Azure File Share backups item %s (Vault %s): Location header missing or empty", id.ProtectionContainerName, id.VaultName)
	}

	opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["backupOperationResults"] // This is different for create and delete requests ¯\_(ツ)_/¯

	if _, err := resourceBackupProtectedFileShareWaitForOperation(ctx, opClient, id.VaultName, id.ResourceGroupName, operationID, d); err != nil {
		return err
	}

	return nil
}

// nolint unused - linter mistakenly things this function isn't used?
func resourceBackupProtectedFileShareWaitForOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string, d *pluginsdk.ResourceData) (backup.OperationStatus, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 10 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"InProgress"},
		Target:     []string{"Succeeded"},
		Refresh:    resourceBackupProtectedFileShareCheckOperation(ctx, client, vaultName, resourceGroup, operationID),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	log.Printf("[DEBUG] Waiting for backup operation %s (Vault %s) to complete", operationID, vaultName)
	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		return resp.(backup.OperationStatus), err
	}
	return resp.(backup.OperationStatus), nil
}

func resourceBackupProtectedFileShareCheckOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, operationID)
		if err != nil {
			return resp, "Error", fmt.Errorf("making Read request on Recovery Service Protection Container operation %q (Vault %q in Resource Group %q): %+v", operationID, vaultName, resourceGroup, err)
		}

		if opErr := resp.Error; opErr != nil {
			errMsg := "No upstream error message"
			if opErr.Message != nil {
				errMsg = *opErr.Message
			}
			err = fmt.Errorf("Azure Backup operation status failed with status %q (Vault %q Resource Group %q Operation ID %q): %+v", resp.Status, vaultName, resourceGroup, operationID, errMsg)
		}

		log.Printf("[DEBUG] Backup operation %s status is %s", operationID, string(resp.Status))
		return resp, string(resp.Status), err
	}
}

func protectionContainerOperationResultsRefreshFunc(ctx context.Context, client *backup.ProtectionContainerOperationResultsClient, vaultName, resourceGroup, containerName string, operationID string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, operationID)
		if err != nil {
			return nil, "Error", fmt.Errorf("making Read request on Recovery Service Protection Container operation %q (Vault %q in Resource Group %q): %+v", operationID, vaultName, resourceGroup, err)
		}

		return resp, strconv.Itoa(resp.StatusCode), err
	}
}
