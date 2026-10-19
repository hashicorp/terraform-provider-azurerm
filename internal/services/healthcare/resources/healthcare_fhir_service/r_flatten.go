// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_fhir_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenFhirAcrLoginServer(acrConfig *fhirservices.FhirServiceAcrConfiguration) []string {
	result := make([]string, 0)
	if acrConfig == nil {
		return result
	}

	if loginServer := acrConfig.LoginServers; loginServer != nil {
		result = append(result, *loginServer...)
	}
	return result
}

func flattenFhirAccessPolicy(policies *[]fhirservices.FhirServiceAccessPolicyEntry) []string {
	result := make([]string, 0)

	if policies == nil {
		return result
	}

	for _, policy := range *policies {
		result = append(result, policy.ObjectId)
	}
	return result
}

func flattenOciArtifacts(artifacts *[]fhirservices.ServiceOciArtifactEntry) []map[string]any {
	result := make([]map[string]any, 0)
	if artifacts == nil {
		return result
	}
	for _, artifact := range *artifacts {
		artifactRaw := make(map[string]any)

		if loginServer := artifact.LoginServer; loginServer != nil {
			artifactRaw["login_server"] = *loginServer
		}
		if imageName := artifact.ImageName; imageName != nil {
			artifactRaw["image_name"] = *imageName
		}
		if digest := artifact.Digest; digest != nil {
			artifactRaw["digest"] = *digest
		}
		result = append(result, artifactRaw)
	}

	return result
}

func flattenFhirCorsConfiguration(corsConfig *fhirservices.FhirServiceCorsConfiguration) []any {
	if corsConfig == nil {
		return []any{}
	}

	if corsConfig.Origins != nil && len(*corsConfig.Origins) == 0 &&
		corsConfig.Methods != nil && len(*corsConfig.Methods) == 0 &&
		corsConfig.Headers != nil && len(*corsConfig.Headers) == 0 &&
		corsConfig.AllowCredentials != nil && !*corsConfig.AllowCredentials {
		return []any{}
	}

	var maxAge int
	if corsConfig.MaxAge != nil {
		maxAge = int(*corsConfig.MaxAge)
	}

	return []any{
		map[string]any{
			"credentials_allowed": pointer.From(corsConfig.AllowCredentials),
			"allowed_headers":     pluginsdk.FlattenSlice(corsConfig.Headers),
			"allowed_methods":     pluginsdk.FlattenSlice(corsConfig.Methods),
			"allowed_origins":     pluginsdk.FlattenSlice(corsConfig.Origins),
			"max_age_in_seconds":  maxAge,
		},
	}
}

func flattenFhirAuthentication(authConfig *fhirservices.FhirServiceAuthenticationConfiguration) []any {
	if authConfig == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"audience":            pointer.From(authConfig.Audience),
			"authority":           pointer.From(authConfig.Authority),
			"smart_proxy_enabled": pointer.From(authConfig.SmartProxyEnabled),
		},
	}
}
