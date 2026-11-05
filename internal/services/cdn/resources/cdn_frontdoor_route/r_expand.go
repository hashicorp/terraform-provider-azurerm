// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandCdnFrontDoorRouteResourceReference(input string) *routes.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &routes.ResourceReference{
		Id: pointer.To(input),
	}
}

func expandCdnFrontDoorRouteEnabled(input bool) *routes.EnabledState {
	if input {
		return pointer.To(routes.EnabledStateEnabled)
	}

	return pointer.To(routes.EnabledStateDisabled)
}

func expandCdnFrontDoorRouteHttpsRedirect(input bool) *routes.HTTPSRedirect {
	if input {
		return pointer.To(routes.HTTPSRedirectEnabled)
	}

	return pointer.To(routes.HTTPSRedirectDisabled)
}

func expandCdnFrontDoorRouteDefaultDomain(input bool) *routes.LinkToDefaultDomain {
	if input {
		return pointer.To(routes.LinkToDefaultDomainEnabled)
	}

	return pointer.To(routes.LinkToDefaultDomainDisabled)
}

func expandCdnFrontDoorRouteEndpointProtocolsArray(input []any) *[]routes.AFDEndpointProtocols {
	results := make([]routes.AFDEndpointProtocols, 0)

	for _, item := range input {
		results = append(results, routes.AFDEndpointProtocols(item.(string)))
	}

	return &results
}

func expandCdnFrontdoorRouteRuleSetReferenceArray(input []any) *[]routes.ResourceReference {
	// NOTE: The Frontdoor service, do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors
	// when they attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]routes.ResourceReference, 0)
	for _, item := range input {
		results = append(results, routes.ResourceReference{
			Id: pointer.To(item.(string)),
		})
	}

	return &results
}

func expandCdnFrontdoorRouteCacheConfiguration(input []any) *routes.AfdRouteCacheConfiguration {
	// NOTE: If this is not an explicit nil you will receive an "Unsupported QueryStringCachingBehavior type:
	// Property 'RouteV2.CacheConfiguration.QueryStringCachingBehavior' is required but it was not set" error.
	// The Frontdoor service treats empty slices as if they are fully defined unlike other services.
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)

	cacheConfiguration := &routes.AfdRouteCacheConfiguration{
		CompressionSettings: &routes.CompressionSettings{
			IsCompressionEnabled: pointer.To(v["compression_enabled"].(bool)),
		},
		QueryParameters:            helpers.ExpandStringSliceToCsvFormat(v["query_strings"].([]any)),
		QueryStringCachingBehavior: pointer.ToEnum[routes.AfdQueryStringCachingBehavior](v["query_string_caching_behavior"].(string)),
	}

	if contentTypes := v["content_types_to_compress"].([]any); len(contentTypes) > 0 {
		cacheConfiguration.CompressionSettings.ContentTypesToCompress = pluginsdk.ExpandStringSlice(contentTypes)
	}

	return cacheConfiguration
}
