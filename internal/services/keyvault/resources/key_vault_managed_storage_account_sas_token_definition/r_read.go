// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_managed_storage_account_sas_token_definition

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultManagedStorageAccountSasTokenDefinitionRead(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SasDefinitionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, id.KeyVaultBaseUrl)
	if err != nil {
		return fmt.Errorf("retrieving the Resource ID of the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}

	keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
	if err != nil {
		return err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Managed Storage Account Sas Definition %q (Key Vault %q): %+v", id.Name, id.KeyVaultBaseUrl, err)
	}
	if !ok {
		log.Printf("[DEBUG] Managed Storage Account Sas Definition %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseUrl)
		d.SetId("")
		return nil
	}

	client := keyVaultsClient.DataPlaneKeyVaultClient.Storage.Clone(id.KeyVaultBaseUrl)
	resp, err := client.GetSasDefinition(ctx, storage.NewSaID(id.KeyVaultBaseUrl, id.StorageAccountName, id.Name))
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Managed Storage Account Sas Definition %q was not found in Key Vault at URI %q - removing from state", id.Name, id.KeyVaultBaseUrl)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("cannot read Managed Storage Account Sas Definition %s: %+v", id.Name, err)
	}

	d.Set("name", id.Name)

	if model := resp.Model; model != nil {
		d.Set("sas_template_uri", model.TemplateUri)
		d.Set("sas_type", pointer.FromEnum(model.SasType))
		d.Set("secret_id", model.Sid)
		d.Set("validity_period", model.ValidityPeriod)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
