// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_vm

import (
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesBackupProtectedVMCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	vaultClient := meta.(*clients.Client).RecoveryServices.VaultsClient
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
				if err = resourceRecoveryServicesVaultBackupProtectedVMRecoverSoftDeleted(ctx, client, id); err != nil {
					return fmt.Errorf("recovering soft deleted %s: %+v", id, err)
				}
			} else {
				return errors.New(optedOutOfRecoveringSoftDeletedBackupProtectedVMFmt(parsedVmId.ID(), vaultName))
			}
		}

		if !isSoftDeleted {
			if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				return tf.ImportAsExistsError("azurerm_backup_protected_vm", id.ID())
			}
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

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, item, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// the protection state cannot be set during initial creation.
	protectionState := d.Get("protection_state").(string)
	protectionStateUpdateRequired := slices.Contains([]string{
		string(protecteditems.ProtectionStateProtectionStopped),
		string(protecteditems.ProtectionStateBackupsSuspended),
	}, protectionState)

	if protectionStateUpdateRequired {
		if protectionState == string(protecteditems.ProtectionStateBackupsSuspended) {
			if err := checkRecoveryServicesVaultIsImmutable(ctx, vaultClient, vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)); err != nil {
				return err
			}
		}

		updateInput := protecteditems.ProtectedItemResource{
			Properties: &protecteditems.AzureIaaSComputeVMProtectedItem{
				ProtectionState:  pointer.ToEnum[protecteditems.ProtectionState](protectionState),
				SourceResourceId: pointer.To(vmId),
			},
		}

		if err := client.CreateOrUpdateThenPoll(ctx, id, updateInput); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}

func expandDiskExclusion(d *pluginsdk.ResourceData) *protecteditems.ExtendedProperties {
	if v, ok := d.GetOk("include_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &protecteditems.ExtendedProperties{
			DiskExclusionProperties: &protecteditems.DiskExclusionProperties{
				DiskLunList:     pluginsdk.ExpandInt64Slice(diskLun),
				IsInclusionList: pointer.To(true),
			},
		}
	}

	if v, ok := d.GetOk("exclude_disk_luns"); ok {
		diskLun := expandDiskLunList(v.(*pluginsdk.Set).List())

		return &protecteditems.ExtendedProperties{
			DiskExclusionProperties: &protecteditems.DiskExclusionProperties{
				DiskLunList:     pluginsdk.ExpandInt64Slice(diskLun),
				IsInclusionList: pointer.To(false),
			},
		}
	}
	return nil
}

func expandDiskLunList(input []any) []any {
	result := make([]any, 0, len(input))
	for _, v := range input {
		result = append(result, v.(int))
	}
	return result
}
