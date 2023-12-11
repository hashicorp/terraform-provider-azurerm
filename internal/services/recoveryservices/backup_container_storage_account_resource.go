// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectionContainerStorageAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectionContainerStorageAccountCreate,
		Read:   resourceBackupProtectionContainerStorageAccountRead,
		Delete: resourceBackupProtectionContainerStorageAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protectioncontainers.ParseProtectionContainerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceBackupProtectionContainerStorageAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
	if d.IsNewResource() {
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

	opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}

	operationID := parsedLocation.Path["operationResults"]
	if err = resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx, opStatusClient, id.VaultName, id.ResourceGroupName, operationID, d); err != nil {
		return err
	}

	d.SetId(handleAzureSdkForGoBug2824(id.ID()))

	return resourceBackupProtectionContainerStorageAccountRead(d, meta)
}

func resourceBackupProtectionContainerStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protectioncontainers.ParseProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on backup protection container %s : %+v", id.String(), err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	if model := resp.Model; model != nil {
		if properties, ok := model.Properties.(protectioncontainers.AzureStorageContainer); ok {
			d.Set("storage_account_id", properties.SourceResourceId)
		}
	}

	return nil
}

func resourceBackupProtectionContainerStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

	opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

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

// nolint unused - linter mistakenly things this function isn't used?
func resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout:                10 * time.Second,
		Delay:                     10 * time.Second,
		Pending:                   []string{"InProgress"},
		Target:                    []string{"Succeeded"},
		Refresh:                   resourceBackupProtectionContainerStorageAccountCheckOperation(ctx, client, vaultName, resourceGroup, operationID),
		ContinuousTargetOccurence: 5, // Without this buffer, file share backups and storage account deletions may fail if performed immediately after creating/destroying the container
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	log.Printf("[DEBUG] Waiting for backup container operation %q (Vault %q) to complete", operationID, vaultName)
	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func resourceBackupProtectionContainerStorageAccountCheckOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string) pluginsdk.StateRefreshFunc {
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
			err = fmt.Errorf("Recovery Service Protection Container operation status failed with status %q (Vault %q Resource Group %q Operation ID %q): %+v", resp.Status, vaultName, resourceGroup, operationID, errMsg)
		}

		return resp, string(resp.Status), err
	}
}
