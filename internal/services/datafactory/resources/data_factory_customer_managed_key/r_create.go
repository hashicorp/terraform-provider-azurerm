// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Factories

			var customerManagedKey DataFactoryCustomerManagedKeyModel
			if err := metadata.Decode(&customerManagedKey); err != nil {
				return err
			}

			id, err := factories.ParseFactoryID(customerManagedKey.DataFactoryID)
			if err != nil {
				return err
			}

			dataFactory, err := client.Get(ctx, *id, factories.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(dataFactory.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if dataFactory.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if dataFactory.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
			}

			payload := dataFactory.Model
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if payload.Properties.Encryption != nil {
					return tf.ImportAsExistsError("azurerm_data_factory_customer_managed_key", id.ID())
				}
			}

			keyVaultKey, err := keyvault.ParseNestedItemID(customerManagedKey.CustomerManagedKeyID, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}
			encryption := &factories.EncryptionConfiguration{
				VaultBaseURL: keyVaultKey.KeyVaultBaseURL,
				KeyName:      keyVaultKey.Name,
				KeyVersion:   pointer.To(keyVaultKey.Version),
			}

			if identityId := customerManagedKey.UserAssignedIdentityID; identityId != "" {
				encryption.Identity = &factories.CMKIdentityDefinition{
					UserAssignedIdentity: pointer.To(identityId),
				}
			}

			payload.Properties.Encryption = encryption
			if _, err = client.CreateOrUpdate(ctx, *id, *payload, factories.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
