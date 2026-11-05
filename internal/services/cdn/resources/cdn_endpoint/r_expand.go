// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandCdnEndpointGeoFilters(d *pluginsdk.ResourceData) *[]cdn.GeoFilter {
	filters := make([]cdn.GeoFilter, 0)

	inputFilters := d.Get("geo_filter").([]any)
	for _, v := range inputFilters {
		input := v.(map[string]any)
		action := input["action"].(string)
		relativePath := input["relative_path"].(string)

		inputCountryCodes := input["country_codes"].([]any)
		countryCodes := make([]string, 0)

		for _, v := range inputCountryCodes {
			if v != nil {
				countryCode := v.(string)
				countryCodes = append(countryCodes, countryCode)
			}
		}

		filter := cdn.GeoFilter{
			Action:       cdn.GeoFilterActions(action),
			RelativePath: pointer.To(relativePath),
			CountryCodes: &countryCodes,
		}
		filters = append(filters, filter)
	}

	return &filters
}

func expandArmCdnEndpointContentTypesToCompress(d *pluginsdk.ResourceData) []string {
	results := make([]string, 0)
	input := d.Get("content_types_to_compress").(*pluginsdk.Set).List()

	for _, v := range input {
		contentType := v.(string)
		results = append(results, contentType)
	}

	return results
}

func expandAzureRmCdnEndpointOrigins(d *pluginsdk.ResourceData) []cdn.DeepCreatedOrigin {
	configs := d.Get("origin").(*pluginsdk.Set).List()
	origins := make([]cdn.DeepCreatedOrigin, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]any)

		name := data["name"].(string)
		hostName := data["host_name"].(string)

		origin := cdn.DeepCreatedOrigin{
			Name: pointer.To(name),
			DeepCreatedOriginProperties: &cdn.DeepCreatedOriginProperties{
				HostName: pointer.To(hostName),
			},
		}

		if v, ok := data["https_port"]; ok {
			port := v.(int)
			origin.HTTPSPort = pointer.To(int32(port))
		}

		if v, ok := data["http_port"]; ok {
			port := v.(int)
			origin.HTTPPort = pointer.To(int32(port))
		}

		origins = append(origins, origin)
	}

	return origins
}

func expandArmCdnEndpointDeliveryPolicy(globalRulesRaw []any, deliveryRulesRaw []any) (*cdn.EndpointPropertiesUpdateParametersDeliveryPolicy, error) {
	deliveryRules := make([]cdn.DeliveryRule, 0)
	deliveryPolicy := cdn.EndpointPropertiesUpdateParametersDeliveryPolicy{
		Description: pointer.To(""),
		Rules:       &deliveryRules,
	}

	if len(globalRulesRaw) > 0 && globalRulesRaw[0] != nil {
		ruleRaw := globalRulesRaw[0].(map[string]any)
		rule, err := expandArmCdnEndpointGlobalDeliveryRule(ruleRaw)
		if err != nil {
			return nil, err
		}
		deliveryRules = append(deliveryRules, *rule)
	}

	for _, ruleV := range deliveryRulesRaw {
		ruleRaw := ruleV.(map[string]any)
		rule, err := helpers.ExpandArmCdnEndpointDeliveryRule(ruleRaw)
		if err != nil {
			return nil, err
		}
		deliveryRules = append(deliveryRules, *rule)
	}

	return &deliveryPolicy, nil
}
