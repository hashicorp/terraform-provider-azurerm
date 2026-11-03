// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_key

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/deletedkeys"
	"github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7-4/keys"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKeyVaultKeyCreate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Key %q vault url from id %q: %+v", name, *keyVaultId, err)
	}
	client := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.Keys.Clone(*keyVaultBaseUri)
	keyId := keys.NewKeyID(*keyVaultBaseUri, name)
	keyVersionId := keys.NewKeyversionID(keyId.BaseURI, keyId.KeyName, "")

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetKey(ctx, keyVersionId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Key %q (Key Vault %q): %s", name, *keyVaultBaseUri, err)
			}
		}

		if model := existing.Model; model != nil && model.Key != nil && pointer.From(model.Key.Kid) != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_key", *model.Key.Kid)
		}
	}

	// TODO: support Importing Keys once this is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/1747
	parameters := keys.KeyCreateParameters{
		Kty:    keys.JsonWebKeyType(d.Get("key_type").(string)),
		KeyOps: expandKeyVaultKeyOptions(d),
		Attributes: &keys.KeyAttributes{
			Enabled: pointer.To(true),
		},
		Tags: pointer.To(tags.ToTypedObject(tags.Expand(d.Get("tags").(map[string]any)))),
	}

	if p := expandKeyVaultKeyReleasePolicy(d.Get("release_policy").([]any)); p != nil {
		parameters.ReleasePolicy = p
		// key must be exportable when `release_policy` is set
		parameters.Attributes.Exportable = pointer.To(true)
	}

	switch parameters.Kty {
	case keys.JsonWebKeyTypeEC, keys.JsonWebKeyTypeECNegativeHSM:
		parameters.Crv = pointer.ToEnum[keys.JsonWebKeyCurveName](d.Get("curve").(string))
	case keys.JsonWebKeyTypeRSA, keys.JsonWebKeyTypeRSANegativeHSM:
		keySize, ok := d.GetOk("key_size")
		if !ok {
			return errors.New("key_size is required when creating an RSA key")
		}
		parameters.KeySize = pointer.To(int64(keySize.(int)))
	}
	// TODO: support `oct` once this is fixed
	// https://github.com/Azure/azure-rest-api-specs/issues/1739#issuecomment-332236257

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Nbf = pointer.To(notBeforeDate.Unix())
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.Attributes.Exp = pointer.To(expirationDate.Unix())
	}

	if resp, err := client.CreateKey(ctx, keyId, parameters); err != nil {
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedKeys && response.WasConflict(resp.HttpResponse) {
			deletedKeysClient := meta.(*clients.Client).KeyVault.DataPlaneKeyVaultClient.DeletedKeys.Clone(*keyVaultBaseUri)
			deletedKeyId := deletedkeys.NewDeletedkeyID(*keyVaultBaseUri, name)
			recoveredKey, err := deletedKeysClient.RecoverDeletedKey(ctx, deletedKeyId)
			if err != nil {
				return err
			}
			if recoveredKey.Model != nil && recoveredKey.Model.Key != nil {
				kid := recoveredKey.Model.Key.Kid
				log.Printf("[DEBUG] Recovering Key %q with ID: %q", name, pointer.From(kid))
				if kid != nil {
					stateConf := &pluginsdk.StateChangeConf{
						Pending:                   []string{"pending"},
						Target:                    []string{"available"},
						Refresh:                   helpers.KeyVaultChildItemRefreshFunc(ctx, *kid),
						Delay:                     30 * time.Second,
						PollInterval:              10 * time.Second,
						ContinuousTargetOccurence: 10,
						Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
					}

					if _, err := stateConf.WaitForStateContext(ctx); err != nil {
						return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
					}
					log.Printf("[DEBUG] Key %q recovered with ID: %q", name, *kid)
				}
			}
		} else {
			return fmt.Errorf("creating Key: %+v", err)
		}
	}

	if v, ok := d.GetOk("rotation_policy"); ok {
		if respPolicy, err := client.UpdateKeyRotationPolicy(ctx, keyId, expandKeyVaultKeyRotationPolicy(v.([]any))); err != nil {
			if response.WasForbidden(respPolicy.HttpResponse) {
				return fmt.Errorf("current client lacks permissions to create Key Rotation Policy for Key %q (%q, Vault url: %q), please update this as described here: %s : %v", name, *keyVaultId, *keyVaultBaseUri, "https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_key#example-usage", err)
			}
			return fmt.Errorf("creating Key Rotation Policy: %+v", err)
		}
	}

	// "" indicates the latest version
	read, err := client.GetKey(ctx, keyVersionId)
	if err != nil {
		return err
	}

	if read.Model == nil || read.Model.Key == nil || read.Model.Key.Kid == nil {
		return fmt.Errorf("cannot read KeyVault Key '%s' (in key vault '%s')", name, *keyVaultBaseUri)
	}
	itemId, err := keyvault.ParseNestedItemID(*read.Model.Key.Kid, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
	if err != nil {
		return err
	}
	d.SetId(itemId.ID())

	return resourceKeyVaultKeyRead(d, meta)
}
