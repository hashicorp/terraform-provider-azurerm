// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_provider

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
)

func flattenCustomProviderResourceType(input *[]customresourceprovider.CustomRPResourceTypeRouteDefinition) []any {
	if input == nil {
		return []any{}
	}

	definitions := make([]any, 0)
	for _, v := range *input {
		definition := make(map[string]any)

		definition["routing_type"] = v.RoutingType
		definition["name"] = v.Name
		definition["endpoint"] = v.Endpoint

		definitions = append(definitions, definition)
	}
	return definitions
}

func flattenCustomProviderAction(input *[]customresourceprovider.CustomRPActionRouteDefinition) []any {
	if input == nil {
		return []any{}
	}

	definitions := make([]any, 0)
	for _, v := range *input {
		definition := make(map[string]any)

		definition["name"] = v.Name
		definition["endpoint"] = v.Endpoint

		definitions = append(definitions, definition)
	}
	return definitions
}

func flattenCustomProviderValidation(input *[]customresourceprovider.CustomRPValidations) []any {
	if input == nil {
		return []any{}
	}

	validations := make([]any, 0)
	for _, v := range *input {
		validation := make(map[string]any)

		validation["specification"] = v.Specification

		validations = append(validations, validation)
	}
	return validations
}
