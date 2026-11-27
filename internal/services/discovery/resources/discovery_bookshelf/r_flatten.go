// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
)

func flattenDiscoveryBookshelf(id bookshelves.BookshelfId, model *bookshelves.Bookshelf) (*DiscoveryBookshelfModel, error) {
	state := DiscoveryBookshelfModel{
		Name:              id.BookshelfName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		if props := model.Properties; props != nil {
			state.LogAnalyticsClusterId = pointer.From(props.LogAnalyticsClusterId)
			state.ManagedResourceGroupName = pointer.From(props.ManagedResourceGroup)
			state.PrivateEndpointSubnetId = pointer.From(props.PrivateEndpointSubnetId)
			state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == bookshelves.PublicNetworkAccessEnabled
			state.SearchSubnetId = pointer.From(props.SearchSubnetId)
			state.Uri = pointer.From(props.BookshelfUri)
			state.WorkloadIdentityIds = flattenDiscoveryBookshelfWorkloadIdentities(props.WorkloadIdentities)
			state.ManagedOnBehalfOfBrokerResourceIds = flattenDiscoveryBookshelfMoboBrokerResourceIds(props.ManagedOnBehalfOfConfiguration)

			if pointer.From(props.CustomerManagedKeys) == bookshelves.CustomerManagedKeysEnabled {
				customerManagedKey, err := flattenDiscoveryBookshelfCustomerManagedKey(props.KeyVaultProperties, props.WorkloadIdentities)
				if err != nil {
					return nil, err
				}
				state.CustomerManagedKey = customerManagedKey
			}
		}
	}

	return &state, nil
}

func flattenDiscoveryBookshelfCustomerManagedKey(input *bookshelves.BookshelfKeyVaultProperties, workloadIdentities *map[string]bookshelves.UserAssignedIdentity) ([]DiscoveryBookshelfCustomerManagedKey, error) {
	if input == nil {
		return []DiscoveryBookshelfCustomerManagedKey{}, nil
	}

	keyId, err := keyvault.NewNestedItemID(input.KeyVaultUri, keyvault.NestedItemTypeKey, input.KeyName, pointer.From(input.KeyVersion))
	if err != nil {
		return nil, fmt.Errorf("parsing `key_vault_key_id`: %+v", err)
	}

	userAssignedIdentityId := ""
	for id, identity := range pointer.From(workloadIdentities) {
		if strings.EqualFold(pointer.From(identity.ClientId), input.IdentityClientId) {
			parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(id)
			if err != nil {
				return nil, err
			}
			userAssignedIdentityId = parsed.ID()
			break
		}
	}

	return []DiscoveryBookshelfCustomerManagedKey{{
		KeyVaultKeyId:          keyId.ID(),
		UserAssignedIdentityId: userAssignedIdentityId,
	}}, nil
}

func flattenDiscoveryBookshelfWorkloadIdentities(input *map[string]bookshelves.UserAssignedIdentity) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}

	for id := range *input {
		if parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(id); err == nil {
			output = append(output, parsed.ID())
			continue
		}
		output = append(output, id)
	}

	return output
}

func flattenDiscoveryBookshelfMoboBrokerResourceIds(input *bookshelves.WithMoboBrokerResources) []string {
	output := make([]string, 0)
	if input == nil || input.MoboBrokerResources == nil {
		return output
	}

	for _, v := range *input.MoboBrokerResources {
		output = append(output, pointer.From(v.Id))
	}

	return output
}
