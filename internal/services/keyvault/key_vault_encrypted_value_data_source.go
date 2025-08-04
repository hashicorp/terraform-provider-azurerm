// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/kermit/sdk/keyvault/7.4/keyvault"
)

var _ sdk.DataSource = EncryptedValueDataSource{}

type EncryptedValueDataSource struct{}

type EncryptedValueDataSourceModel struct {
	KeyVaultKeyId         string `tfschema:"key_vault_key_id"`
	Algorithm             string `tfschema:"algorithm"`
	EncryptedData         string `tfschema:"encrypted_data"`
	PlainTextValue        string `tfschema:"plain_text_value"`
	DecodedPlainTextValue string `tfschema:"decoded_plain_text_value"`
}

func (EncryptedValueDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key_vault_key_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validate.NestedItemId,
		},
		"algorithm": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(keyvault.JSONWebKeyEncryptionAlgorithmRSA15),
				string(keyvault.JSONWebKeyEncryptionAlgorithmRSAOAEP),
				string(keyvault.JSONWebKeyEncryptionAlgorithmRSAOAEP256),
			}, false),
		},
		"encrypted_data": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"plain_text_value": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
}

func (EncryptedValueDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"decoded_plain_text_value": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (EncryptedValueDataSource) ModelObject() interface{} {
	return &EncryptedValueDataSourceModel{}
}

func (e EncryptedValueDataSource) ResourceType() string {
	return "azurerm_key_vault_encrypted_value"
}

func (EncryptedValueDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KeyVault.ManagementClient

			var model EncryptedValueDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.EncryptedData == "" && model.PlainTextValue == "" {
				return fmt.Errorf("one of `encrypted_data` or `plain_text_value` must be specified - both were empty")
			}
			if model.EncryptedData != "" && model.PlainTextValue != "" {
				return fmt.Errorf("only one of `encrypted_data` or `plain_text_value` must be specified - both were specified")
			}

			keyVaultKeyId, err := parse.ParseNestedItemID(model.KeyVaultKeyId)
			if err != nil {
				return err
			}

			if model.EncryptedData != "" {
				params := keyvault.KeyOperationsParameters{
					Algorithm: keyvault.JSONWebKeyEncryptionAlgorithm(model.Algorithm),
					Value:     utils.String(model.EncryptedData),
				}
				result, err := client.Decrypt(ctx, keyVaultKeyId.KeyVaultBaseUrl, keyVaultKeyId.Name, keyVaultKeyId.Version, params)
				if err != nil {
					return fmt.Errorf("decrypting plain-text value using Key Vault Key ID %q: %+v", model.KeyVaultKeyId, err)
				}
				if result.Result == nil {
					return fmt.Errorf("decrypting plain-text value using Key Vault Key ID %q: `result` was nil", model.KeyVaultKeyId)
				}
				model.PlainTextValue = *result.Result
				if decodedResult, err := base64.RawURLEncoding.DecodeString(*result.Result); err == nil {
					model.DecodedPlainTextValue = string(decodedResult)
				} else {
					log.Printf("[WARN] Failed to decode plain-text value: %+v", err)
				}
			} else {
				params := keyvault.KeyOperationsParameters{
					Algorithm: keyvault.JSONWebKeyEncryptionAlgorithm(model.Algorithm),
					Value:     utils.String(model.PlainTextValue),
				}
				result, err := client.Encrypt(ctx, keyVaultKeyId.KeyVaultBaseUrl, keyVaultKeyId.Name, keyVaultKeyId.Version, params)
				if err != nil {
					return fmt.Errorf("encrypting plain-text value using Key Vault Key ID %q: %+v", model.KeyVaultKeyId, err)
				}
				if result.Result == nil {
					return fmt.Errorf("encrypting plain-text value using Key Vault Key ID %q: `result` was nil", model.KeyVaultKeyId)
				}
				model.EncryptedData = *result.Result
			}

			metadata.ResourceData.SetId(fmt.Sprintf("%s-%s-%s", model.KeyVaultKeyId, model.Algorithm, sha1.Sum([]byte(model.EncryptedData))))
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}
