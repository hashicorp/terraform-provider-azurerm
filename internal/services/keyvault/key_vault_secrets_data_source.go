// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

func dataSourceKeyVaultSecrets() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKeyVaultSecretsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"key_vault_id": commonschema.ResourceIDReferenceRequired(commonids.KeyVaultId{}),

			"names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"secrets": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceKeyVaultSecretsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).KeyVault.ManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultId, err := commonids.ParseKeyVaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return fmt.Errorf("fetching base vault url from id %q: %+v", *keyVaultId, err)
	}

	secretList, err := client.GetSecretsComplete(ctx, *keyVaultBaseUri, utils.Int32(25))
	if err != nil {
		return fmt.Errorf("making Read request on Azure KeyVault %q: %+v", *keyVaultId, err)
	}

	d.SetId(keyVaultId.ID())

	var names []string
	var secrets []map[string]interface{}

	if secretList.Response().Value != nil {
		for secretList.NotDone() {
			for _, v := range *secretList.Response().Value {
				name, err := parseNameFromSecretUrl(*v.ID)
				if err != nil {
					return err
				}
				names = append(names, *name)
				secrets = append(secrets, expandSecrets(*name, v))
				err = secretList.NextWithContext(ctx)
				if err != nil {
					return fmt.Errorf("listing secrets on Azure KeyVault %q: %+v", *keyVaultId, err)
				}
			}
		}
	}

	d.Set("names", names)
	d.Set("secrets", secrets)
	d.Set("key_vault_id", keyVaultId.ID())

	return nil
}

func parseNameFromSecretUrl(input string) (*string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	// https://favoretti-keyvault.vault.azure.net/secrets/secret-name
	segments := strings.Split(uri.Path, "/")
	if len(segments) != 3 {
		return nil, fmt.Errorf("expected a Path in the format `/secrets/secret-name` but got %q", uri.Path)
	}
	return &segments[2], nil
}

func expandSecrets(name string, item keyvault.SecretItem) map[string]interface{} {
	res := map[string]interface{}{
		"id":   *item.ID,
		"name": name,
	}
	if item.Attributes != nil && item.Attributes.Enabled != nil {
		res["enabled"] = *item.Attributes.Enabled
	}
	return res
}
