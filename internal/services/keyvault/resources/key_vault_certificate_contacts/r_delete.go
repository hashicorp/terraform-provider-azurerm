// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_contacts

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KeyVault.ManagementClient

			id, err := parse.CertificateContactsID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			// check if any contacts exist in vault, if they do not then nothing to delete
			resp, err := client.GetCertificateContacts(ctx, id.KeyVaultBaseUrl)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					log.Printf("[DEBUG] Key Vault Certificate Contact %q was not found in Key Vault at URI %q - removing from state", id.ID(), id.KeyVaultBaseUrl)
					return nil
				}

				return fmt.Errorf("checking if any contacts exist in key vault at url %q: %v", id.KeyVaultBaseUrl, err)
			}

			if _, err := client.DeleteCertificateContacts(ctx, id.KeyVaultBaseUrl); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
