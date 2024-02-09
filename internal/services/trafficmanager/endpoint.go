// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package trafficmanager

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandEndpointCustomHeaderConfig(input []interface{}) *[]endpoints.EndpointPropertiesCustomHeadersInlined {
	output := make([]endpoints.EndpointPropertiesCustomHeadersInlined, 0)

	for _, header := range input {
		headerBlock := header.(map[string]interface{})
		output = append(output, endpoints.EndpointPropertiesCustomHeadersInlined{
			Name:  utils.String(headerBlock["name"].(string)),
			Value: utils.String(headerBlock["value"].(string)),
		})
	}

	return &output
}

func flattenEndpointCustomHeaderConfig(input *[]endpoints.EndpointPropertiesCustomHeadersInlined) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}
	for _, header := range *input {
		name := ""
		if header.Name != nil {
			name = *header.Name
		}

		value := ""
		if header.Value != nil {
			value = *header.Value
		}
		result = append(result, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}
	return result
}

func expandEndpointSubnetConfig(input []interface{}) *[]endpoints.EndpointPropertiesSubnetsInlined {
	output := make([]endpoints.EndpointPropertiesSubnetsInlined, 0)

	for _, subnet := range input {
		subnetBlock := subnet.(map[string]interface{})
		if subnetBlock["scope"].(int) == 0 && subnetBlock["first"].(string) != "0.0.0.0" {
			output = append(output, endpoints.EndpointPropertiesSubnetsInlined{
				First: utils.String(subnetBlock["first"].(string)),
				Last:  utils.String(subnetBlock["last"].(string)),
			})
		} else {
			output = append(output, endpoints.EndpointPropertiesSubnetsInlined{
				First: utils.String(subnetBlock["first"].(string)),
				Scope: utils.Int64(int64(subnetBlock["scope"].(int))),
			})
		}
	}

	return &output
}

func flattenEndpointSubnetConfig(input *[]endpoints.EndpointPropertiesSubnetsInlined) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}
	for _, subnet := range *input {
		first := ""
		if subnet.First != nil {
			first = *subnet.First
		}

		last := ""
		if subnet.Last != nil {
			last = *subnet.Last
		}
		scope := 0
		if subnet.Scope != nil {
			scope = int(*subnet.Scope)
		}
		result = append(result, map[string]interface{}{
			"first": first,
			"last":  last,
			"scope": scope,
		})
	}
	return result
}
