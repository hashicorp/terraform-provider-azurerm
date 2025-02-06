// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/configurationstores"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2023-03-01/replicas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

func flattenAppConfigurationEncryption(input *configurationstores.EncryptionProperties) []interface{} {
	if input == nil || input.KeyVaultProperties == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"identity_client_id":       input.KeyVaultProperties.IdentityClientId,
			"key_vault_key_identifier": input.KeyVaultProperties.KeyIdentifier,
		},
	}
}

func flattenAppConfigurationReplicas(input []replicas.Replica) ([]interface{}, error) {
	results := make([]interface{}, 0)
	for _, v := range input {
		if v.Properties == nil {
			return results, fmt.Errorf("retrieving Replica %s Properties is nil", *v.Id)
		}

		replicaId, err := replicas.ParseReplicaIDInsensitively(pointer.From(v.Id))
		if err != nil {
			return results, err
		}

		result := map[string]interface{}{
			"name":     pointer.From(v.Name),
			"location": location.Normalize(pointer.From(v.Location)),
			"endpoint": pointer.From(v.Properties.Endpoint),
			"id":       replicaId.ID(),
		}
		results = append(results, result)
	}
	return results, nil
}

func resourceConfigurationStoreReplicaHash(input interface{}) int {
	var buf bytes.Buffer
	if rawData, ok := input.(map[string]interface{}); ok {
		buf.WriteString(rawData["name"].(string))
		buf.WriteString(location.Normalize(rawData["location"].(string)))
	}
	return pluginsdk.HashString(buf.String())
}

func appConfigurationGetKeyRefreshFunc(ctx context.Context, client *appconfiguration.BaseClient, key, label string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.GetKeyValue(ctx, key, label, "", "", "", []appconfiguration.KeyValueFields{})
		if err != nil {
			if v, ok := err.(autorest.DetailedError); ok {
				if response.WasForbidden(v.Response) {
					return "Forbidden", "Forbidden", nil
				}
				if response.WasNotFound(v.Response) {
					return "NotFound", "NotFound", nil
				}
			}
			return res, "Error", nil
		}

		return res, "Exists", nil
	}
}
