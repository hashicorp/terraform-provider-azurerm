// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_transparent_data_encryption

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/encryptionprotectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlTransparentDataEncryptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.EncryptionProtectorClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EncryptionProtectorID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.Get(ctx, serverId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("server_id", serverId.ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			var key *keyvault.NestedItemID
			if props.ServerKeyType == encryptionprotectors.ServerKeyTypeAzureKeyVault && props.Uri != nil {
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

func resourceMsSqlTransparentDataEncryptionKeyVaultName(keyVaultURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(keyVaultURL)
	if err != nil {
		return "", fmt.Errorf("parsing Key Vault URL (%s): %+v", keyVaultURL, err)
	}

	return strings.Split(parsedURL.Host, ".")[0], nil
}
