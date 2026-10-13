// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_authorization_server

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
)

func flattenApiManagementAuthorizationServerGrantTypes(input []authorizationserver.GrantType) []any {
	outputs := make([]any, 0)

	for _, v := range input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func flattenApiManagementAuthorizationServerAuthorizationMethods(input *[]authorizationserver.AuthorizationMethod) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func flattenApiManagementAuthorizationServerBearerTokenSendingMethods(input *[]authorizationserver.BearerTokenSendingMethod) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func flattenApiManagementAuthorizationServerClientAuthenticationMethods(input *[]authorizationserver.ClientAuthenticationMethod) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		outputs = append(outputs, string(v))
	}

	return outputs
}

func flattenApiManagementAuthorizationServerTokenBodyParameters(input *[]authorizationserver.TokenBodyParameterContract) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]any)

		output["name"] = v.Name
		output["value"] = v.Value

		outputs = append(outputs, output)
	}

	return outputs
}
