// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenAccessPolicies(policies *[]resource.ServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		result = append(result, policy.ObjectId)
	}

	return result
}

func flattenAuthentication(input *resource.ServiceAuthenticationConfigurationInfo) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"audience":            pointer.From(input.Audience),
			"authority":           pointer.From(input.Authority),
			"smart_proxy_enabled": pointer.From(input.SmartProxyEnabled),
		},
	}
}

func flattenCorsConfig(input *resource.ServiceCorsConfigurationInfo) []any {
	if input == nil {
		return []any{}
	}

	maxAge := 0
	if input.MaxAge != nil {
		maxAge = int(*input.MaxAge)
	}
	return []any{
		map[string]any{
			"allow_credentials":  pointer.From(input.AllowCredentials),
			"allowed_headers":    pluginsdk.FlattenSlice(input.Headers),
			"allowed_methods":    pluginsdk.FlattenSlice(input.Methods),
			"allowed_origins":    pluginsdk.FlattenSlice(input.Origins),
			"max_age_in_seconds": maxAge,
		},
	}
}
