// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appconfiguration

import (
	"bytes"
	"context"
	"fmt"
	"log"

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
		log.Printf("[DEBUG] Refresh App Configuration status")
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

func appConfigurationGetKeyRefreshFuncForUpdate(ctx context.Context, client *appconfiguration.BaseClient, key, label string, model appconfiguration.KeyValue) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refresh App Configuration to ensure all properties are synced")
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

		if !appConfigurationKeyValueEuqals(res, model) {
			return "Syncing", "Syncing", nil
		}

		return res, "Synced", nil
	}
}

func appConfigurationKeyValueEuqals(kv1, kv2 appconfiguration.KeyValue) bool {
	if (kv1.ContentType == nil) != (kv2.ContentType == nil) || pointer.From(kv1.ContentType) != pointer.From(kv2.ContentType) {
		log.Printf("[DEBUG] Syncing App Configuration Key `content_type`: one with value %q, another with value %q", pointer.From(kv1.ContentType), pointer.From(kv2.ContentType))
		return false
	}

	if (kv1.Locked == nil) != (kv2.Locked == nil) || pointer.From(kv1.Locked) != pointer.From(kv2.Locked) {
		log.Printf("[DEBUG] Syncing App Configuration Key `locked`: one with value %q, another with value %q", pointer.From(kv1.ContentType), pointer.From(kv2.ContentType))
		return false
	}

	if (kv1.Value == nil) != (kv2.Value == nil) || pointer.From(kv1.Value) != pointer.From(kv2.Value) {
		log.Printf("[DEBUG] Syncing App Configuration Key `value` field: one with value %q, another with value %q", pointer.From(kv1.Value), pointer.From(kv2.Value))
		return false
	}

	if (kv1.Tags == nil) != (kv2.Tags == nil) || len(kv1.Tags) != len(kv2.Tags) {
		log.Printf("[DEBUG] Syncing App Configuration Key `tags` field: one with length %q, another with length %q", len(kv1.Tags), len(kv2.Tags))
		return false
	}

	for k, v := range kv1.Tags {
		if v != nil {
			v2, ok := kv2.Tags[k]
			if !ok || v2 == nil || *v != *v2 {
				log.Printf("[DEBUG] Syncing App Configuration Key `tags` field: one with value %q, another with value %q", pointer.From(v), pointer.From(v2))
				return false
			}
		}
	}

	return true
}
