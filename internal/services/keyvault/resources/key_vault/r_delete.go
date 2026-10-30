// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault

import (
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2026-02-01/deletedvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	deletedVaultsClient := meta.(*clients.Client).KeyVault.DeletedVaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKeyVaultID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VaultName, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultName, keyVaultResourceName)

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] Key Vault %q was not found - removing from state", id.VaultName)
			return nil
		}
		return fmt.Errorf("checking if key vault %q exists: %v", *id, err)
	}

	location := ""
	purgeProtectionEnabled := false
	softDeleteEnabled := false
	virtualNetworkNames := make([]string, 0)
	if model := read.Model; model != nil {
		if model.Location != nil {
			location = *model.Location
		}

		// Check to see if purge protection is enabled or not...
		if model.Properties.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *model.Properties.EnablePurgeProtection
		}
		if model.Properties.EnableSoftDelete != nil {
			softDeleteEnabled = *model.Properties.EnableSoftDelete
		}

		// ensure we lock on the latest network names, to ensure we handle Azure's networking layer being limited to one change at a time
		if acls := model.Properties.NetworkAcls; acls != nil {
			if rules := acls.VirtualNetworkRules; rules != nil {
				for _, v := range *rules {
					subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
					if err != nil {
						return err
					}

					if !slices.Contains(virtualNetworkNames, subnetId.VirtualNetworkName) {
						virtualNetworkNames = append(virtualNetworkNames, subnetId.VirtualNetworkName)
					}
				}
			}
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// Purge the soft deleted key vault permanently if the feature flag is enabled
	if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeleteOnDestroy && softDeleteEnabled {
		deletedVaultId := deletedvaults.NewDeletedVaultID(id.SubscriptionId, location, id.VaultName)

		// KeyVaults with Purge Protection Enabled cannot be deleted unless done by Azure
		if purgeProtectionEnabled {
			deletedInfo, err := getSoftDeletedStateForKeyVault(ctx, deletedVaultsClient, deletedVaultId)
			if err != nil {
				return fmt.Errorf("retrieving the Deletion Details for %s: %+v", *id, err)
			}

			// in the future it'd be nice to raise a warning, but this is the best we can do for now
			if deletedInfo != nil {
				log.Printf("[DEBUG] The Key Vault %q has Purge Protection Enabled and was deleted on %q. Azure will purge this on %q", id.VaultName, deletedInfo.deleteDate, deletedInfo.purgeDate)
			} else {
				log.Printf("[DEBUG] The Key Vault %q has Purge Protection Enabled and will be purged automatically by Azure", id.VaultName)
			}
			return nil
		}

		log.Printf("[DEBUG] KeyVault %q marked for purge - executing purge", id.VaultName)
		if err := deletedVaultsClient.VaultsPurgeDeletedThenPoll(ctx, deletedVaultId); err != nil {
			return fmt.Errorf("purging %s: %+v", *id, err)
		}
		log.Printf("[DEBUG] Purged KeyVault %q.", id.VaultName)
	}

	meta.(*clients.Client).KeyVault.Purge(*id)

	return nil
}
