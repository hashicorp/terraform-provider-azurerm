// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
)

func flattenIoTHubDPSSku(input iotdpsresource.IotDpsSkuInfo) []any {
	output := make(map[string]any)

	name := ""
	if input.Name != nil {
		name = string(*input.Name)
	}
	output["name"] = name

	if capacity := input.Capacity; capacity != nil {
		output["capacity"] = int(*capacity)
	}

	return []any{output}
}

func flattenIoTHubDPSLinkedHub(input *[]iotdpsresource.IotHubDefinitionDescription) []any {
	linkedHubs := make([]any, 0)
	if input == nil {
		return linkedHubs
	}

	for _, attr := range *input {
		linkedHub := make(map[string]any)

		if attr.Name != nil {
			linkedHub["hostname"] = *attr.Name
		}
		if attr.ApplyAllocationPolicy != nil {
			linkedHub["apply_allocation_policy"] = *attr.ApplyAllocationPolicy
		}
		if attr.AllocationWeight != nil {
			linkedHub["allocation_weight"] = *attr.AllocationWeight
		}

		linkedHub["connection_string"] = attr.ConnectionString
		linkedHub["location"] = location.Normalize(attr.Location)

		linkedHubs = append(linkedHubs, linkedHub)
	}

	return linkedHubs
}

func flattenDpsIPFilterRules(in *[]iotdpsresource.IPFilterRule) []any {
	rules := make([]any, 0)
	if in == nil {
		return rules
	}

	for _, r := range *in {
		rawRule := make(map[string]any)

		rawRule["name"] = r.FilterName
		rawRule["action"] = string(r.Action)
		rawRule["ip_mask"] = r.IPMask

		if r.Target != nil && *r.Target != "" {
			rawRule["target"] = string(*r.Target)
		}

		rules = append(rules, rawRule)
	}
	return rules
}
