// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/secrets"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultSecretUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
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

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Secrets.Clone(id.KeyVaultBaseURL)
	secretId := secrets.NewSecretID(id.KeyVaultBaseURL, id.Name)
	secretVersionId := secrets.NewSecretversionID(id.KeyVaultBaseURL, id.Name, "")

	value := d.Get("value").(string)
	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]any)

	secretAttributes := &secrets.SecretAttributes{}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		secretAttributes.Nbf = pointer.To(notBeforeDate.Unix())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		secretAttributes.Exp = pointer.To(expirationDate.Unix())
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
		parameters := secrets.SecretSetParameters{
			Value:       value,
			ContentType: pointer.To(contentType),
			Tags:        tags.Expand(t),
			Attributes:  secretAttributes,
		}

		if _, err = client.SetSecret(ctx, secretId, parameters); err != nil {
			return err
		}
	} else {
		parameters := secrets.SecretUpdateParameters{
			ContentType: pointer.To(contentType),
			Tags:        tags.Expand(t),
			Attributes:  secretAttributes,
		}

		if _, err = client.UpdateSecret(ctx, secretVersionId, parameters); err != nil {
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, secretVersionId)
	if err != nil {
		return fmt.Errorf("getting Key Vault Secret %q : %+v", id.Name, err)
	}

	if read.Model == nil || read.Model.Id == nil {
		return fmt.Errorf("retrieving %s: unable to determine secret ID", secretId)
	}

	parsedSecretId, err := keyvault.ParseNestedItemID(*read.Model.Id, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	// the ID is suffixed with the secret version
	d.SetId(parsedSecretId.ID())

	return resourceKeyVaultSecretRead(d, meta)
}
