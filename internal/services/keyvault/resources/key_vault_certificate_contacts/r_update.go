// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_contacts

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KeyVault.DataPlaneKeyVaultClient.Certificates

			var state KeyVaultCertificateContactsResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id, err := parse.CertificateContactsID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Clone(id.KeyVaultBaseUrl).GetCertificateContacts(ctx)
			if err != nil {
				return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", id.KeyVaultBaseUrl, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving Certificate Contacts (Key Vault %q): `model` was nil", id.KeyVaultBaseUrl)
			}
			contacts := *existing.Model

			if metadata.ResourceData.HasChange("contact") {
				contacts.Contacts = expandKeyVaultCertificateContactsContact(state.Contact)
			}

			if _, err := client.Clone(id.KeyVaultBaseUrl).SetCertificateContacts(ctx, contacts); err != nil {
				return fmt.Errorf("updating Key Vault Certificate Contacts %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
