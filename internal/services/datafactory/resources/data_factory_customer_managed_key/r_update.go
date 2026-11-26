// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_customer_managed_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Factories

			var customerManagedKey DataFactoryCustomerManagedKeyModel
			if err := metadata.Decode(&customerManagedKey); err != nil {
				return err
			}

			id, err := factories.ParseFactoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dataFactory, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(dataFactory.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if dataFactory.Model == nil || dataFactory.Model.Properties == nil || dataFactory.Model.Properties.Encryption == nil {
				return fmt.Errorf("retrieving encryption for %s", id)
			}

			payload := dataFactory.Model
			encryption := payload.Properties.Encryption

			if metadata.ResourceData.HasChange("customer_managed_key_id") {
				keyVaultKey, err := keyvault.ParseNestedItemID(customerManagedKey.CustomerManagedKeyID, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
				if err != nil {
					return fmt.Errorf("could not parse Key Vault Key ID: %+v", err)
				}
				encryption.VaultBaseURL = keyVaultKey.KeyVaultBaseURL
				encryption.KeyName = keyVaultKey.Name
				encryption.KeyVersion = pointer.To(keyVaultKey.Version)
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_id") {
				encryption.Identity = helpers.ExpandDataFactoryEncryptionIdentity(customerManagedKey.UserAssignedIdentityID)
			}

			payload.Properties.Encryption = encryption
			if _, err = client.CreateOrUpdate(ctx, *id, *payload, factories.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
