// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_account_encryption

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.AccountClient

			id, err := netappaccounts.ParseNetAppAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppAccountEncryption
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.AccountsGet(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model.Properties.Encryption == nil {
				return fmt.Errorf("encryption information does not exist for %s", id)
			}

			anfAccountIdentityFlattened, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(existing.Model.Identity)
			if err != nil {
				return err
			}

			encryptionKey, federatedClientID, err := flattenEncryption(existing.Model.Properties.Encryption)
			if err != nil {
				return err
			}

			model := models.NetAppAccountEncryption{
				NetAppAccountID:   id.ID(),
				EncryptionKey:     encryptionKey,
				FederatedClientID: federatedClientID,
			}

			// Populate cross-tenant key vault resource ID only for cross-tenant scenarios (when federated_client_id is present)
			if federatedClientID != "" && existing.Model.Properties.Encryption.KeyVaultProperties != nil && existing.Model.Properties.Encryption.KeyVaultProperties.KeyVaultResourceId != nil {
				model.CrossTenantKeyVaultResourceID = pointer.From(existing.Model.Properties.Encryption.KeyVaultProperties.KeyVaultResourceId)
			}

			if len(anfAccountIdentityFlattened) > 0 {
				if anfAccountIdentityFlattened[0].Type == identity.TypeSystemAssigned {
					model.SystemAssignedIdentityPrincipalID = anfAccountIdentityFlattened[0].PrincipalId
				}

				if anfAccountIdentityFlattened[0].Type == identity.TypeUserAssigned {
					if len(anfAccountIdentityFlattened[0].IdentityIds) > 0 {
						model.UserAssignedIdentityID = anfAccountIdentityFlattened[0].IdentityIds[0]
					}
				}
			}

			metadata.SetID(id)

			return metadata.Encode(pointer.To(model))
		},
	}
}

func flattenEncryption(encryptionProperties *netappaccounts.AccountEncryption) (string, string, error) {
	if encryptionProperties == nil || pointer.From(encryptionProperties.KeySource) == netappaccounts.KeySourceMicrosoftPointNetApp || encryptionProperties.KeyVaultProperties == nil {
		return "", "", nil
	}
	props := encryptionProperties.KeyVaultProperties

	keyVaultKeyId, err := keyvault.NewNestedItemID(props.KeyVaultUri, keyvault.NestedItemTypeKey, props.KeyName, "")
	if err != nil {
		return "", "", fmt.Errorf("parsing key vault key id: %+v", err)
	}

	federatedClientID := ""
	if encryptionProperties.Identity != nil && encryptionProperties.Identity.FederatedClientId != nil {
		federatedClientID = pointer.From(encryptionProperties.Identity.FederatedClientId)
	}

	return keyVaultKeyId.VersionlessID(), federatedClientID, nil
}
