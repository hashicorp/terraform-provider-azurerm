// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_vm

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	providerFeatures "github.com/hashicorp/terraform-provider-azurerm/internal/features"
	recoveryservicesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesBackupProtectedVMDelete(d *pluginsdk.ResourceData, meta any) error {
	client := recoveryservicesClient.From(meta).ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	features := providerFeatures.From(meta).RecoveryService

	if features.VMBackupStopProtectionAndRetainDataOnDestroy || features.VMBackupSuspendProtectionAndRetainDataOnDestroy {
		log.Printf("[DEBUG] Retaining Data and Stopping Protection for %s", id)

		existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(existing.HttpResponse) {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("retrieving %s: %+v", *id, err)
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
						return fmt.Errorf("setting protection to %s and retaining data for %s: %+v", desiredState, *id, err)
					}

					return nil
				}
			}
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
