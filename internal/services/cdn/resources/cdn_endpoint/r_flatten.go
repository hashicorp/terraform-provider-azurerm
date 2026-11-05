// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
)

func flattenCdnEndpointGeoFilters(input *[]cdn.GeoFilter) []any {
	results := make([]any, 0)

	if filters := input; filters != nil {
		for _, filter := range *filters {
			outputCodes := make([]any, 0)
			if codes := filter.CountryCodes; codes != nil {
				for _, code := range *codes {
					outputCodes = append(outputCodes, code)
				}
			}

			results = append(results, map[string]any{
				"action":        string(filter.Action),
				"country_codes": outputCodes,
				"relative_path": pointer.From(filter.RelativePath),
			})
		}
	}

	return results
}

func flattenAzureRMCdnEndpointContentTypes(input *[]string) []any {
	output := make([]any, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}

func flattenAzureRMCdnEndpointOrigin(input *[]cdn.DeepCreatedOrigin) []any {
	results := make([]any, 0)

	if list := input; list != nil {
		for _, i := range *list {
			hostName := ""
			httpPort := 80
			httpsPort := 443
			if props := i.DeepCreatedOriginProperties; props != nil {
				if props.HostName != nil {
					hostName = *props.HostName
				}
				if port := props.HTTPPort; port != nil {
					httpPort = int(*port)
				}
				if port := props.HTTPSPort; port != nil {
					httpsPort = int(*port)
				}
			}

			results = append(results, map[string]any{
				"name":       pointer.From(i.Name),
				"host_name":  hostName,
				"http_port":  httpPort,
				"https_port": httpsPort,
			})
		}
	}

	return results
}

func flattenEndpointDeliveryPolicy(input *cdn.EndpointPropertiesUpdateParametersDeliveryPolicy) (*flattenedEndpointDeliveryPolicies, error) {
	output := flattenedEndpointDeliveryPolicies{
		globalDeliveryRules: make([]any, 0),
		deliveryRules:       make([]any, 0),
	}
	if input == nil || input.Rules == nil {
		return &output, nil
	}

	for _, rule := range *input.Rules {
		if rule.Order == nil {
			continue
		}

		if int(*rule.Order) == 0 {
			flattenedRule, err := flattenArmCdnEndpointGlobalDeliveryRule(rule)
			if err != nil {
				return nil, err
			}

			output.globalDeliveryRules = append(output.globalDeliveryRules, flattenedRule)
			continue
		}

		flattenedRule, err := helpers.FlattenArmCdnEndpointDeliveryRule(rule)
		if err != nil {
			return nil, err
		}

		output.deliveryRules = append(output.deliveryRules, flattenedRule)
	}

	return &output, nil
}

type flattenedEndpointDeliveryPolicies struct {
	globalDeliveryRules []any
	deliveryRules       []any
}
