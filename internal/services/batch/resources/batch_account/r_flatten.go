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

func flattenEncryption(encryptionProperties *batchaccount.EncryptionProperties) []any {
	if encryptionProperties == nil || *encryptionProperties.KeySource == batchaccount.KeySourceMicrosoftPointBatch {
		return []any{}
	}

	return []any{
		map[string]any{
			"key_vault_key_id": *encryptionProperties.KeyVaultProperties.KeyIdentifier,
		},
	}
}

func flattenBatchAccountNetworkProfile(input *batchaccount.NetworkProfile) []any {
	if input == nil || input.AccountAccess == nil && input.NodeManagementAccess == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"account_access":         flattenBatchAccountEndpointAccessProfile(input.AccountAccess),
			"node_management_access": flattenBatchAccountEndpointAccessProfile(input.NodeManagementAccess),
		},
	}
}

func flattenBatchAccountEndpointAccessProfile(input *batchaccount.EndpointAccessProfile) []any {
	if input == nil {
		return []any{}
	}

	ipRules := make([]any, 0)
	if input.IPRules != nil {
		for _, ipRule := range *input.IPRules {
			flattenedIpRule := map[string]any{
				"action":   string(ipRule.Action),
				"ip_range": ipRule.Value,
			}
			ipRules = append(ipRules, flattenedIpRule)
		}
	}

	return []any{
		map[string]any{
			"default_action": string(input.DefaultAction),
			"ip_rule":        ipRules,
		},
	}
}

func resourceBatchAccountFlatten(ctx context.Context, client *batchaccount.BatchAccountClient, d *pluginsdk.ResourceData, id *batchaccount.BatchAccountId, model *batchaccount.BatchAccount, includeResource bool) error {
	d.Set("name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
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
