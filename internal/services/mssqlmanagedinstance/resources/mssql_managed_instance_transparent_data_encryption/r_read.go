// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_transparent_data_encryption

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceencryptionprotectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlManagedInstanceTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQLManagedInstance.ManagedInstanceEncryptionProtectorClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedInstanceEncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	managedInstanceId := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

	resp, err := client.Get(ctx, managedInstanceId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("managed_instance_id", managedInstanceId.ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			var key *keyvault.NestedItemID
			if props.ServerKeyType == managedinstanceencryptionprotectors.ServerKeyTypeAzureKeyVault && props.Uri != nil {
				key, err = keyvault.ParseNestedItemID(*props.Uri, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
				if err != nil {
					return err
				}
			}

			var keyVaultKeyId string
			if key != nil {
				keyVaultKeyId = key.ID()
			}

			if err := d.Set("key_vault_key_id", keyVaultKeyId); err != nil {
				return fmt.Errorf("setting `key_vault_key_id`: %+v", err)
			}

			if err := d.Set("auto_rotation_enabled", pointer.From(props.AutoRotationEnabled)); err != nil {
				return fmt.Errorf("setting `auto_rotation_enabled`: %+v", err)
			}
		}
	}
	return nil
}

func resourceMsSqlManagedInstanceTransparentDataEncryptionKeyVaultName(keyVaultURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(keyVaultURL)
	if err != nil {
		return "", fmt.Errorf("parsing Key Vault URL (%s): %+v", keyVaultURL, err)
	}

	return strings.Split(parsedURL.Host, ".")[0], nil
}

// If versionless KV key is used in config/state, suppress the diff with versioned key.
// This is needed because backend API will return back a versioned key if it's given versionless one
func diffSuppressKeyVaultVersionedKey(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if d == nil || !d.Get("auto_rotation_enabled").(bool) {
		return false
	}

	newId, err := keyvault.ParseNestedItemID(newValue, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return false
	}

	if newId.Version != "" {
		return false
	}

	oldId, err := keyvault.ParseNestedItemID(oldValue, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return false
	}

	return strings.EqualFold(oldId.KeyVaultBaseURL, newId.KeyVaultBaseURL) && strings.EqualFold(oldId.Name, newId.Name)
}

func resourceMsSqlManagedInstanceTransparentDataEncryptionVersionedKey(ctx context.Context, id *keyvault.NestedItemID, fieldName string, autoRotationEnabled bool, meta any) (*keyvault.NestedItemID, error) {
	if id.Version != "" {
		return id, nil
	}

	if !autoRotationEnabled {
		return nil, fmt.Errorf("`%s` must be a versioned Key Vault Key ID when `auto_rotation_enabled` is false", fieldName)
	}

	var keyUrl string
	if id.IsManagedHSM() {
		resp, err := meta.(*clients.Client).ManagedHSMs.DataPlaneKeysClient.GetKey(ctx, id.KeyVaultBaseURL, id.Name, "")
		if err != nil {
			return nil, err
		}

		if resp.Key != nil {
			keyUrl = pointer.From(resp.Key.Kid)
		}
	} else {
		resp, err := meta.(*clients.Client).KeyVault.ManagementClient.GetKey(ctx, id.KeyVaultBaseURL, id.Name, "")
		if err != nil {
			return nil, err
		}

		if resp.Key != nil {
			keyUrl = pointer.From(resp.Key.Kid)
		}
	}

	if keyUrl == "" {
		return nil, fmt.Errorf("retrieving latest version for Key Vault Key %q at %q: empty key ID returned", id.Name, id.KeyVaultBaseURL)
	}

	versionedId, err := keyvault.ParseNestedItemID(keyUrl, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, err
	}

	return versionedId, nil
}
