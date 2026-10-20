// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_customer_managed_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataFactory.Factories

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

			encryption := dataFactory.Model.Properties.Encryption

			customerManagedKeyId := ""
			customerManagedKeyIdentityId := ""
			if encryption.VaultBaseURL != "" && encryption.KeyName != "" && encryption.KeyVersion != nil {
				version := pointer.From(encryption.KeyVersion)

				keyId, err := keyvault.NewNestedItemID(encryption.VaultBaseURL, keyvault.NestedItemTypeKey, encryption.KeyName, version)
				if err != nil {
					return fmt.Errorf("parsing Nested Item ID: %+v", err)
				}
				customerManagedKeyId = keyId.ID()
			}

			if encIdentity := encryption.Identity; encIdentity != nil && pointer.From(encIdentity.UserAssignedIdentity) != "" {
				parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(pointer.From(encIdentity.UserAssignedIdentity))
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", *encIdentity.UserAssignedIdentity, err)
				}
				customerManagedKeyIdentityId = parsed.ID()
			}

			state := DataFactoryCustomerManagedKeyModel{
				DataFactoryID:          id.ID(),
				CustomerManagedKeyID:   customerManagedKeyId,
				UserAssignedIdentityID: customerManagedKeyIdentityId,
			}

			return metadata.Encode(&state)
		},
	}
}
