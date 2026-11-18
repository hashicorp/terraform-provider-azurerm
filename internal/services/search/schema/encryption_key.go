// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SearchDatasourceEncryptionKeyModel struct {
	KeyVaultKeyId string `tfschema:"key_vault_key_id"`
	ClientId      string `tfschema:"client_id"`
	ClientSecret  string `tfschema:"client_secret"`
}

func SearchDatasourceEncryptionKeySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"key_vault_key_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
				},

				"client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
					RequiredWith: []string{"encryption_key.0.client_secret"},
				},

				"client_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{"encryption_key.0.client_id"},
				},
			},
		},
	}
}

func ExpandSearchDatasourceEncryptionKey(input []SearchDatasourceEncryptionKeyModel) (*datasources.SearchResourceEncryptionKey, error) {
	if len(input) == 0 {
		return nil, nil
	}

	ek := input[0]
	keyVaultKeyId, err := keyvault.ParseNestedItemID(ek.KeyVaultKeyId, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
	if err != nil {
		return nil, err
	}

	result := &datasources.SearchResourceEncryptionKey{
		KeyVaultKeyName:    keyVaultKeyId.Name,
		KeyVaultKeyVersion: keyVaultKeyId.Version,
		KeyVaultUri:        keyVaultKeyId.KeyVaultBaseURL,
	}
	if ek.ClientId != "" {
		result.AccessCredentials = &datasources.AzureActiveDirectoryApplicationCredentials{
			ApplicationId:     ek.ClientId,
			ApplicationSecret: pointer.To(ek.ClientSecret),
		}
	}

	return result, nil
}

func FlattenSearchDatasourceEncryptionKey(input *datasources.SearchResourceEncryptionKey, d *pluginsdk.ResourceData) ([]SearchDatasourceEncryptionKeyModel, error) {
	if input == nil {
		return []SearchDatasourceEncryptionKeyModel{}, nil
	}

	keyVaultKeyId, err := keyvault.NewNestedItemID(input.KeyVaultUri, keyvault.NestedItemTypeKey, input.KeyVaultKeyName, input.KeyVaultKeyVersion)
	if err != nil {
		return []SearchDatasourceEncryptionKeyModel{}, err
	}

	ekModel := SearchDatasourceEncryptionKeyModel{
		KeyVaultKeyId: keyVaultKeyId.ID(),
	}
	if ac := input.AccessCredentials; ac != nil {
		ekModel.ClientId = ac.ApplicationId
		if v, ok := d.GetOk("encryption_key.0.client_secret"); ok {
			ekModel.ClientSecret = v.(string)
		}
	}

	return []SearchDatasourceEncryptionKeyModel{ekModel}, nil
}
