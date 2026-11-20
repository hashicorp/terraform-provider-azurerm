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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.CertificateContactsID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			vaultClient := metadata.Client.KeyVault
			client := metadata.Client.KeyVault.DataPlaneKeyVaultClient.Certificates.Clone(id.KeyVaultBaseUrl)
			subscriptionId := metadata.Client.Account.SubscriptionId

			subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
			keyVaultIdRaw, err := vaultClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving resource ID of the Key Vault at URL %s: %+v", id.KeyVaultBaseUrl, err)
			}
			if keyVaultIdRaw == nil {
				metadata.Logger.Infof("Unable to determine the Resource ID for the Key Vault at URL %s - removing from state!", id.KeyVaultBaseUrl)
				return metadata.MarkAsGone(id)
			}
			keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
			if err != nil {
				return fmt.Errorf("parsing Key Vault ID: %+v", err)
			}

			existing, err := client.GetCertificateContacts(ctx)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					metadata.Logger.Infof("No Certificate Contacts could be found at %s - removing from state!", id.KeyVaultBaseUrl)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("checking for presence of existing Certificate Contacts (Key Vault %q): %s", id.KeyVaultBaseUrl, err)
			}

			var contacts *[]certificates.Contact
			if existing.Model != nil {
				contacts = existing.Model.Contacts
			}

			state := KeyVaultCertificateContactsResourceModel{
				KeyVaultId: keyVaultId.ID(),
				Contact:    flattenKeyVaultCertificateContactsContact(contacts),
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func flattenKeyVaultCertificateContactsContact(input *[]certificates.Contact) []Contact {
	result := make([]Contact, 0)
	if input == nil {
		return result
	}

	for _, item := range *input {
		result = append(result, Contact{
			Email: pointer.From(item.Email),
			Name:  pointer.From(item.Name),
			Phone: pointer.From(item.Phone),
		})
	}

	return result
}

type Contact struct {
	Email string `tfschema:"email"`
	Name  string `tfschema:"name"`
	Phone string `tfschema:"phone"`
}
