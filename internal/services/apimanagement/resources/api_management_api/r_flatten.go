// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenApiManagementApiProtocols(input *[]api.Protocol) []string {
	if input == nil {
		return []string{}
	}

	results := make([]string, 0)
	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func flattenApiManagementApiSubscriptionKeyParamNames(paramNames *api.SubscriptionKeyParameterNamesContract) []any {
	if paramNames == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)

	result["header"] = pointer.From(paramNames.Header)
	result["query"] = pointer.From(paramNames.Query)

	return []any{result}
}

func flattenApiManagementOAuth2Authorization(input *api.OAuth2AuthenticationSettingsContract) []any {
	if input == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)

	result["authorization_server_name"] = pointer.From(input.AuthorizationServerId)
	result["scope"] = pointer.From(input.Scope)

	return []any{result}
}

func flattenApiManagementOpenIDAuthentication(input *api.OpenIdAuthenticationSettingsContract) []any {
	if input == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)

	result["openid_provider_name"] = pointer.From(input.OpenidProviderId)

	bearerTokenSendingMethods := make([]any, 0)
	if s := input.BearerTokenSendingMethods; s != nil {
		for _, v := range *s {
			bearerTokenSendingMethods = append(bearerTokenSendingMethods, string(v))
		}
	}
	result["bearer_token_sending_methods"] = pluginsdk.NewSet(pluginsdk.HashString, bearerTokenSendingMethods)

	return []any{result}
}

func flattenApiManagementApiContact(contact *api.ApiContactInformation) []any {
	if contact == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)

	result["email"] = pointer.From(contact.Email)
	result["name"] = pointer.From(contact.Name)
	result["url"] = pointer.From(contact.Url)

	return []any{result}
}

func flattenApiManagementApiLicense(license *api.ApiLicenseInformation) []any {
	if license == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)

	result["name"] = pointer.From(license.Name)
	result["url"] = pointer.From(license.Url)

	return []any{result}
}
