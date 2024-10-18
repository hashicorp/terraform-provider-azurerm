// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"fmt"
	"sort"
	"strings"
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
	MaxResults int64                `tfschema:"max_results"`
	Versions   []secretVersionModel `tfschema:"versions"`
}

type secretVersionModel struct {
	ID             string `tfschema:"id"`
	CreatedDate    string `tfschema:"created_date"`
	Enabled        bool   `tfschema:"enabled"`
	NotBeforeDate  string `tfschema:"not_before_date"`
	ExpirationDate string `tfschema:"expiration_date"`
	UpdatedDate    string `tfschema:"updated_date"`
	Uri            string `tfschema:"uri"`
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
			ValidateFunc: commonids.ValidateKeyVaultID,
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
					"uri": {
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
			client := metadata.Client.KeyVault.ManagementClient

			var model KeyVaultSecretVersionsDataSourceModel

			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			keyVaultId, err := commonids.ParseKeyVaultID(model.KeyVaultId)
			if err != nil {
				return err
			}

			keyVaultUri, err := metadata.Client.KeyVault.BaseUriForKeyVault(ctx, *keyVaultId)
			if err != nil {
				return err
			}

			maxResults32 := int32(model.MaxResults)

			resp, err := client.GetSecretVersions(ctx, *keyVaultUri, model.Name, &maxResults32)
			if err != nil {
				return fmt.Errorf("making List Versions request on Azure KeyVault Secret %s: %+v", model.Name, err)
			}

			if resp.Values() != nil {
				for resp.NotDone() {
					for _, v := range resp.Values() {
						model.Versions = append(model.Versions, expandSecretVersion(&v))
					}
					err = resp.NextWithContext(ctx)
					if err != nil {
						return fmt.Errorf("iterating over Secret Versions: %+v", err)
					}
				}

				var errors []error
				model.Versions, errors = sortSecretVersions(model.Versions)
				if len(errors) > 0 {
					return fmt.Errorf("sorting Secret Versions: %+v", errors)
				}
			}

			metadata.ResourceData.SetId(fmt.Sprintf("%s/%s", model.KeyVaultId, model.Name))

			return metadata.Encode(&model)
		},
	}
}

func expandSecretVersion(v *keyvault.SecretItem) secretVersionModel {
	var item secretVersionModel
	item.Uri = *v.ID
	item.ID = (*v.ID)[strings.LastIndex(*v.ID, "/")+1:]
	item.CreatedDate = time.Time(*v.Attributes.Created).Format(time.RFC3339)
	item.Enabled = *v.Attributes.Enabled
	if notBefore := v.Attributes.NotBefore; notBefore != nil {
		item.NotBeforeDate = time.Time(*notBefore).Format(time.RFC3339)
	}
	if expires := v.Attributes.Expires; expires != nil {
		item.ExpirationDate = time.Time(*expires).Format(time.RFC3339)
	}
	if updated := v.Attributes.Updated; updated != nil {
		item.UpdatedDate = time.Time(*updated).Format(time.RFC3339)
	}

	return item
}

func sortSecretVersions(values []secretVersionModel) ([]secretVersionModel, []error) {
	errors := make([]error, 0)
	sort.Slice(values, func(i, j int) bool {
		// Sort by CreatedDate in descending order
		timeA, err := time.Parse(time.RFC3339, values[i].CreatedDate)
		if err != nil {
			errors = append(errors, err)
			return false
		}

		timeB, err := time.Parse(time.RFC3339, values[j].CreatedDate)
		if err != nil {
			errors = append(errors, err)
			return false
		}

		return timeA.After(timeB)

	})

	if len(errors) > 0 {
		return values, errors
	}

	return values, nil
}
