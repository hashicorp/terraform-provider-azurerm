// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRecoveryServicesBackupProtectedVM() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRecoveryServicesBackupProtectedVMCreate,
		Read:   resourceRecoveryServicesBackupProtectedVMRead,
		Update: resourceRecoveryServicesBackupProtectedVMUpdate,
		Delete: resourceRecoveryServicesBackupProtectedVMDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := protecteditems.ParseProtectedItemID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(120 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(120 * time.Minute),
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

func resourceRecoveryServicesBackupProtectedVMCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	// source_vm_id must be specified at creation time but can be removed during update
	if _, ok := d.GetOk("source_vm_id"); !ok {
		return fmt.Errorf("`source_vm_id` must be specified when creating")
	}

	vmId := d.Get("source_vm_id").(string)
	policyId := d.Get("backup_policy_id").(string)

	if policyId == "" {
		return fmt.Errorf("`backup_policy_id` must be specified during creation")
	}

	// get VM name from id
	parsedVmId, err := commonids.ParseVirtualMachineID(vmId)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to parse source_vm_id '%s': %+v", vmId, err)
	}

	protectedItemName := fmt.Sprintf("VM;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroupName, parsedVmId.VirtualMachineName)
	containerName := fmt.Sprintf("iaasvmcontainer;iaasvmcontainerv2;%s;%s", parsedVmId.ResourceGroupName, parsedVmId.VirtualMachineName)

	log.Printf("[DEBUG] Creating Azure Backup Protected VM %s (resource group %q)", protectedItemName, resourceGroup)

	id := protecteditems.NewProtectedItemID(subscriptionId, resourceGroup, vaultName, "Azure", containerName, protectedItemName)

	existing, err := client.Get(ctx, id, protecteditems.GetOperationOptions{})
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		isSoftDeleted := false
		if existing.Model != nil && existing.Model.Properties != nil {
			if prop, ok := existing.Model.Properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
				isSoftDeleted = pointer.From(prop.IsScheduledForDeferredDelete)
			}
		}

		if isSoftDeleted {
			if meta.(*clients.Client).Features.RecoveryServicesVault.RecoverSoftDeletedBackupProtectedVM {
				err = resourceRecoveryServicesVaultBackupProtectedVMRecoverSoftDeleted(ctx, client, id)
				if err != nil {
					return fmt.Errorf("recovering soft deleted %s: %+v", id, err)
				}
			} else {
				return errors.New(optedOutOfRecoveringSoftDeletedBackupProtectedVMFmt(parsedVmId.ID(), vaultName))
			}
		}

		if !isSoftDeleted {
			return tf.ImportAsExistsError("azurerm_backup_protected_vm", id.ID())
		}
	}

	item := protecteditems.ProtectedItemResource{
		Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
			PolicyId:           &policyId,
			WorkloadType:       pointer.To(protecteditems.DataSourceTypeVM),
			SourceResourceId:   pointer.To(vmId),
			FriendlyName:       pointer.To(parsedVmId.VirtualMachineName),
			ExtendedProperties: expandDiskExclusion(d),
			VirtualMachineId:   pointer.To(vmId),
		},
	}

	protectionState, ok := d.GetOk("protection_state")
	protectionStopped := strings.EqualFold(protectionState.(string), string(protecteditems.ProtectionStateProtectionStopped))
	requireUpdateProtectionState := ok && protectionStopped

	if err := client.CreateOrUpdateThenPoll(ctx, id, item); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// the protection state will be updated in the additional update.
	if requireUpdateProtectionState {
		p := protecteditems.ProtectionState(protectionState.(string))
		updateInput := protecteditems.ProtectedItemResource{
			Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
				ProtectionState:  &p,
				SourceResourceId: utils.String(vmId),
			},
		}

		if err := client.CreateOrUpdateThenPoll(ctx, id, updateInput); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

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

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			if vm, ok := properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
				if vm.IsScheduledForDeferredDelete != nil && *vm.IsScheduledForDeferredDelete {
					d.SetId("")
					return nil
				}

				d.Set("source_vm_id", vm.SourceResourceId)
				d.Set("protection_state", pointer.From(vm.ProtectionState))

				backupPolicyId := ""
				if policyId := pointer.From(vm.PolicyId); policyId != "" {
					parsedPolicyId, err := protectionpolicies.ParseBackupPolicyIDInsensitively(policyId)
					if err != nil {
						return fmt.Errorf("parsing policy ID %q: %+v", policyId, err)
					}
					backupPolicyId = parsedPolicyId.ID()
				}
				d.Set("backup_policy_id", backupPolicyId)

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

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)

	return nil
}

func resourceRecoveryServicesBackupProtectedVMUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		return err
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if _, ok := existing.Model.Properties.(protecteditems.AzureIaaSComputeVMProtectedItem); !ok {
		return fmt.Errorf("retrieving %s: `properties` was not a VM Protected Item", id)
	}

	model := *existing.Model
	properties := existing.Model.Properties.(protecteditems.AzureIaaSComputeVMProtectedItem)
	updateProtectedBackup := false

	if d.HasChange("backup_policy_id") {
		properties.PolicyId = pointer.To(d.Get("backup_policy_id").(string))
		updateProtectedBackup = true
	}

	if d.HasChange("exclude_disk_luns") || d.HasChange("include_disk_luns") {
		properties.ExtendedProperties = expandDiskExclusion(d)
		updateProtectedBackup = true
	}

	model.Properties = properties

	if updateProtectedBackup {
		if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	protectionState := string(pointer.From(properties.ProtectionState))
	protectionStopped := false
	if d.HasChange("protection_state") {
		protectionState = d.Get("protection_state").(string)
		protectionStopped = strings.EqualFold(protectionState, string(protecteditems.ProtectionStateProtectionStopped))
	}

	// the protection state will be updated in the additional update.
	if protectionStopped {
		p := protecteditems.ProtectionState(protectionState)
		vmId := d.Get("source_vm_id").(string)
		updateInput := protecteditems.ProtectedItemResource{
			Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
				ProtectionState:  &p,
				SourceResourceId: utils.String(vmId),
			},
		}

		if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func resourceRecoveryServicesBackupProtectedVMDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	features := meta.(*clients.Client).Features.RecoveryService

	if features.VMBackupStopProtectionAndRetainDataOnDestroy || features.VMBackupSuspendProtectionAndRetainDataOnDestroy {
		log.Printf("[DEBUG] Retaining Data and Stopping Protection for %s", id)

		existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(existing.HttpResponse) {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("making Read request on %s: %+v", id, err)
		}

		desiredState := protecteditems.ProtectionStateProtectionStopped
		if features.VMBackupSuspendProtectionAndRetainDataOnDestroy {
			desiredState = protecteditems.ProtectionStateBackupsSuspended
		}

		if model := existing.Model; model != nil {
			if properties := model.Properties; properties != nil {
				if vm, ok := properties.(protecteditems.AzureIaaSComputeVMProtectedItem); ok {
					updateInput := protecteditems.ProtectedItemResource{
						Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
							ProtectionState:  pointer.To(desiredState),
							SourceResourceId: vm.SourceResourceId,
						},
					}

					if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
						return fmt.Errorf("setting protection to %s and retaining data for %s: %+v", desiredState, id, err)
					}

					return nil
				}
			}
		}
	}

	log.Printf("[DEBUG] Deleting %s", id)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
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

func resourceRecoveryServicesVaultBackupProtectedVMRecoverSoftDeleted(ctx context.Context, client *protecteditems.ProtectedItemsClient, id protecteditems.ProtectedItemId) error {
	payload := protecteditems.ProtectedItemResource{
		Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
			IsRehydrate: pointer.To(true),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("recovering soft-deleted %s: %+v", id, err)
	}

	return nil
}

func optedOutOfRecoveringSoftDeletedBackupProtectedVMFmt(vmId string, vaultName string) string {
	return fmt.Sprintf(`
An existing soft-deleted Backup Protected VM exists with the source VM %q in the recovery services
vault %q, however automatically recovering this Backup Protected VM has been disabled via the 
"features" block.

Terraform can automatically recover the soft-deleted Backup Protected VM when this behaviour is
enabled within the "features" block (located within the "provider" block) - more
information can be found here:

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block

Alternatively you can manually recover this (e.g. using the Azure CLI) and then import
this into Terraform via "terraform import".
`, vmId, vaultName)
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
			Optional:     true,
			ValidateFunc: protectionpolicies.ValidateBackupPolicyID,
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

		"protection_state": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(protecteditems.ProtectedItemStateIRPending),
				string(protecteditems.ProtectedItemStateProtected),
				string(protecteditems.ProtectedItemStateProtectionError),
				string(protecteditems.ProtectedItemStateProtectionStopped),
				string(protecteditems.ProtectedItemStateProtectionPaused),
				string(protecteditems.ProtectionStateInvalid),
			}, false),
		},
	}
}
