// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_issuer

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultCertificateIssuerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).KeyVault.ManagementClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IssuerID(d.Id())
	if err != nil {
		return err
	}
	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	if keyVaultIdRaw == nil {
		log.Printf("[DEBUG] Unable to determine the Resource ID for the Key Vault at URL %q - removing from state!", id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if %s for Certificate %q exists: %v", *keyVaultId, id.Name, err)
	}
	if !ok {
		log.Printf("[DEBUG] Certificate %q was not found %s - removing from state", id.Name, *keyVaultId)
		d.SetId("")
		return nil
	}

	resp, err := client.GetCertificateIssuer(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[DEBUG] KeyVault Certificate Issuer %q (KeyVault URI %q) does not exist - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("failed to make Read request on Azure KeyVault Certificate Issuer %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)

	if resp.Provider != nil {
		d.Set("provider_name", resp.Provider)
	}
	if resp.OrganizationDetails != nil {
		if resp.OrganizationDetails.ID != nil {
			d.Set("org_id", resp.OrganizationDetails.ID)
		}
		d.Set("admin", flattenKeyVaultCertificateIssuerAdmins(resp.OrganizationDetails.AdminDetails))
	}
	if resp.Credentials != nil {
		d.Set("account_id", resp.Credentials.AccountID)
	}

	return nil
}

func flattenKeyVaultCertificateIssuerAdmins(input *[]keyvault.AdministratorDetails) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, admin := range *input {
		results = append(results, map[string]any{
			"email_address": pointer.From(admin.EmailAddress),
			"first_name":    pointer.From(admin.FirstName),
			"last_name":     pointer.From(admin.LastName),
			"phone":         pointer.From(admin.Phone),
		})
	}

	return results
}
