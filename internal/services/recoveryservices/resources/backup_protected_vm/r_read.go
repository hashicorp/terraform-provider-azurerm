// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2024-10-01/protectionpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesBackupProtectedVMRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
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
						if err := d.Set("include_disk_luns", pluginsdk.FlattenSlice(v.DiskExclusionProperties.DiskLunList)); err != nil {
							return fmt.Errorf("setting include_disk_luns: %+v", err)
						}
					} else {
						if err := d.Set("exclude_disk_luns", pluginsdk.FlattenSlice(v.DiskExclusionProperties.DiskLunList)); err != nil {
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

func checkRecoveryServicesVaultIsImmutable(ctx context.Context, client *vaults.VaultsClient, vaultId vaults.VaultId) error {
	// While not ideal, if `protection_state` = `BackupsSuspended`, we get the recovery vault so we can ensure it's in an immutable state
	// We're doing this here because the error message provided by Azure if it is not in an immutable state is confusing.
	// Relevant issue: https://github.com/Azure/azure-rest-api-specs/issues/32688

	existingVault, err := client.Get(ctx, vaultId)
	if err != nil {
		if !response.WasNotFound(existingVault.HttpResponse) {
			return fmt.Errorf("checking for presence of Recovery Services Vault %s: %+v", vaultId, err)
		}
	}

	if existingVault.Model != nil &&
		existingVault.Model.Properties != nil &&
		existingVault.Model.Properties.SecuritySettings != nil &&
		existingVault.Model.Properties.SecuritySettings.ImmutabilitySettings != nil {
		immutabilityState := pointer.From(existingVault.Model.Properties.SecuritySettings.ImmutabilitySettings.State)
		if immutabilityState == vaults.ImmutabilityStateDisabled {
			return errors.New("`protection_state` cannot be set to `BackupsSuspended` while the Recovery Services Vault is not in an immutable (`Locked` / `Unlocked`) state")
		}
	}

	return nil
}
