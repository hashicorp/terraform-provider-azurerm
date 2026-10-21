// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventhub_namespace_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventHubNamespaceCustomerManagedKeyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := namespaces.ParseNamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if resp.Model.Properties == nil && resp.Model.Properties.Encryption == nil {
		d.SetId("")
		return nil
	}

	d.Set("eventhub_namespace_id", id.ID())

	if props := resp.Model.Properties; props != nil {
		keyVaultKeyIds, err := flattenEventHubNamespaceKeyVaultKeyIds(props.Encryption)
		if err != nil {
			return err
		}

		d.Set("key_vault_key_ids", keyVaultKeyIds)
		d.Set("infrastructure_encryption_enabled", props.Encryption.RequireInfrastructureEncryption)

		if kvprops := props.Encryption.KeyVaultProperties; kvprops != nil {
			// we can only have a single user managed id for N number of keys, azure portal only allows setting a single one and then applies it to each key
			for _, item := range *kvprops {
				if item.Identity != nil && item.Identity.UserAssignedIdentity != nil {
					userAssignedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*item.Identity.UserAssignedIdentity)
					if err != nil {
						return fmt.Errorf("parsing `user_assigned_identity_id`: %+v", err)
					}
					if err := d.Set("user_assigned_identity_id", userAssignedId.ID()); err != nil {
						return fmt.Errorf("setting `user_assigned_identity_id`: %+v", err)
					}

					break
				}
			}
		}
	}

	return nil
}

func flattenEventHubNamespaceKeyVaultKeyIds(input *namespaces.Encryption) ([]string, error) {
	results := make([]string, 0)
	if input == nil || input.KeyVaultProperties == nil {
		return results, nil
	}

	for _, item := range *input.KeyVaultProperties {
		keyVaultKeyId, err := keyvault.NewNestedItemID(pointer.From(item.KeyVaultUri), keyvault.NestedItemTypeKey, pointer.From(item.KeyName), pointer.From(item.KeyVersion))
		if err != nil {
			return nil, err
		}

		results = append(results, keyVaultKeyId.ID())
	}

	return results, nil
}
