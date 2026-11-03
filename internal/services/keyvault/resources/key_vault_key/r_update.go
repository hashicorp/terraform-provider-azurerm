// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultKeyUpdate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := keyvault.ParseNestedItemID(d.Id(), keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
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
		return fmt.Errorf("checking if key vault %q for Key %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseURL, err)
	}
	if !ok {
		log.Printf("[DEBUG] Key %q Key Vault %q was not found in Key Vault at URI %q - removing from state", id.Name, *keyVaultId, id.KeyVaultBaseURL)
		d.SetId("")
		return nil
	}

	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Keys.Clone(id.KeyVaultBaseURL)
	keyVersionId := keys.NewKeyversionID(id.KeyVaultBaseURL, id.Name, "")
	keyId := keys.NewKeyID(id.KeyVaultBaseURL, id.Name)

	parameters := keys.KeyUpdateParameters{
		KeyOps: expandKeyVaultKeyOptions(d),
		Attributes: &keys.KeyAttributes{
			Enabled: pointer.To(true),
		},
		Tags: pointer.To(tags.ToTypedObject(tags.Expand(d.Get("tags").(map[string]any)))),
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Nbf = pointer.To(notBeforeDate.Unix())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Exp = pointer.To(expirationDate.Unix())
	}

	if d.HasChange("release_policy") {
		parameters.ReleasePolicy = expandKeyVaultKeyReleasePolicy(d.Get("release_policy").([]any))
	}

	if _, err = client.UpdateKey(ctx, keyVersionId, parameters); err != nil {
		return err
	}

	if d.HasChange("rotation_policy") {
		if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, keyId, expandKeyVaultKeyRotationPolicy(d.Get("rotation_policy").([]any))); err != nil {
			if response.WasForbidden(respPolicy.HttpResponse) {
				return fmt.Errorf("current client lacks permissions to update Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", id.Name, *keyVaultId, id.KeyVaultBaseURL, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
			}
			return fmt.Errorf("creating Key Rotation Policy: %+v", err)
		}
	}

	return resourceKeyVaultKeyRead(d, meta)
}
