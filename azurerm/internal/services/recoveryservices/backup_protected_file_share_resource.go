package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceBackupProtectedFileShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceBackupProtectedFileShareCreateUpdate,
		Read:   resourceBackupProtectedFileShareRead,
		Update: resourceBackupProtectedFileShareCreateUpdate,
		Delete: resourceBackupProtectedFileShareDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(80 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(80 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},

			"source_storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"source_file_share_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareName,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceBackupProtectedFileShareCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	protectableClient := meta.(*clients.Client).RecoveryServices.ProtectableItemsClient
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)

	vaultName := d.Get("recovery_vault_name").(string)
	storageAccountID := d.Get("source_storage_account_id").(string)
	fileShareName := d.Get("source_file_share_name").(string)
	policyID := d.Get("backup_policy_id").(string)

	// get storage account name from id
	parsedStorageAccountID, err := azure.ParseAzureResourceID(storageAccountID)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_storage_account_id '%s': %+v", storageAccountID, err)
	}
	accountName, hasName := parsedStorageAccountID.Path["storageAccounts"]
	if !hasName {
		return fmt.Errorf("[ERROR] parsed source_storage_account_id '%s' doesn't contain 'storageAccounts'", storageAccountID)
	}

	// The fileshare has a user defined name, but its system name (fileShareSystemName) is only known to Azure Backup
	filter := fmt.Sprintf("backupManagementType eq 'AzureStorage' and friendlyName eq '%s'", fileShareName)
	backupProtectableItemsResponse, err := protectableClient.List(ctx, vaultName, resourceGroup, filter, "")
	if err != nil {
		return fmt.Errorf("Error checking for protectable fileshares in Recovery Service Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}
	backupProtectableItems := backupProtectableItemsResponse.Values()

	backupProtectedItemsResponse, err := protectedClient.List(ctx, vaultName, resourceGroup, filter, "")
	if err != nil {
		return fmt.Errorf("Error checking for protected fileshares in Recovery Service Vault %q (Resource Group %q): %+v", vaultName, resourceGroup, err)
	}
	backupProtectedItems := backupProtectedItemsResponse.Values()

	if backupProtectedItems == nil && backupProtectableItems == nil {
		return fmt.Errorf("[ERROR] fileshare '%s' not found in protectable or protected fileshares, make sure Storage Account %q is registered with Recovery Service Vault %q (Resource Group %q)", fileShareName, accountName, vaultName, resourceGroup)
	}
	if len(backupProtectableItems)+len(backupProtectedItems) > 1 {
		return fmt.Errorf("[ERROR] multiple fileshares found after filtering protectable or protected fileshares where only one is expected")
	}

	fileShareSystemName := ""
	if backupProtectableItems != nil && *backupProtectableItems[0].Name != "" {
		fileShareSystemName = *backupProtectableItems[0].Name
	}
	if backupProtectedItems != nil && *backupProtectedItems[0].Name != "" {
		fileShareSystemName = *backupProtectedItems[0].Name
	}

	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, accountName)

	log.Printf("[DEBUG] creating/updating Recovery Service Protected File Share %q (Container Name %q)", fileShareName, containerName)

	if d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, fileShareSystemName, "")
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Recovery Service Protected File Share %q (Resource Group %q): %+v", fileShareName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_protected_file_share", *existing.ID)
		}
	}

	item := backup.ProtectedItemResource{
		Properties: &backup.AzureFileshareProtectedItem{
			PolicyID:          &policyID,
			ProtectedItemType: backup.ProtectedItemTypeAzureFileShareProtectedItem,
			WorkloadType:      backup.DataSourceTypeAzureFileShare,
			SourceResourceID:  utils.String(storageAccountID),
			FriendlyName:      utils.String(fileShareName),
		},
	}

	resp, err := client.CreateOrUpdate(ctx, vaultName, resourceGroup, "Azure", containerName, fileShareSystemName, item)
	if err != nil {
		return fmt.Errorf("Error creating/updating Recovery Service Protected File Share %q (Resource Group %q): %+v", fileShareName, resourceGroup, err)
	}

	locationURL, err := resp.Response.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("Error creating/updating Azure File Share backup item %q (Vault %q): Location header missing or empty", containerName, vaultName)
	}

	opResourceID := azure.HandleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["operationResults"]

	if _, err := resourceBackupProtectedFileShareWaitForOperation(ctx, opClient, vaultName, resourceGroup, operationID, d); err != nil {
		return err
	}

	resp, err = client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, fileShareSystemName, "")

	if err != nil {
		return fmt.Errorf("Error creating/updating Azure File Share backup item %q (Vault %q): %+v", fileShareName, vaultName, err)
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1) // This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
	d.SetId(id)

	return resourceBackupProtectedFileShareRead(d, meta)
}

func resourceBackupProtectedFileShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	fileShareSystemName := id.Path["protectedItems"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Reading Recovery Service Protected File Share %q (resource group %q)", fileShareSystemName, resourceGroup)

	resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, fileShareSystemName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Recovery Service Protected File Share %q (Vault %q Resource Group %q): %+v", fileShareSystemName, vaultName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("recovery_vault_name", vaultName)

	if properties := resp.Properties; properties != nil {
		if item, ok := properties.AsAzureFileshareProtectedItem(); ok {
			sourceResourceID := strings.Replace(*item.SourceResourceID, "Microsoft.storage", "Microsoft.Storage", 1) // The SDK is returning inconsistent capitalization
			d.Set("source_storage_account_id", sourceResourceID)
			d.Set("source_file_share_name", item.FriendlyName)

			if v := item.PolicyID; v != nil {
				d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
			}
		}
	}

	return nil
}

func resourceBackupProtectedFileShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opClient := meta.(*clients.Client).RecoveryServices.BackupOperationStatusesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	fileShareSystemName := id.Path["protectedItems"]
	resourceGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Deleting Recovery Service Protected Item %q (resource group %q)", fileShareSystemName, resourceGroup)

	resp, err := client.Delete(ctx, vaultName, resourceGroup, "Azure", containerName, fileShareSystemName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Recovery Service Protected File Share %q (Resource Group %q): %+v", fileShareSystemName, resourceGroup, err)
		}
	}

	locationURL, err := resp.Response.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("Error deleting Azure File Share backups item %s (Vault %s): Location header missing or empty", containerName, vaultName)
	}

	opResourceID := azure.HandleAzureSdkForGoBug2824(locationURL.Path)

	parsedLocation, err := azure.ParseAzureResourceID(opResourceID)
	if err != nil {
		return err
	}
	operationID := parsedLocation.Path["backupOperationResults"] // This is different for create and delete requests ¯\_(ツ)_/¯

	if _, err := resourceBackupProtectedFileShareWaitForOperation(ctx, opClient, vaultName, resourceGroup, operationID, d); err != nil {
		return err
	}

	return nil
}

// nolint unused - linter mistakenly things this function isn't used?
func resourceBackupProtectedFileShareWaitForOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string, d *schema.ResourceData) (backup.OperationStatus, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 10 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"InProgress"},
		Target:     []string{"Succeeded"},
		Refresh:    resourceBackupProtectedFileShareCheckOperation(ctx, client, vaultName, resourceGroup, operationID),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	log.Printf("[DEBUG] Waiting for backup operation %s (Vault %s) to complete", operationID, vaultName)
	resp, err := state.WaitForState()
	if err != nil {
		return resp.(backup.OperationStatus), err
	}
	return resp.(backup.OperationStatus), nil
}

func resourceBackupProtectedFileShareCheckOperation(ctx context.Context, client *backup.OperationStatusesClient, vaultName, resourceGroup, operationID string) resource.StateRefreshFunc {
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
			err = fmt.Errorf("Azure Backup operation status failed with status %q (Vault %q Resource Group %q Operation ID %q): %+v", resp.Status, vaultName, resourceGroup, operationID, errMsg)
		}

		log.Printf("[DEBUG] Backup operation %s status is %s", operationID, string(resp.Status))
		return resp, string(resp.Status), err
	}
}
