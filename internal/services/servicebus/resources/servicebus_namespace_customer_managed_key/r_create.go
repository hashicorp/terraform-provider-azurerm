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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if resp.Model.Properties != nil && resp.Model.Properties.Encryption != nil && resp.Model.Properties.Encryption.KeyVaultProperties != nil && len(*resp.Model.Properties.Encryption.KeyVaultProperties) > 0 {
					return metadata.ResourceRequiresImport(r.ResourceType(), *id)
				}
			}

			payload := resp.Model

			keyId, err := keyvault.ParseNestedItemID(cmk.KeyVaultKeyID, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
			if err != nil {
				return err
			}

			payload.Properties.Encryption = &namespaces.Encryption{
				RequireInfrastructureEncryption: pointer.To(cmk.InfrastructureEncryptionEnabled),
				KeySource:                       pointer.To(namespaces.KeySourceMicrosoftPointKeyVault),
				KeyVaultProperties: &[]namespaces.KeyVaultProperties{
					{
						KeyName:     pointer.To(keyId.Name),
						KeyVersion:  pointer.To(keyId.Version),
						KeyVaultUri: pointer.To(keyId.KeyVaultBaseURL),
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, *id, *payload, metadata.SetIDAndIdentityCallback(id)); err != nil {
				return fmt.Errorf("creating Customer Managed Key for %s: %+v", *id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return nil
		},
	}
}
