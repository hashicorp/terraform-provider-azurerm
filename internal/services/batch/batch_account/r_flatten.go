// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_account

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceBatchAccountFlatten(ctx context.Context, client *batchaccount.BatchAccountClient, d *pluginsdk.ResourceData, id *batchaccount.BatchAccountId, model *batchaccount.BatchAccount, includeResource bool) error {
	d.Set("name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		ident, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", ident); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("account_endpoint", props.AccountEndpoint)
			if autoStorage := props.AutoStorage; autoStorage != nil {
				d.Set("storage_account_id", autoStorage.StorageAccountId)
				d.Set("storage_account_authentication_mode", string(pointer.From(autoStorage.AuthenticationMode)))

				if autoStorage.NodeIdentityReference != nil {
					d.Set("storage_account_node_identity", autoStorage.NodeIdentityReference.ResourceId)
				}
			} else {
				d.Set("storage_account_authentication_mode", "")
				d.Set("storage_account_id", "")
			}

			if v := props.PublicNetworkAccess; v != nil {
				d.Set("public_network_access_enabled", *v == batchaccount.PublicNetworkAccessTypeEnabled)
			}

			if err := d.Set("network_profile", flattenBatchAccountNetworkProfile(props.NetworkProfile)); err != nil {
				return fmt.Errorf("setting `network_profile`: %+v", err)
			}

			d.Set("pool_allocation_mode", string(pointer.From(props.PoolAllocationMode)))

			if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}

			if err := d.Set("allowed_authentication_modes", flattenAllowedAuthenticationModes(props.AllowedAuthenticationModes)); err != nil {
				return fmt.Errorf("setting `allowed_authentication_modes`: %+v", err)
			}
			if includeResource {
				if d.Get("pool_allocation_mode").(string) == string(batchaccount.PoolAllocationModeBatchService) &&
					isShardKeyAllowed(d.Get("allowed_authentication_modes").(*pluginsdk.Set).List()) {
					keys, err := client.GetKeys(ctx, *id)
					if err != nil {
						return fmt.Errorf("cannot read keys for Batch account %s: %v", *id, err)
					}

					if keysModel := keys.Model; keysModel != nil {
						d.Set("primary_access_key", keysModel.Primary)
						d.Set("secondary_access_key", keysModel.Secondary)
					}
				}
			}
			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenAllowedAuthenticationModes(input *[]batchaccount.AuthenticationMode) []string {
	if input == nil || len(*input) == 0 {
		return []string{}
	}

	allowedAuthModes := make([]string, 0)
	for _, mode := range *input {
		allowedAuthModes = append(allowedAuthModes, string(mode))
	}
	return allowedAuthModes
}

func flattenEncryption(encryptionProperties *batchaccount.EncryptionProperties) []interface{} {
	if encryptionProperties == nil || *encryptionProperties.KeySource == batchaccount.KeySourceMicrosoftPointBatch {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"key_vault_key_id": *encryptionProperties.KeyVaultProperties.KeyIdentifier,
		},
	}
}

func flattenBatchAccountNetworkProfile(input *batchaccount.NetworkProfile) []interface{} {
	if input == nil || input.AccountAccess == nil && input.NodeManagementAccess == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"account_access":         flattenBatchAccountEndpointAccessProfile(input.AccountAccess),
			"node_management_access": flattenBatchAccountEndpointAccessProfile(input.NodeManagementAccess),
		},
	}
}

func flattenBatchAccountEndpointAccessProfile(input *batchaccount.EndpointAccessProfile) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	if input.IPRules != nil {
		for _, ipRule := range *input.IPRules {
			flattenedIpRule := map[string]interface{}{
				"action":   string(ipRule.Action),
				"ip_range": ipRule.Value,
			}
			ipRules = append(ipRules, flattenedIpRule)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"default_action": string(input.DefaultAction),
			"ip_rule":        ipRules,
		},
	}
}

func flattenBatchAccountKeyvaultReference(keyVaultReference *batchaccount.KeyVaultReference) interface{} {
	result := make(map[string]interface{})

	if keyVaultReference == nil {
		return []interface{}{}
	}

	if keyVaultReference.Id != "" {
		result["id"] = keyVaultReference.Id
	}

	if keyVaultReference.Url != "" {
		result["url"] = keyVaultReference.Url
	}

	return []interface{}{result}
}
