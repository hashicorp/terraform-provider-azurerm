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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	vmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
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
			_, err := protecteditems.ParseProtectedItemID(id)
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
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

	id := protecteditems.NewProtectedItemID(subscriptionId, resourceGroup, vaultName, "Azure", containerName, protectedItemName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_backup_protected_vm", id.ID())
		}
	}

	var protectedItem protecteditems.ProtectedItem = protecteditems.AzureIaaSComputeVMProtectedItem{
		PolicyId:           &policyId,
		WorkloadType:       pointer.To(protecteditems.DataSourceTypeVM),
		SourceResourceId:   pointer.To(vmId),
		FriendlyName:       pointer.To(parsedVmId.Name),
		ExtendedProperties: expandDiskExclusion(d),
		VirtualMachineId:   pointer.To(vmId),
	}
	item := protecteditems.ProtectedItemResource{
		Properties: &protectedItem,
	}

	if _, err = client.CreateOrUpdate(ctx, id, item); err != nil {
		return fmt.Errorf("creating/updating Azure Backup Protected VM %q (Resource Group %q): %+v", protectedItemName, resourceGroup, err)
	}

	if err = resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx, client, id, d); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func resourceRecoveryServicesBackupProtectedVMRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Reading %s", id)

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
			if vm, ok := (*properties).(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
				d.Set("source_vm_id", vm.SourceResourceId)

				if v := vm.PolicyId; v != nil {
					// TODO: update to use an ID Parser
					d.Set("backup_policy_id", strings.Replace(*v, "Subscriptions", "subscriptions", 1))
				}

				if v := vm.ExtendedProperties; v != nil && v.DiskExclusionProperties != nil {
					if *v.DiskExclusionProperties.IsInclusionList {
						if err := d.Set("include_disk_luns", utils.FlattenInt64Slice(v.DiskExclusionProperties.DiskLunList)); err != nil {
							return fmt.Errorf("setting include_disk_luns: %+v", err)
						}
					} else {
						if err := d.Set("exclude_disk_luns", utils.FlattenInt64Slice(v.DiskExclusionProperties.DiskLunList)); err != nil {
							return fmt.Errorf("setting exclude_disk_luns: %+v", err)
						}
					}
				}
			}
		}
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	opResultClient := meta.(*clients.Client).RecoveryServices.BackupOperationResultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", id)

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("issuing delete request for %s: %+v", id, err)
		}
	}

	locationURL, err := resp.HttpResponse.Location()
	if err != nil || locationURL == nil {
		return fmt.Errorf("deleting %s: Location header missing or empty", id)
	}

	parsedLocation, err := azure.ParseAzureResourceID(handleAzureSdkForGoBug2824(locationURL.Path))
	if err != nil {
		return err
	}

	if err = resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx, client, opResultClient, *id, parsedLocation.Path["backupOperationResults"], d); err != nil {
		return err
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForStateCreateUpdate(ctx context.Context, client *protecteditems.ProtectedItemsClient, id protecteditems.ProtectedItemId, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx, client, id),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to provision: %+v", id, err)
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMWaitForDeletion(ctx context.Context, client *protecteditems.ProtectedItemsClient, opResultClient *backup.OperationResultsClient, id protecteditems.ProtectedItemId, operationId string, d *pluginsdk.ResourceData) error {
	state := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Pending"},
		Target:     []string{"NotFound", "Stopped"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, "NotFound", nil
				}

				return resp, "Error", fmt.Errorf("making Read request on %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					if vm, ok := (*properties).(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
						if vm.ProtectionState != nil && strings.EqualFold(string(*vm.ProtectionState), string(backup.ProtectionStateProtectionStopped)) {
							return resp, "Stopped", nil
						}
					}
				}
			}
			return resp, "Pending", nil
		},

		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}

	_, err := state.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s: %+v", id, err)
	}

	// we should also wait for the operation to complete, or it will fail when creating a new backup vm with the same vm in different vault immediately.
	opState := &pluginsdk.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"202"},
		Target:     []string{"200", "204"},
		Refresh: func() (interface{}, string, error) {
			resp, err := opResultClient.Get(ctx, id.VaultName, id.ResourceGroupName, operationId)
			if err != nil {
				return nil, "Error", fmt.Errorf("making Read request on Recovery Service Protected Item operation %q for %s: %+v", operationId, id, err)
			}
			return resp, strconv.Itoa(resp.StatusCode), err
		},

		Timeout: d.Timeout(pluginsdk.TimeoutDelete),
	}

	_, err = opState.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the Recovery Service Protected Item operation to be deleted for %s: %+v", id, err)
	}

	return nil
}

func resourceRecoveryServicesBackupProtectedVMRefreshFunc(ctx context.Context, client *protecteditems.ProtectedItemsClient, id protecteditems.ProtectedItemId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("making Read request on %s: %+v", id, err)
		}
		return resp, "Found", nil
	}
}

func expandDiskExclusion(d *pluginsdk.ResourceData) *protecteditems.ExtendedProperties {
	if v, ok := d.GetOk("include_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &protecteditems.ExtendedProperties{
			DiskExclusionProperties: &protecteditems.DiskExclusionProperties{
				DiskLunList:     utils.ExpandInt64Slice(diskLun),
				IsInclusionList: utils.Bool(true),
			},
		}
	}

	if v, ok := d.GetOk("exclude_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &protecteditems.ExtendedProperties{
			DiskExclusionProperties: &protecteditems.DiskExclusionProperties{
				DiskLunList:     utils.ExpandInt64Slice(diskLun),
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
		"resource_group_name": commonschema.ResourceGroupName(),

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
