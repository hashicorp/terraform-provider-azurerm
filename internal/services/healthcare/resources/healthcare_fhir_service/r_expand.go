// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_fhir_service

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandOciArtifacts(input []any) *[]fhirservices.ServiceOciArtifactEntry {
	output := make([]fhirservices.ServiceOciArtifactEntry, 0)

	for _, artifactSet := range input {
		artifactRaw := artifactSet.(map[string]any)

		artifact := fhirservices.ServiceOciArtifactEntry{
			LoginServer: pointer.To(artifactRaw["login_server"].(string)),
		}
		if image := artifactRaw["image_name"].(string); image != "" {
			artifact.ImageName = &image
		}
		if digest := artifactRaw["digest"].(string); digest != "" {
			artifact.Digest = &digest
		}

		output = append(output, artifact)
	}

	return &output
}

func expandFhirAuthentication(input []any) *fhirservices.FhirServiceAuthenticationConfiguration {
	authConfig := input[0].(map[string]any)
	authority := authConfig["authority"].(string)
	audience := authConfig["audience"].(string)
	smartProxyEnabled := authConfig["smart_proxy_enabled"].(bool)

	return &fhirservices.FhirServiceAuthenticationConfiguration{
		Authority:         pointer.To(authority),
		Audience:          pointer.To(audience),
		SmartProxyEnabled: pointer.To(smartProxyEnabled),
	}
}

func expandAccessPolicy(input []any) *[]fhirservices.FhirServiceAccessPolicyEntry {
	if len(input) == 0 {
		return nil
	}

	accessPolicySet := make([]fhirservices.FhirServiceAccessPolicyEntry, 0)

	for _, objectId := range input {
		accessPolicyObjectId := fhirservices.FhirServiceAccessPolicyEntry{
			ObjectId: objectId.(string),
		}
		accessPolicySet = append(accessPolicySet, accessPolicyObjectId)
	}

	return &accessPolicySet
}

func expandFhirCorsConfiguration(input []any) *fhirservices.FhirServiceCorsConfiguration {
	if len(input) == 0 {
		return &fhirservices.FhirServiceCorsConfiguration{
			Origins:          &[]string{},
			Headers:          &[]string{},
			Methods:          &[]string{},
			AllowCredentials: pointer.To(false),
		}
	}

	block := input[0].(map[string]any)

	allowedOrigins := *pluginsdk.ExpandStringSlice(block["allowed_origins"].(*pluginsdk.Set).List())
	allowedHeaders := *pluginsdk.ExpandStringSlice(block["allowed_headers"].(*pluginsdk.Set).List())
	allowedMethods := *pluginsdk.ExpandStringSlice(block["allowed_methods"].(*pluginsdk.Set).List())

	cors := &fhirservices.FhirServiceCorsConfiguration{
		Origins:          &allowedOrigins,
		Headers:          &allowedHeaders,
		Methods:          &allowedMethods,
		AllowCredentials: pointer.To(block["credentials_allowed"].(bool)),
	}

	if v, ok := block["max_age_in_seconds"]; ok {
		cors.MaxAge = pointer.To(int64(v.(int)))
	}

	return cors
}

func expandFhirAcrLoginServer(input []any) *[]string {
	acrLoginServers := make([]string, 0)

	if len(input) == 0 {
		return &acrLoginServers
	}

	for _, item := range input {
		acrLoginServers = append(acrLoginServers, item.(string))
	}
	return &acrLoginServers
}
