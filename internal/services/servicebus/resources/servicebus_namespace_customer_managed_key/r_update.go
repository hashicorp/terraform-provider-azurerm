// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_customer_managed_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceBus.NamespacesClient
			var cmk ServiceBusNamespaceCustomerManagedKeyModel

			if err := metadata.Decode(&cmk); err != nil {
				return err
			}

			id, err := namespaces.ParseNamespaceID(cmk.NamespaceID)
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", *id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` is nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` is nil", *id)
			}

			if resp.Model.Properties.Encryption == nil || resp.Model.Properties.Encryption.KeyVaultProperties == nil || len(*resp.Model.Properties.Encryption.KeyVaultProperties) == 0 {
				return fmt.Errorf("retrieving %s: Customer Managed Key was not found", *id)
			}

			payload := resp.Model

			if metadata.ResourceData.HasChange("key_vault_key_id") {
				keyId, err := keyvault.ParseNestedItemID(cmk.KeyVaultKeyID, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
				if err != nil {
					return err
				}

				(*payload.Properties.Encryption.KeyVaultProperties)[0].KeyName = pointer.To(keyId.Name)
				(*payload.Properties.Encryption.KeyVaultProperties)[0].KeyVersion = pointer.To(keyId.Version)
				(*payload.Properties.Encryption.KeyVaultProperties)[0].KeyVaultUri = pointer.To(keyId.KeyVaultBaseURL)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating Customer Managed Key for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
