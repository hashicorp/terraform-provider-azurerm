package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-07-01/backup"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBackupProtectionContainerStorageAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBackupProtectionContainerStorageAccountCreate,
		Read:   resourceBackupProtectionContainerStorageAccountRead,
		Update: nil,
		Delete: resourceBackupProtectionContainerStorageAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ProtectionContainerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

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
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	storageAccountID := d.Get("storage_account_id").(string)

	parsedStorageAccountID, err := storageParse.StorageAccountID(storageAccountID)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse storage_account_id '%s': %+v", storageAccountID, err)
	}

	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, parsedStorageAccountID.Name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, vaultName, resGroup, "Azure", containerName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing recovery services protection container %s (Vault %s): %+v", containerName, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_protection_container_storage", handleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	parameters := backup.ProtectionContainerResource{
		Properties: &backup.AzureStorageContainer{
			SourceResourceID:     &storageAccountID,
			FriendlyName:         &parsedStorageAccountID.Name,
			BackupManagementType: backup.ManagementTypeAzureStorage,
			ContainerType:        backup.ContainerTypeBasicProtectionContainerContainerTypeStorageContainer,
		},
	}

	resp, err := client.Register(ctx, vaultName, resGroup, "Azure", containerName, parameters)
	if err != nil {
		return fmt.Errorf("registering backup protection container %s (Vault %s): %+v", containerName, vaultName, err)
	}

	locationURL, err := resp.Response.Location() // Operation ID found in the Location header
	if locationURL == nil || err != nil {
		return fmt.Errorf("Unable to determine operation URL for protection container registration status for %s. (Vault %s): Location header missing or empty", containerName, vaultName)
	}

	opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}

	operationID := parsedLocation.Path["operationResults"]
	if _, err = resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx, opStatusClient, vaultName, resGroup, operationID, d); err != nil {
		return err
	}

	resp, err = client.Get(ctx, vaultName, resGroup, "Azure", containerName)
	if err != nil {
		return fmt.Errorf("retrieving site recovery protection container %s (Vault %s): %+v", containerName, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceBackupProtectionContainerStorageAccountRead(d, meta)
}

func resourceBackupProtectionContainerStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.VaultName, id.ResourceGroup, id.BackupFabricName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on backup protection container %s : %+v", id.String(), err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)

	if properties, ok := resp.Properties.AsAzureStorageContainer(); ok && properties != nil {
		d.Set("storage_account_id", properties.SourceResourceID)
	}

	return nil
}

func resourceBackupProtectionContainerStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ProtectionContainerID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Unregister(ctx, id.VaultName, id.ResourceGroup, id.BackupFabricName, id.Name)
	if err != nil {
		return fmt.Errorf("deregistering backup protection container %s (Vault %s): %+v", id.Name, id.VaultName, err)
	}

	locationURL, err := resp.Response.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("unregistering backup protection container %s : Location header missing or empty", id.String())
	}

	opResourceID := handleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["backupOperationResults"]

	if _, err = resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx, opClient, id.VaultName, id.ResourceGroup, operationID, d); err != nil {
		return err
	}

	return nil
}

// nolint unused - linter mistakenly things this function isn't used?
func resourceBackupProtectionContainerStorageAccountWaitForOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string, d *pluginsdk.ResourceData) (backup.OperationStatus, error) {
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
	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		return resp.(backup.OperationStatus), err
	}
	return resp.(backup.OperationStatus), nil
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
