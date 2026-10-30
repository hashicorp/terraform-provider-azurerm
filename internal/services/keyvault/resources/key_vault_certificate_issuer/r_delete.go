// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultCertificateIssuerDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IssuerID(d.Id())
	if err != nil {
		return err
	}

	// we verify it exists
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Issuer %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	if !ok {
		// parent key vault does not exist, nothing to delete
		return nil
	}

	resp, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] Certificate Issuer %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			return nil
		}

		return fmt.Errorf("checking if Issuer %q in key vault %q at url %q exists: %v", id.Name, *keyVaultId, id.KeyVaultBaseUrl, err)
	}

	_, err = client.DeleteCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	return err
}
