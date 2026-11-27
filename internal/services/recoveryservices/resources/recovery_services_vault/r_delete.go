// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package recovery_services_vault

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2025-08-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	recoveryservicesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRecoveryServicesVaultDelete(d *pluginsdk.ResourceData, meta any) error {
	client := recoveryservicesClient.From(meta).VaultsClient
	protectedItemsClient := recoveryservicesClient.From(meta).ProtectedItemsGroupClient
	protectedItemClient := recoveryservicesClient.From(meta).ProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := vaults.ParseVaultID(d.Id())
	if err != nil {
		return err
	}

	if features.From(meta).RecoveryService.PurgeProtectedItemsFromVaultOnDestroy {
		log.Printf("[DEBUG] Purging Protected Items from %s", id)

		vaultId := backupprotecteditems.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

		protectedItems, err := protectedItemsClient.ListComplete(ctx, vaultId, backupprotecteditems.ListOperationOptions{})
		if err != nil {
			return fmt.Errorf("listing protected items in %s: %+v", id, err)
		}

		for _, item := range protectedItems.Items {
			if item.Id != nil {
				protectedItemId, err := protecteditems.ParseProtectedItemID(pointer.From(item.Id))
				if err != nil {
					return err
				}

				log.Printf("[DEBUG] Purging %s from %s", protectedItemId, id)

				if err := protectedItemClient.DeleteThenPoll(ctx, *protectedItemId); err != nil {
					return fmt.Errorf("deleting %s: %+v", protectedItemId, err)
				}
			}
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
