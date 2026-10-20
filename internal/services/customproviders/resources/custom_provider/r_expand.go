// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_provider

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
)

func expandCustomProviderResourceType(input []any) *[]customresourceprovider.CustomRPResourceTypeRouteDefinition {
	if len(input) == 0 {
		return nil
	}
	definitions := make([]customresourceprovider.CustomRPResourceTypeRouteDefinition, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]any)
		definitions = append(definitions, customresourceprovider.CustomRPResourceTypeRouteDefinition{
			RoutingType: pointer.ToEnum[customresourceprovider.ResourceTypeRouting](attrs["routing_type"].(string)),
			Name:        attrs["name"].(string),
			Endpoint:    attrs["endpoint"].(string),
		})
	}

	return &definitions
}

func expandCustomProviderAction(input []any) *[]customresourceprovider.CustomRPActionRouteDefinition {
	if len(input) == 0 {
		return nil
	}
	definitions := make([]customresourceprovider.CustomRPActionRouteDefinition, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]any)
		definitions = append(definitions, customresourceprovider.CustomRPActionRouteDefinition{
			Name:     attrs["name"].(string),
			Endpoint: attrs["endpoint"].(string),
		})
	}

	return &definitions
}

func expandCustomProviderValidation(input []any) *[]customresourceprovider.CustomRPValidations {
	if len(input) == 0 {
		return nil
	}

	validations := make([]customresourceprovider.CustomRPValidations, 0)

	for _, v := range input {
		if v == nil {
			continue
		}

		attrs := v.(map[string]any)
		validations = append(validations, customresourceprovider.CustomRPValidations{
			Specification: attrs["specification"].(string),
		})
	}

	return &validations
}
