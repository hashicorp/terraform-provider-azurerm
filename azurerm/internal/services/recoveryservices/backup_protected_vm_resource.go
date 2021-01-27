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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceRecoveryServicesBackupProtectedVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecoveryServicesBackupProtectedVMCreateUpdate,
		Read:   resourceRecoveryServicesBackupProtectedVMRead,
		Update: resourceRecoveryServicesBackupProtectedVMCreateUpdate,
		Delete: resourceRecoveryServicesBackupProtectedVMDelete,

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

			"source_vm_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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

func resourceRecoveryServicesBackupProtectedVMCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	vaultName := d.Get("recovery_vault_name").(string)
	vmId := d.Get("source_vm_id").(string)
	policyId := d.Get("backup_policy_id").(string)

	// get VM name from id
	parsedVmId, err := azure.ParseAzureResourceID(vmId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_vm_id '%s': %+v", vmId, err)
	}
	vmName, hasName := parsedVmId.Path["virtualMachines"]
	if !hasName {
		return fmt.Errorf("[ERROR] parsed source_vm_id '%s' doesn't contain 'virtualMachines'", vmId)
	}

	protectedItemName := fmt.Sprintf("VM;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, vmName)
	containerName := fmt.Sprintf("iaasvmcontainer;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, vmName)

	log.Printf("[DEBUG] Creating/updating Azure Backup Protected VM %s (resource group %q)", protectedItemName, resourceGroup)

	if d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_protected_vm", *existing.ID)
		}
	}

	item := backup.ProtectedItemResource{
		Tags: tags.Expand(t),
		Properties: &backup.AzureIaaSComputeVMProtectedItem{
			PolicyID:          &policyId,
			ProtectedItemType: backup.ProtectedItemTypeMicrosoftClassicComputevirtualMachines,
			WorkloadType:      backup.DataSourceTypeVM,
			SourceResourceID:  utils.String(vmId),
			FriendlyName:      utils.String(vmName),
			VirtualMachineID:  utils.String(vmId),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, item); err != nil {
		return fmt.Errorf("Error creating/updating Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	resp, err := resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyId, d)
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1) // This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
	d.SetId(id)

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func resourceRecoveryServicesBackupProtectedVMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	protectedItemName := id.Path["protectedItems"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Reading Azure Backup Protected VM %q (resource group %q)", protectedItemName, resourceGroup)

	resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("recovery_vault_name", vaultName)

	if properties := resp.Properties; properties != nil {
		if vm, ok := properties.AsAzureIaaSComputeVMProtectedItem(); ok {
			d.Set("source_vm_id", vm.SourceResourceID)

			if v := vm.PolicyID; v != nil {
				d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceRecoveryServicesBackupProtectedVMDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	protectedItemName := id.Path["protectedItems"]
	resourceGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Deleting Azure Backup Protected Item %q (resource group %q)", protectedItemName, resourceGroup)

	resp, err := client.Delete(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
		}
	}

	if _, err := resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, "", d); err != nil {
		return err
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyId string, d *schema.ResourceData) (backup.ProtectedItemResource, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyId, true),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	resp, err := state.WaitForState()
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("Error waiting for the Azure Backup Protected VM %q to be true (Resource Group %q) to provision: %+v", protectedItemName, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyId string, d *schema.ResourceData) (backup.ProtectedItemResource, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyId, false),
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	resp, err := state.WaitForState()
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("Error waiting for the Azure Backup Protected VM %q to be false (Resource Group %q) to provision: %+v", protectedItemName, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}

func resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyId string, newResource bool) resource.StateRefreshFunc {
	// TODO: split this into two functions
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("Error making Read request on Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
		} else if !newResource && policyId != "" {
			if properties := resp.Properties; properties != nil {
				if vm, ok := properties.AsAzureIaaSComputeVMProtectedItem(); ok {
					if v := vm.PolicyID; v != nil {
						if strings.Replace(*v, "Subscriptions", "subscriptions", 1) != policyId {
							return resp, "NotFound", nil
						}
					} else {
						return resp, "Error", fmt.Errorf("Error reading policy ID attribute nil on Azure Backup Protected VM %q (Resource Group %q)", protectedItemName, resourceGroup)
					}
				} else {
					return resp, "Error", fmt.Errorf("Error reading properties on Azure Backup Protected VM %q (Resource Group %q)", protectedItemName, resourceGroup)
				}
			} else {
				return resp, "Error", fmt.Errorf("Error reading properties on empty Azure Backup Protected VM %q (Resource Group %q)", protectedItemName, resourceGroup)
			}
		}
		return resp, "Found", nil
	}
}
