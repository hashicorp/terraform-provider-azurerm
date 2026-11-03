// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_secret

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	kv "github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

func resourceKeyVaultSecretCreate(d *pluginsdk.ResourceData, meta any) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUrl, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetSecret(ctx, *keyVaultBaseUrl, name, "")
		if err != nil {
			if !response.WasNotFound(existing.Response.Response) {
				return fmt.Errorf("checking for presence of existing Secret %q (Key Vault %q): %s", name, *keyVaultBaseUrl, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_key_vault_secret", *existing.ID)
		}
	}

	value := d.Get("value").(string)

	valueWo, err := pluginsdk.GetWriteOnly(d, "value_wo", cty.String)
	if err != nil {
		return err
	}
	if !valueWo.IsNull() {
		value = valueWo.AsString()
	}

	contentType := d.Get("content_type").(string)
	t := d.Get("tags").(map[string]any)

	parameters := kv.SecretSetParameters{
		Value:            pointer.To(value),
		ContentType:      pointer.To(contentType),
		Tags:             tags.Expand(t),
		SecretAttributes: &kv.SecretAttributes{},
	}

	if v, ok := d.GetOk("not_before_date"); ok {
		notBeforeDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.SecretAttributes.NotBefore = pointer.To(date.UnixTime(notBeforeDate))
	}

	if v, ok := d.GetOk("expiration_date"); ok {
		expirationDate, _ := time.Parse(time.RFC3339, v.(string)) // validated by schema
		parameters.SecretAttributes.Expires = pointer.To(date.UnixTime(expirationDate))
	}

	if resp, err := client.SetSecret(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
		// In the case that the Secret already exists in a Soft Deleted / Recoverable state we check if `recover_soft_deleted_key_vaults` is set
		// and attempt recovery where appropriate
		if meta.(*clients.Client).Features.KeyVault.RecoverSoftDeletedSecrets && response.WasConflict(resp.Response.Response) {
			recoveredSecret, err := client.RecoverDeletedSecret(ctx, *keyVaultBaseUrl, name)
			if err != nil {
				return err
			}
			log.Printf("[DEBUG] Recovering Secret %q with ID: %q", name, *recoveredSecret.ID)
			// We need to wait for consistency, recovered Key Vault Child items are not as readily available as newly created
			if secret := recoveredSecret.ID; secret != nil {
				stateConf := &pluginsdk.StateChangeConf{
					Pending:                   []string{"pending"},
					Target:                    []string{"available"},
					Refresh:                   helpers.KeyVaultChildItemRefreshFunc(ctx, *secret),
					Delay:                     30 * time.Second,
					PollInterval:              10 * time.Second,
					ContinuousTargetOccurence: 10,
					Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for Key Vault Secret %q to become available: %s", name, err)
				}
				log.Printf("[DEBUG] Secret %q recovered with ID: %q", name, *recoveredSecret.ID)

				if _, err := client.SetSecret(ctx, *keyVaultBaseUrl, name, parameters); err != nil {
					return err
				}
			}
		} else {
			// If the error response was anything else, or `recover_soft_deleted_key_vaults` is `false` just return the error
			return err
		}
	}

	// "" indicates the latest version
	read, err := client.GetSecret(ctx, *keyVaultBaseUrl, name, "")
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read KeyVault Secret '%s' (in key vault '%s')", name, *keyVaultBaseUrl)
	}

	secretId, err := keyvault.ParseNestedItemID(*read.ID, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeSecret)
	if err != nil {
		return err
	}

	d.SetId(secretId.ID())

	return resourceKeyVaultSecretRead(d, meta)
}
