// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_contacts

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			vaultClient := metadata.Client.KeyVault
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

			client := metadata.Client.KeyVault.DataPlaneKeyVaultClient.Certificates.Clone(*keyVaultBaseUri)

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetCertificateContacts(ctx)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", *keyVaultBaseUri, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					if existing.Model != nil && existing.Model.Contacts != nil && len(*existing.Model.Contacts) != 0 {
						return tf.ImportAsExistsError(r.ResourceType(), id.ID())
					}
				}
			}

			contacts := certificates.Contacts{
				Contacts: expandKeyVaultCertificateContactsContact(state.Contact),
			}

			if _, err := client.SetCertificateContacts(ctx, contacts); err != nil {
				return fmt.Errorf("creating Key Vault Certificate Contacts %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandKeyVaultCertificateContactsContact(input []Contact) *[]certificates.Contact {
	results := make([]certificates.Contact, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		results = append(results, certificates.Contact{
			Email: pointer.To(item.Email),
			Name:  pointer.To(item.Name),
			Phone: pointer.To(item.Phone),
		})
	}

	return &results
}
