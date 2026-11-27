// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) flatten(metadata sdk.ResourceMetaData, id *namespaces.NamespaceId, model *namespaces.SBNamespace) error {
	state := ServiceBusNamespaceCustomerManagedKeyModel{
		NamespaceID: id.ID(),
	}

	if model != nil {
		if props := model.Properties; props != nil && props.Encryption != nil {
			encryption := props.Encryption
			if keyVaultProperties := encryption.KeyVaultProperties; keyVaultProperties != nil && len(*keyVaultProperties) > 0 {
				keyVaultKeyId, err := keyvault.NewNestedItemID(pointer.From((*keyVaultProperties)[0].KeyVaultUri), keyvault.NestedItemTypeKey, pointer.From((*keyVaultProperties)[0].KeyName), pointer.From((*keyVaultProperties)[0].KeyVersion))
				if err != nil {
					return fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
				}

				state.KeyVaultKeyID = keyVaultKeyId.ID()
			}

			state.InfrastructureEncryptionEnabled = pointer.From(encryption.RequireInfrastructureEncryption)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}
