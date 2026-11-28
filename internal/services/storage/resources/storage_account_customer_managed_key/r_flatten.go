// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_customer_managed_key

import (
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceStorageAccountCustomerManagedKeyFlatten(d *pluginsdk.ResourceData, id *commonids.StorageAccountId, storageAccount *storageaccounts.StorageAccount) error {
	d.Set("storage_account_id", id.ID())

	enabled := false
	if storageAccount != nil {
		if props := storageAccount.Properties; props != nil {
			if encryption := props.Encryption; encryption != nil && pointer.From(encryption.KeySource) == storageaccounts.KeySourceMicrosoftPointKeyvault {
				enabled = true

				if kvProps := encryption.Keyvaultproperties; kvProps != nil {
					keyID, err := keyvault.NewNestedItemID(pointer.From(kvProps.Keyvaulturi), keyvault.NestedItemTypeKey, pointer.From(kvProps.Keyname), pointer.From(kvProps.Keyversion))
					if err != nil {
						return err
					}
					d.Set("key_vault_key_id", keyID.ID())
				}

				if identityProps := encryption.Identity; identityProps != nil {
					d.Set("user_assigned_identity_id", identityProps.UserAssignedIdentity)
					d.Set("federated_identity_client_id", identityProps.FederatedIdentityClientId)
				}
			}
		}
	}

	if !enabled {
		log.Printf("[DEBUG] Customer Managed Key was not defined for %s - removing from state!", id)
		d.SetId("")
		return nil
	}

	return pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual)
}
