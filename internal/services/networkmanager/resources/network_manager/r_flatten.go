// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
)

func flattenNetworkManagerScope(input networkmanagers.NetworkManagerPropertiesNetworkManagerScopes) []ManagerScopeModel {
	return []ManagerScopeModel{{
		ManagementGroups: pointer.From(input.ManagementGroups),
		Subscriptions:    pointer.From(input.Subscriptions),
	}}
}

func flattenNetworkManagerScopeAccesses(input *[]networkmanagers.ConfigurationType) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenNetworkManagerCrossTenantScopes(input *[]networkmanagers.CrossTenantScopes) []ManagerCrossTenantScopeModel {
	if input == nil {
		return make([]ManagerCrossTenantScopeModel, 0)
	}

	results := make([]ManagerCrossTenantScopeModel, 0, len(*input))
	for _, v := range *input {
		results = append(results, ManagerCrossTenantScopeModel{
			TenantId:         pointer.From(v.TenantId),
			ManagementGroups: pointer.From(v.ManagementGroups),
			Subscriptions:    pointer.From(v.Subscriptions),
		})
	}
	return results
}
