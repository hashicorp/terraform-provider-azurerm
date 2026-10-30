// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_vm

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesBackupProtectedVMUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	vaultClient := meta.(*clients.Client).RecoveryServices.VaultsClient
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

	if d.HasChange("backup_policy_id") {
		properties.PolicyId = pointer.To(d.Get("backup_policy_id").(string))
	}

	if d.HasChanges("exclude_disk_luns", "include_disk_luns") {
		properties.ExtendedProperties = expandDiskExclusion(d)
	}

	if d.HasChange("protection_state") {
		protectionState := d.Get("protection_state").(string)
		if protectionState == string(protecteditems.ProtectionStateBackupsSuspended) {
			if err := checkRecoveryServicesVaultIsImmutable(ctx, vaultClient, vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)); err != nil {
				return err
			}
		}
		properties.ProtectionState = pointer.ToEnum[protecteditems.ProtectionState](protectionState)
	}
	model.Properties = properties

	if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceRecoveryServicesBackupProtectedVMRead(d, meta)
}
