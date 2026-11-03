// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_account_encryption

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient
			keyVaultsClient := metadata.Client.KeyVault
			subscriptionId := commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId)

			var model models.NetAppAccountEncryption
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.NetAppAccountID == "" {
				return fmt.Errorf("netapp_account_id is empty")
			}

			accountID, err := netappaccounts.ParseNetAppAccountID(model.NetAppAccountID)
			if err != nil {
				return fmt.Errorf("error parsing account id %s: %+v", model.NetAppAccountID, err)
			}

			locks.ByID(accountID.ID())
			defer locks.UnlockByID(accountID.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.AccountsGet(ctx, pointer.From(accountID))
				if err != nil {
					if response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("not found %s: %s", accountID.ID(), err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					if existing.Model.Properties.Encryption != nil && existing.Model.Properties.Encryption.KeySource != nil && pointer.From(existing.Model.Properties.Encryption.KeySource) == netappaccounts.KeySourceMicrosoftPointKeyVault {
						return tf.ImportAsExistsError(r.ResourceType(), accountID.ID())
					}
				}
			}

			update := netappaccounts.NetAppAccountPatch{
				Properties: &netappaccounts.AccountProperties{},
			}

			encryptionExpanded, err := expandEncryption(ctx, model.EncryptionKey, keyVaultsClient, subscriptionId, pointer.To(model))
			if err != nil {
				return fmt.Errorf("error expanding encryption: %+v", err)
			}

			update.Properties.Encryption = encryptionExpanded

			if err := client.AccountsUpdateCallbackThenPoll(ctx, pointer.From(accountID), update, metadata.SetIDCallback(accountID)); err != nil {
				return fmt.Errorf("updating %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)

			return nil
		},
	}
}

func expandEncryption(ctx context.Context, input string, keyVaultsClient *keyVaultClient.Client, subscriptionID commonids.SubscriptionId, model *models.NetAppAccountEncryption) (*netappaccounts.AccountEncryption, error) {
	encryptionProperty := netappaccounts.AccountEncryption{
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointNetApp),
	}

	if input == "" {
		return &encryptionProperty, nil
	}

	keyId, err := keyvault.ParseNestedItemID(input, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, fmt.Errorf("parsing `key_vault_key_id` %q: %+v", input, err)
	}

	var keyVaultResourceID string

	// For cross-tenant scenarios: using CrossTenantKeyVaultResourceID
	if model.FederatedClientID != "" {
		keyVaultResourceID = model.CrossTenantKeyVaultResourceID
	} else {
		// Same-tenant scenario: lookup the key vault ID
		keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionID, keyId.KeyVaultBaseURL)
		if err != nil {
			return nil, fmt.Errorf("retrieving the resource id the key vault at url %q: %s", keyId.KeyVaultBaseURL, err)
		}

		if keyVaultID == nil {
			return nil, fmt.Errorf("keyVaultID is nil for key vault url %q", keyId.KeyVaultBaseURL)
		}

		keyVaultResourceID = pointer.From(keyVaultID)
	}

	encryptionIdentity := &netappaccounts.EncryptionIdentity{}

	if model.UserAssignedIdentityID != "" {
		encryptionIdentity = &netappaccounts.EncryptionIdentity{
			UserAssignedIdentity: pointer.To(model.UserAssignedIdentityID),
		}
		if model.FederatedClientID != "" {
			encryptionIdentity.FederatedClientId = pointer.To(model.FederatedClientID)
		}
	}

	encryptionProperty = netappaccounts.AccountEncryption{
		Identity:  encryptionIdentity,
		KeySource: pointer.To(netappaccounts.KeySourceMicrosoftPointKeyVault),
		KeyVaultProperties: &netappaccounts.KeyVaultProperties{
			KeyName:            keyId.Name,
			KeyVaultUri:        keyId.KeyVaultBaseURL,
			KeyVaultResourceId: pointer.To(keyVaultResourceID),
		},
	}

	return &encryptionProperty, nil
}
