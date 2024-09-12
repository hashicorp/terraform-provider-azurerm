// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/tombuildsstuff/kermit/sdk/keyvault/7.4/keyvault"
)

var _ sdk.DataSource = KeyVaultSecretVersionsDataSource{}

type KeyVaultSecretVersionsDataSource struct{}

type KeyVaultSecretVersionsDataSourceModel struct {
	Name       string               `tfschema:"name"`
	KeyVaultId string               `tfschema:"key_vault_id"`
	Versions   []secretVersionModel `tfschema:"versions"`
}

type secretVersionModel struct {
	ID             string `tfschema:"id"`
	CreatedDate    string `tfschema:"created_date"`
	Enabled        bool   `tfschema:"enabled"`
	NotBeforeDate  string `tfschema:"not_before_date"`
	ExpirationDate string `tfschema:"expiration_date"`
	UpdatedDate    string `tfschema:"updated_date"`
}

func (r KeyVaultSecretVersionsDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: keyVaultValidate.NestedItemName,
		},

		"key_vault_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: keyVaultValidate.NestedItemId,
		},

		"max_results": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  25,
		},
	}
}

func (r KeyVaultSecretVersionsDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"versions": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"created_date": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"enabled": {
						Type:     schema.TypeBool,
						Computed: true,
					},

					"not_before_date": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"expiration_date": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"updated_date": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r KeyVaultSecretVersionsDataSource) ResourceType() string {
	return "azurerm_key_vault_secret_versions"
}

func (r KeyVaultSecretVersionsDataSource) ModelObject() interface{} {
	return &KeyVaultSecretVersionsDataSourceModel{}
}

func (r KeyVaultSecretVersionsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			keyVaultsClient := metadata.Client.KeyVault.VaultsClient
			client := metadata.Client.KeyVault.ManagementClient
			name := metadata.ResourceData.Get("name").(string)
			maxResults := metadata.ResourceData.Get("max_results").(int32)
			keyVaultId, err := commonids.ParseKeyVaultID(metadata.ResourceData.Get("key_vault_id").(string))
			if err != nil {
				return err
			}

			keyVaultBaseUri, err := keyVaultsClient.Get(ctx, *keyVaultId)
			if err != nil {
				return fmt.Errorf("looking up Secret %q vault url from id %q: %+v", name, *keyVaultId, err)
			}

			resp, err := client.GetSecretVersions(ctx, *keyVaultBaseUri.Model.Properties.VaultUri, name, &maxResults)
			if err != nil {
				return fmt.Errorf("making List Versions request on Azure KeyVault Secret %s: %+v", name, err)
			}

			var versions []map[string]interface{}

			if resp.Values() != nil {
				for resp.NotDone() {
					for _, v := range resp.Values() {
						versions = append(versions, expandSecretVersion(&v))
					}
					err = resp.NextWithContext(ctx)
					if err != nil {
						return fmt.Errorf("iterating over Secret Versions: %+v", err)
					}
				}
			}

			metadata.ResourceData.Set("versions", versions)

			return nil
		},
	}
}

func expandSecretVersion(v *keyvault.SecretItem) map[string]interface{} {
	item := make(map[string]interface{})
	item["id"] = *v.ID
	item["created_date"] = time.Time(*v.Attributes.Created).Format(time.RFC3339)
	item["enabled"] = *v.Attributes.Enabled
	if notBefore := v.Attributes.NotBefore; notBefore != nil {
		item["not_before_date"] = time.Time(*notBefore).Format(time.RFC3339)
	}
	if expires := v.Attributes.Expires; expires != nil {
		item["expiration_date"] = time.Time(*expires).Format(time.RFC3339)
	}
	if updated := v.Attributes.Updated; updated != nil {
		item["updated_date"] = time.Time(*updated).Format(time.RFC3339)
	}

	return item
}
