package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBackupProtectionContainerStorageAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBackupProtectionContainerStorageAccountCreate,
		Read:   resourceArmBackupProtectionContainerStorageAccountRead,
		Update: nil,
		Delete: resourceArmBackupProtectionContainerStorageAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},
			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmBackupProtectionContainerStorageAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	opStatusClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	storageAccountID := d.Get("storage_account_id").(string)

	parsedStorageAccountID, err := azure.ParseAzureResourceID(storageAccountID)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse storage_account_id '%s': %+v", storageAccountID, err)
	}
	accountName, hasName := parsedStorageAccountID.Path["storageAccounts"]
	if !hasName {
		return fmt.Errorf("[ERROR] parsed storage_account_id '%s' doesn't contain 'storageAccounts'", storageAccountID)
	}

	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, accountName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, vaultName, resGroup, "Azure", containerName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing recovery services protection container %s (Vault %s): %+v", containerName, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_protection_container_storage", azure.HandleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	parameters := backup.ProtectionContainerResource{
		Properties: &backup.AzureStorageContainer{
			SourceResourceID:     &storageAccountID,
			FriendlyName:         &accountName,
			BackupManagementType: backup.ManagementTypeAzureStorage,
			ContainerType:        backup.ContainerTypeStorageContainer1,
		},
	}

	resp, err := client.Register(ctx, vaultName, resGroup, "Azure", containerName, parameters)
	if err != nil {
		return fmt.Errorf("Error registering backup protection container %s (Vault %s): %+v", containerName, vaultName, err)
	}

	locationURL, err := resp.Response.Location() // Operation ID found in the Location header
	if locationURL == nil || err != nil {
		return fmt.Errorf("Unable to determine operation URL for protection container registration status for %s. (Vault %s): Location header missing or empty", containerName, vaultName)
	}

	opResourceID := azure.HandleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}

	operationID := parsedLocation.Path["operationResults"]
	if _, err = resourceArmBackupProtectionContainerStorageAccountWaitForOperation(ctx, opStatusClient, vaultName, resGroup, operationID, d); err != nil {
		return err
	}

	resp, err = client.Get(ctx, vaultName, resGroup, "Azure", containerName)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery protection container %s (Vault %s): %+v", containerName, vaultName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceArmBackupProtectionContainerStorageAccountRead(d, meta)
}

func resourceArmBackupProtectionContainerStorageAccountRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["backupFabrics"]
	containerName := id.Path["protectionContainers"]

	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, vaultName, resGroup, fabricName, containerName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on backup protection container %s (Vault %s): %+v", containerName, vaultName, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)

	if properties, ok := resp.Properties.AsAzureStorageContainer(); ok && properties != nil {
		d.Set("storage_account_id", properties.SourceResourceID)
	}

	return nil
}

func resourceArmBackupProtectionContainerStorageAccountDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["backupFabrics"]
	containerName := id.Path["protectionContainers"]

	client := meta.(*clients.Client).RecoveryServices.BackupProtectionContainersClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Unregister(ctx, vaultName, resGroup, fabricName, containerName)
	if err != nil {
		return fmt.Errorf("Error deregistering backup protection container %s (Vault %s): %+v", containerName, vaultName, err)
	}

	locationURL, err := resp.Response.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("Error unregistering backup protection container %s (Vault %s): Location header missing or empty", containerName, vaultName)
	}

	opResourceID := azure.HandleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["backupOperationResults"]

	if _, err = resourceArmBackupProtectionContainerStorageAccountWaitForOperation(ctx, opClient, vaultName, resGroup, operationID, d); err != nil {
		return err
	}

	return nil
}

// nolint unused - linter mistakenly things this function isn't used?
func resourceArmBackupProtectionContainerStorageAccountWaitForOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string, d *schema.ResourceData) (backup.OperationStatus, error) {
	state := &resource.StateChangeConf{
		MinTimeout:                10 * time.Second,
		Delay:                     10 * time.Second,
		Pending:                   []string{"InProgress"},
		Target:                    []string{"Succeeded"},
		Refresh:                   resourceArmBackupProtectionContainerStorageAccountCheckOperation(ctx, client, vaultName, resourceGroup, operationID),
		ContinuousTargetOccurence: 5, // Without this buffer, file share backups and storage account deletions may fail if performed immediately after creating/destroying the container
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	log.Printf("[DEBUG] Waiting for backup container operation %q (Vault %q) to complete", operationID, vaultName)
	resp, err := state.WaitForState()
	if err != nil {
		return resp.(backup.OperationStatus), err
	}
	return resp.(backup.OperationStatus), nil
}

func resourceArmBackupProtectionContainerStorageAccountCheckOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, operationID)
		if err != nil {
			return resp, "Error", fmt.Errorf("Error making Read request on Recovery Service Protection Container operation %q (Vault %q in Resource Group %q): %+v", operationID, vaultName, resourceGroup, err)
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
