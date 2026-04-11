// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SearchDatasourceEncryptionKeyModel struct {
	KeyName           string `tfschema:"key_name"`
	KeyVersion        string `tfschema:"key_version"`
	KeyVaultUri       string `tfschema:"key_vault_uri"`
	ApplicationId     string `tfschema:"application_id"`
	ApplicationSecret string `tfschema:"application_secret"`
}

func SearchDatasourceEncryptionKeySchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"key_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: keyVaultValidate.NestedItemName,
				},

				"key_vault_uri": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsURLWithHTTPS,
				},

				"key_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"application_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
					RequiredWith: []string{"encryption_key.0.application_secret"},
				},

				"application_secret": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Sensitive:    true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{"encryption_key.0.application_id"},
				},
			},
		},
	}
}

func ExpandSearchDatasourceEncryptionKey(input []SearchDatasourceEncryptionKeyModel) *datasources.SearchResourceEncryptionKey {
	if len(input) == 0 {
		return nil
	}

	ek := input[0]
	result := &datasources.SearchResourceEncryptionKey{
		KeyVaultKeyName:    ek.KeyName,
		KeyVaultKeyVersion: ek.KeyVersion,
		KeyVaultUri:        ek.KeyVaultUri,
	}

	if ek.ApplicationId != "" {
		result.AccessCredentials = &datasources.AzureActiveDirectoryApplicationCredentials{
			ApplicationId:     ek.ApplicationId,
			ApplicationSecret: pointer.To(ek.ApplicationSecret),
		}
	}

	return result
}

func FlattenSearchDatasourceEncryptionKey(input *datasources.SearchResourceEncryptionKey, d *pluginsdk.ResourceData) []SearchDatasourceEncryptionKeyModel {
	if input == nil {
		return []SearchDatasourceEncryptionKeyModel{}
	}

	ekModel := SearchDatasourceEncryptionKeyModel{
		KeyName:     input.KeyVaultKeyName,
		KeyVersion:  input.KeyVaultKeyVersion,
		KeyVaultUri: input.KeyVaultUri,
	}

	if ac := input.AccessCredentials; ac != nil {
		ekModel.ApplicationId = ac.ApplicationId
		if v, ok := d.GetOk("encryption_key.0.application_secret"); ok {
			ekModel.ApplicationSecret = v.(string)
		}
	}

	return []SearchDatasourceEncryptionKeyModel{ekModel}
}
