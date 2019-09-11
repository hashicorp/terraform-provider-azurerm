package azurerm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2017-07-01/backup"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesProtectedVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesProtectedVmCreateUpdate,
		Read:   resourceArmRecoveryServicesProtectedVmRead,
		Update: resourceArmRecoveryServicesProtectedVmCreateUpdate,
		Delete: resourceArmRecoveryServicesProtectedVmDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceArmRecoveryServicesProtectedVmCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServices.ProtectedItemsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	vaultName := d.Get("recovery_vault_name").(string)
	vmId := d.Get("source_vm_id").(string)
	policyId := d.Get("backup_policy_id").(string)

	//get VM name from id
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

	log.Printf("[DEBUG] Creating/updating Recovery Service Protected VM %s (resource group %q)", protectedItemName, resourceGroup)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Recovery Service Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_recovery_services_protected_vm", *existing.ID)
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
		return fmt.Errorf("Error creating/updating Recovery Service Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	resp, err := resourceArmRecoveryServicesProtectedVmWaitForState(client, ctx, true, vaultName, resourceGroup, containerName, protectedItemName, policyId, d.IsNewResource())
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1) // This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
	d.SetId(id)

	return resourceArmRecoveryServicesProtectedVmRead(d, meta)
}

func resourceArmRecoveryServicesProtectedVmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServices.ProtectedItemsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	protectedItemName := id.Path["protectedItems"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup
	containerName := id.Path["protectionContainers"]

	log.Printf("[DEBUG] Reading Recovery Service Protected VM %q (resource group %q)", protectedItemName, resourceGroup)

	resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Recovery Service Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
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

func resourceArmRecoveryServicesProtectedVmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServices.ProtectedItemsClient
	ctx := meta.(*ArmClient).StopContext

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
			return fmt.Errorf("Error issuing delete request for Recovery Service Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
		}
	}

	if _, err := resourceArmRecoveryServicesProtectedVmWaitForState(client, ctx, false, vaultName, resourceGroup, containerName, protectedItemName, "", false); err != nil {
		return err
	}

	return nil
}

func resourceArmRecoveryServicesProtectedVmWaitForState(client *backup.ProtectedItemsGroupClient, ctx context.Context, found bool, vaultName, resourceGroup, containerName, protectedItemName string, policyId string, newResource bool) (backup.ProtectedItemResource, error) {
	state := &resource.StateChangeConf{
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Refresh: func() (interface{}, string, error) {

			resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, "NotFound", nil
				}

				return resp, "Error", fmt.Errorf("Error making Read request on Recovery Service Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
			} else if !newResource && policyId != "" {
				if properties := resp.Properties; properties != nil {
					if vm, ok := properties.AsAzureIaaSComputeVMProtectedItem(); ok {
						if v := vm.PolicyID; v != nil {
							if strings.Replace(*v, "Subscriptions", "subscriptions", 1) != policyId {
								return resp, "NotFound", nil
							}
						} else {
							return resp, "Error", fmt.Errorf("Error reading policy ID attribute nil on Recovery Service Protected VM %q (Resource Group %q)", protectedItemName, resourceGroup)
						}
					} else {
						return resp, "Error", fmt.Errorf("Error reading properties on Recovery Service Protected VM %q (Resource Group %q)", protectedItemName, resourceGroup)
					}
				} else {
					return resp, "Error", fmt.Errorf("Error reading properties on empty Recovery Service Protected VM %q (Resource Group %q)", protectedItemName, resourceGroup)
				}
			}
			return resp, "Found", nil
		},
	}

	if found {
		state.Pending = []string{"NotFound"}
		state.Target = []string{"Found"}
	} else {
		state.Pending = []string{"Found"}
		state.Target = []string{"NotFound"}
	}

	resp, err := state.WaitForState()
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("Error waiting for the Recovery Service Protected VM %q to be %t (Resource Group %q) to provision: %+v", protectedItemName, found, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}
