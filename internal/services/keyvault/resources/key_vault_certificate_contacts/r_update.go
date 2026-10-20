// Copyright IBM Corp. 2014, 2025
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
			client := metadata.Client.KeyVault.ManagementClient

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

			existing, err := client.GetCertificateContacts(ctx, id.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", id.KeyVaultBaseUrl, err)
			}

			if metadata.ResourceData.HasChange("contact") {
				existing.ContactList = expandKeyVaultCertificateContactsContact(state.Contact)
			}

			if _, err := client.SetCertificateContacts(ctx, id.KeyVaultBaseUrl, existing); err != nil {
				return fmt.Errorf("updating Key Vault Certificate Contacts %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
