package azurerm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBackupProtectedFileShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBackupProtectedFileShareCreateUpdate,
		Read:   resourceArmBackupProtectedFileShareRead,
		Update: resourceArmBackupProtectedFileShareCreateUpdate,
		Delete: resourceArmBackupProtectedFileShareDelete,

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
				ValidateFunc: validateArmStorageShareName,
			},

			"backup_policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmBackupProtectedFileShareCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	vaultName := d.Get("recovery_vault_name").(string)
	storageAccountID := d.Get("source_storage_account_id").(string)
	fileShareName := d.Get("source_file_share_name").(string)
	policyID := d.Get("backup_policy_id").(string)

	//get storage account name from id
	parsedStorageAccountID, err := azure.ParseAzureResourceID(storageAccountID)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_storage_account_id '%s': %+v", storageAccountID, err)
	}
	accountName, hasName := parsedStorageAccountID.Path["storageAccounts"]
	if !hasName {
		return fmt.Errorf("[ERROR] parsed source_storage_account_id '%s' doesn't contain 'storageAccounts'", storageAccountID)
	}

	protectedItemName := fmt.Sprintf("AzureFileShare;%s", fileShareName)
	containerName := fmt.Sprintf("StorageContainer;storage;%s;%s", parsedStorageAccountID.ResourceGroup, accountName)

	log.Printf("[DEBUG] Creating/updating Recovery Service Protected File Share %s (Container Name %q)", protectedItemName, containerName)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Recovery Service Protected File Share %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_protected_file_share", *existing.ID)
		}
	}

	item := backup.ProtectedItemResource{
		Tags: tags.Expand(t),
		Properties: &backup.AzureFileshareProtectedItem{
			PolicyID:          &policyID,
			ProtectedItemType: backup.ProtectedItemTypeAzureFileShareProtectedItem,
			WorkloadType:      backup.DataSourceTypeAzureFileShare,
			SourceResourceID:  utils.String(storageAccountID),
			FriendlyName:      utils.String(fileShareName),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, item); err != nil {
		return fmt.Errorf("Error creating/updating Recovery Service Protected File Share %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	resp, err := resourceArmBackupProtectedFileShareWaitForStateCreateUpdate(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyID, d)
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1) // This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
	d.SetId(id)

	return resourceArmBackupProtectedFileShareRead(d, meta)
}

func resourceArmBackupProtectedFileShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	protectedItemName := id.Path["protectedItems"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Reading Recovery Service Protected File Share %q (resource group %q)", protectedItemName, resourceGroup)

	resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Recovery Service Protected File Share %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("recovery_vault_name", vaultName)

	if properties := resp.Properties; properties != nil {
		if item, ok := properties.AsAzureFileshareProtectedItem(); ok {
			d.Set("source_storage_account_id", item.SourceResourceID)
			d.Set("source_file_share_name", item.FriendlyName)

			if v := item.PolicyID; v != nil {
				d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmBackupProtectedFileShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	protectedItemName := id.Path["protectedItems"]
	resourceGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Deleting Recovery Service Protected Item %q (resource group %q)", protectedItemName, resourceGroup)

	resp, err := client.Delete(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Recovery Service Protected File Share %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
		}
	}

	if _, err := resourceArmBackupProtectedFileShareWaitForDeletion(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, "", d); err != nil {
		return err
	}

	return nil
}

func resourceArmBackupProtectedFileShareWaitForStateCreateUpdate(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyID string, d *schema.ResourceData) (backup.ProtectedItemResource, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceArmBackupProtectedFileShareRefreshFunc(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyID, true),
	}

	if features.SupportsCustomTimeouts() {
		if d.IsNewResource() {
			state.Timeout = d.Timeout(schema.TimeoutCreate)
		} else {
			state.Timeout = d.Timeout(schema.TimeoutUpdate)
		}
	} else {
		state.Timeout = 30 * time.Minute
	}

	resp, err := state.WaitForState()
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("Error waiting for the Recovery Service Protected File Share %q to be true (Resource Group %q) to provision: %+v", protectedItemName, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}

func resourceArmBackupProtectedFileShareWaitForDeletion(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyID string, d *schema.ResourceData) (backup.ProtectedItemResource, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceArmBackupProtectedFileShareRefreshFunc(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyID, false),
	}

	if features.SupportsCustomTimeouts() {
		state.Timeout = d.Timeout(schema.TimeoutDelete)
	} else {
		state.Timeout = 30 * time.Minute
	}

	resp, err := state.WaitForState()
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("Error waiting for the Recovery Service Protected File Share %q to be false (Resource Group %q) to provision: %+v", protectedItemName, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}

func resourceArmBackupProtectedFileShareRefreshFunc(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyID string, newResource bool) resource.StateRefreshFunc {
	// TODO: split this into two functions
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("Error making Read request on Recovery Service Protected File Share %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
		} else if !newResource && policyID != "" {
			if properties := resp.Properties; properties != nil {
				if vm, ok := properties.AsAzureFileshareProtectedItem(); ok {
					if v := vm.PolicyID; v != nil {
						if strings.Replace(*v, "Subscriptions", "subscriptions", 1) != policyID {
							return resp, "NotFound", nil
						}
					} else {
						return resp, "Error", fmt.Errorf("Error reading policy ID attribute nil on Recovery Service Protected File Share %q (Resource Group %q)", protectedItemName, resourceGroup)
					}
				} else {
					return resp, "Error", fmt.Errorf("Error reading properties on Recovery Service Protected File Share %q (Resource Group %q)", protectedItemName, resourceGroup)
				}
			} else {
				return resp, "Error", fmt.Errorf("Error reading properties on empty Recovery Service Protected File Share %q (Resource Group %q)", protectedItemName, resourceGroup)
			}
		}
		return resp, "Found", nil
	}
}
