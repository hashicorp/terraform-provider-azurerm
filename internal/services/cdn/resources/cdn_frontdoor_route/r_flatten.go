// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_route

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenCdnFrontDoorRouteRuleSetResourceArray(input *[]routes.ResourceReference) ([]any, error) {
	results := make([]any, 0)
	if input == nil {
		return results, nil
	}

	// Normalize these values in the configuration file we know they are valid because they were set on the
	// resource... if these are modified in the portal they will all be lowercased...
	// Issue: https://github.com/Azure/azure-sdk-for-go/issues/19378
	for _, ruleSet := range *input {
		id, err := rulesets.ParseRuleSetIDInsensitively(pointer.From(ruleSet.Id))
		if err != nil {
			return results, err
		}
		results = append(results, id.ID())
	}

	return results, nil
}

func flattenCdnFrontDoorRouteOriginGroupResourceReference(input *routes.ResourceReference) (string, error) {
	if input != nil && input.Id != nil {
		// The 5.0 go-azure-sdk migration changed this read path to strict parsing, regressing support for inconsistently cased IDs returned by Azure; see #32953.
		id, err := afdorigingroups.ParseOriginGroupIDInsensitively(*input.Id)
		if err != nil {
			return "", err
		}

		return id.ID(), nil
	}

	return "", nil
}

func flattenCdnFrontDoorRouteCacheConfiguration(input *routes.AfdRouteCacheConfiguration) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	compressionEnabled := false
	contentTypesToCompress := make([]any, 0)
	if v := input.CompressionSettings; v != nil {
		compressionEnabled = pointer.From(v.IsCompressionEnabled)
		contentTypesToCompress = pluginsdk.FlattenSlice(v.ContentTypesToCompress)
	}

	return []any{
		map[string]any{
			"compression_enabled":           compressionEnabled,
			"content_types_to_compress":     contentTypesToCompress,
			"query_string_caching_behavior": pointer.FromEnum(input.QueryStringCachingBehavior),
			"query_strings":                 helpers.FlattenCsvToStringSlice(input.QueryParameters),
		},
	}
}
