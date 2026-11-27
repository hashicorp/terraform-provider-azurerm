// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package backup_protected_file_share

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBackupProtectedFileShareDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := protecteditems.ParseProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	features := meta.(*clients.Client).Features.RecoveryService

	if features.FileShareBackupStopProtectionAndRetainDataOnDestroy {
		existing, err := client.Get(ctx, *id, protecteditems.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(existing.HttpResponse) {
				return nil
			}
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if existing.Model == nil {
			return fmt.Errorf("retrieving %s: model was nil", *id)
		}

		item, ok := existing.Model.Properties.(protecteditems.AzureFileshareProtectedItem)
		if !ok {
			return fmt.Errorf("retrieving %s: expected AzureFileshareProtectedItem, got %T", *id, existing.Model.Properties)
		}

		item.ProtectionState = pointer.To(protecteditems.ProtectionStateProtectionStopped)
		update := protecteditems.ProtectedItemResource{
			Properties: item,
		}

		if err := client.CreateOrUpdateThenPoll(ctx, *id, update); err != nil {
			return fmt.Errorf("stopping protection and retaining data for %s: %+v", *id, err)
		}

		return nil
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
