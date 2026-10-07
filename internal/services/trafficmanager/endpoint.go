// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/trafficmanagers"
)

func expandEndpointCustomHeaderConfig(input []any) *[]trafficmanagers.EndpointPropertiesCustomHeadersItem {
	output := make([]trafficmanagers.EndpointPropertiesCustomHeadersItem, 0)

	for _, header := range input {
		headerBlock := header.(map[string]any)
		output = append(output, trafficmanagers.EndpointPropertiesCustomHeadersItem{
			Name:  pointer.To(headerBlock["name"].(string)),
			Value: pointer.To(headerBlock["value"].(string)),
		})
	}

	return &output
}

func flattenEndpointCustomHeaderConfig(input *[]trafficmanagers.EndpointPropertiesCustomHeadersItem) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}
	for _, header := range *input {
		name := pointer.From(header.Name)

		value := pointer.From(header.Value)
		result = append(result, map[string]any{
			"name":  name,
			"value": value,
		})
	}
	return result
}

func expandEndpointSubnetConfig(input []any) *[]trafficmanagers.EndpointPropertiesSubnetsItem {
	output := make([]trafficmanagers.EndpointPropertiesSubnetsItem, 0)

	for _, subnet := range input {
		subnetBlock := subnet.(map[string]any)
		if subnetBlock["scope"].(int) == 0 && subnetBlock["first"].(string) != "0.0.0.0" {
			output = append(output, trafficmanagers.EndpointPropertiesSubnetsItem{
				First: pointer.To(subnetBlock["first"].(string)),
				Last:  pointer.To(subnetBlock["last"].(string)),
			})
		} else {
			output = append(output, trafficmanagers.EndpointPropertiesSubnetsItem{
				First: pointer.To(subnetBlock["first"].(string)),
				Scope: pointer.To(int64(subnetBlock["scope"].(int))),
			})
		}
	}

	return &output
}

func flattenEndpointSubnetConfig(input *[]trafficmanagers.EndpointPropertiesSubnetsItem) []any {
	result := make([]any, 0)
	if input == nil {
		return result
	}
	for _, subnet := range *input {
		first := pointer.From(subnet.First)

		last := pointer.From(subnet.Last)
		scope := 0
		if subnet.Scope != nil {
			scope = int(*subnet.Scope)
		}
		result = append(result, map[string]any{
			"first": first,
			"last":  last,
			"scope": scope,
		})
	}
	return result
}
