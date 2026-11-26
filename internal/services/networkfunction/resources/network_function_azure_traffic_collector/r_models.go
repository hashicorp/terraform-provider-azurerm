// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_azure_traffic_collector

type NetworkFunctionAzureTrafficCollectorModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	CollectorPolicies []string          `tfschema:"collector_policy_ids"`
	VirtualHub        []string          `tfschema:"virtual_hub_id"`
}

func (r Resource) ModelObject() any {
	return &NetworkFunctionAzureTrafficCollectorModel{}
}
