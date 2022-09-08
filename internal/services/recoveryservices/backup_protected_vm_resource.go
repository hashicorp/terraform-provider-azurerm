package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	vmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRecoveryServicesBackupProtectedVM() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesBackupProtectedVMCreateUpdate,
		Read:   resourceRecoveryServicesBackupProtectedVMRead,
		Update: resourceRecoveryServicesBackupProtectedVMCreateUpdate,
		Delete: resourceRecoveryServicesBackupProtectedVMDelete,

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

		Schema: resourceRecoveryServicesBackupProtectedVMSchema(),

		// It's possible to remove the associated vm from the protected backup so we'll only ForceNew this attribute if it's
		// changing to something other than empty.
		CustomizeDiff: pluginsdk.ForceNewIfChange("source_vm_id", func(ctx context.Context, old, new, meta interface{}) bool {
			return new.(string) != "" && old.(string) != new.(string)
		}),
	}
}

func resourceRecoveryServicesBackupProtectedVMCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)

	vaultName := d.Get("recovery_vault_name").(string)

	// source_vm_id must be specified at creation time but can be removed during update
	if d.IsNewResource() {
		if _, ok := d.GetOk("source_vm_id"); !ok {
			return fmt.Errorf("`source_vm_id` must be specified when creating")
		}
	}
	vmId := d.Get("source_vm_id").(string)
	policyId := d.Get("backup_policy_id").(string)

	// get VM name from id
	parsedVmId, err := vmParse.VirtualMachineID(vmId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_vm_id '%s': %+v", vmId, err)
	}

	protectedItemName := fmt.Sprintf("VM;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, parsedVmId.Name)
	containerName := fmt.Sprintf("iaasvmcontainer;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroup, parsedVmId.Name)

	log.Printf("[DEBUG] Creating/updating Azure Backup Protected VM %s (resource group %q)", protectedItemName, resourceGroup)

	if d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_protected_vm", *existing.ID)
		}
	}

	item := backup.ProtectedItemResource{
		Properties: &backup.AzureIaaSComputeVMProtectedItem{
			PolicyID:           &policyId,
			ProtectedItemType:  backup.ProtectedItemTypeMicrosoftClassicComputevirtualMachines,
			WorkloadType:       backup.DataSourceTypeVM,
			SourceResourceID:   utils.String(vmId),
			FriendlyName:       utils.String(parsedVmId.Name),
			ExtendedProperties: expandDiskExclusion(d),
			VirtualMachineID:   utils.String(vmId),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, item); err != nil {
		return fmt.Errorf("creating/updating Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	resp, err := resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyId, d)
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1) // This code is a workaround for this bug https://github.com/Azure/azure-sdk-for-go/issues/2824
	d.SetId(id)

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func resourceRecoveryServicesBackupProtectedVMRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading Azure Backup Protected VM %q (resource group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Get(ctx, id.VaultName, id.ResourceGroup, "Azure", id.ProtectionContainerName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Backup Protected VM %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)

	if properties := resp.Properties; properties != nil {
		if vm, ok := properties.AsAzureIaaSComputeVMProtectedItem(); ok {
			d.Set("source_vm_id", vm.SourceResourceID)

			if v := vm.PolicyID; v != nil {
				d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
			}

			if v := vm.ExtendedProperties; v != nil && v.DiskExclusionProperties != nil {
				if *v.DiskExclusionProperties.IsInclusionList {
					if err := d.Set("include_disk_luns", utils.FlattenInt32Slice(v.DiskExclusionProperties.DiskLunList)); err != nil {
						return fmt.Errorf("setting include_disk_luns: %+v", err)
					}
				} else {
					if err := d.Set("exclude_disk_luns", utils.FlattenInt32Slice(v.DiskExclusionProperties.DiskLunList)); err != nil {
						return fmt.Errorf("setting exclude_disk_luns: %+v", err)
					}
				}
			}
		}
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Azure Backup Protected Item %q (resource group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.VaultName, id.ResourceGroup, "Azure", id.ProtectionContainerName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("issuing delete request for Azure Backup Protected VM %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if _, err := resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx, client, id.VaultName, id.ResourceGroup, id.ProtectionContainerName, id.Name, d); err != nil {
		return err
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyId string, d *pluginsdk.ResourceData) (backup.ProtectedItemResource, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx, client, vaultName, resourceGroup, containerName, protectedItemName, policyId, true),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("waiting for the Azure Backup Protected VM %q to be true (Resource Group %q) to provision: %+v", protectedItemName, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, d *pluginsdk.ResourceData) (backup.ProtectedItemResource, error) {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Pending"},
		Target:     []string{"NotFound", "Stopped"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, "NotFound", nil
				}

				return resp, "Error", fmt.Errorf("making Read request on Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
			}

			if properties := resp.Properties; properties != nil {
				if vm, ok := properties.AsAzureIaaSComputeVMProtectedItem(); ok {
					if strings.EqualFold(string(vm.ProtectionState), string(backup.ProtectionStateProtectionStopped)) {
						return resp, "Stopped", nil
					}
				}
			}
			return resp, "Pending", nil
		},

		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}

	resp, err := state.WaitForStateContext(ctx)
	if err != nil {
		i, _ := resp.(backup.ProtectedItemResource)
		return i, fmt.Errorf("waiting for the Azure Backup Protected VM %q to be deleted (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	return resp.(backup.ProtectedItemResource), nil
}

func resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx context.Context, client *backup.ProtectedItemsClient, vaultName, resourceGroup, containerName, protectedItemName string, policyId string, newResource bool) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, "Azure", containerName, protectedItemName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
		}
		return resp, "Found", nil
	}
}

func expandDiskExclusion(d *pluginsdk.ResourceData) *backup.ExtendedProperties {
	if v, ok := d.GetOk("include_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &backup.ExtendedProperties{
			DiskExclusionProperties: &backup.DiskExclusionProperties{
				DiskLunList:     utils.ExpandInt32Slice(diskLun),
				IsInclusionList: utils.Bool(true),
			},
		}
	}

	if v, ok := d.GetOk("exclude_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &backup.ExtendedProperties{
			DiskExclusionProperties: &backup.DiskExclusionProperties{
				DiskLunList:     utils.ExpandInt32Slice(diskLun),
				IsInclusionList: utils.Bool(false),
			},
		}
	}
	return nil
}

func expandDiskLunList(input []interface{}) []interface{} {
	result := make([]interface{}, 0, len(input))
	for _, v := range input {
		result = append(result, v.(int))
	}
	return result
}

func resourceRecoveryServicesBackupProtectedVMSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": azure.SchemaResourceGroupName(),

		"recovery_vault_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.RecoveryServicesVaultName,
		},

		"source_vm_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validation.StringIsEmpty,
				azure.ValidateResourceID,
			),
			// TODO: make this case sensitive once the API's fixed https://github.com/Azure/azure-rest-api-specs/issues/10357
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"backup_policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"exclude_disk_luns": {
			Type:          pluginsdk.TypeSet,
			ConflictsWith: []string{"include_disk_luns"},
			Optional:      true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},

		"include_disk_luns": {
			Type:          pluginsdk.TypeSet,
			ConflictsWith: []string{"exclude_disk_luns"},
			Optional:      true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeInt,
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}
