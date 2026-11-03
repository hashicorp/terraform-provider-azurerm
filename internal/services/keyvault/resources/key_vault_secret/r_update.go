// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultSecretUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	meta.(*clients.Client).KeyVault.AddToCache(*keyVaultId, id.KeyVaultBaseURL)

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("checking if key vault %q for Secret %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseURL, err)
	}
	if !ok {
		log.Printf("[DEBUG] Secret %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseURL)
		d.SetId("")
		return nil
	}

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]any)

	secretAttributes := &kv.SecretAttributes{}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		secretAttributes.NotBefore = pointer.To(date.UnixTime(notBeforeDate))
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		secretAttributes.Expires = pointer.To(date.UnixTime(expirationDate))
	}

	if d.HasChanges("value", "value_wo_version") {
		valueWo, err := pluginsdk.GetWriteOnly(d, "value_wo", cty.String)
		if err != nil {
			return err
		}
		if !valueWo.IsNull() {
			value = valueWo.AsString()
		}
		// for changing the value of the secret we need to create a new version
		parameters := kv.SecretSetParameters{
			Value:            pointer.To(value),
			ContentType:      pointer.To(contentType),
			Tags:             tags.Expand(t),
			SecretAttributes: secretAttributes,
		}

		if _, err = client.SetSecret(ctx, id.KeyVaultBaseURL, id.Name, parameters); err != nil {
			return err
		}
	} else {
		parameters := kv.SecretUpdateParameters{
			ContentType:      pointer.To(contentType),
			Tags:             tags.Expand(t),
			SecretAttributes: secretAttributes,
		}

		if _, err = client.UpdateSecret(ctx, id.KeyVaultBaseURL, id.Name, "", parameters); err != nil {
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, id.KeyVaultBaseURL, id.Name, "")
	if err != nil {
		return fmt.Errorf("getting Key Vault Secret %q : %+v", id.Name, err)
	}

	secretId, err := keyvault.ParseNestedItemID(*read.ID, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	// the ID is suffixed with the secret version
	d.SetId(secretId.ID())

	return resourceKeyVaultSecretRead(d, meta)
}
