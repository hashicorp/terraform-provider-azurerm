// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_contacts

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			vaultClient := metadata.Client.KeyVault
			client := metadata.Client.KeyVault.ManagementClient
			var state KeyVaultCertificateContactsResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			keyVaultId, err := commonids.ParseKeyVaultID(state.KeyVaultId)
			if err != nil {
				return fmt.Errorf("parsing `key_vault_id`, %+v", err)
			}

			keyVaultBaseUri, err := vaultClient.BaseUriForKeyVault(ctx, *keyVaultId)
			if err != nil {
				return fmt.Errorf("looking up Base URI for Key Vault Certificate Contacts from %s: %+v", *keyVaultId, err)
			}

			id, err := parse.NewCertificateContactsID(*keyVaultBaseUri)
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetCertificateContacts(ctx, *keyVaultBaseUri)
				if err != nil {
					if !response.WasNotFound(existing.Response.Response) {
						return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", *keyVaultBaseUri, err)
					}
				}

				if !response.WasNotFound(existing.Response.Response) {
					if existing.ContactList != nil && len(*existing.ContactList) != 0 {
						return tf.ImportAsExistsError(r.ResourceType(), id.ID())
					}
				}
			}

			contacts := keyvault.Contacts{
				ContactList: expandKeyVaultCertificateContactsContact(state.Contact),
			}

			if _, err := client.SetCertificateContacts(ctx, *keyVaultBaseUri, contacts); err != nil {
				return fmt.Errorf("creating Key Vault Certificate Contacts %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandKeyVaultCertificateContactsContact(input []Contact) *[]keyvault.Contact {
	results := make([]keyvault.Contact, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		results = append(results, keyvault.Contact{
			EmailAddress: pointer.To(item.Email),
			Name:         pointer.To(item.Name),
			Phone:        pointer.To(item.Phone),
		})
	}

	return &results
}
