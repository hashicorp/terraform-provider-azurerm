// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_authorization_server

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/authorizationserver"
)

func expandApiManagementAuthorizationServerGrantTypes(input []any) *[]authorizationserver.GrantType {
	outputs := make([]authorizationserver.GrantType, 0)

	for _, v := range input {
		grantType := authorizationserver.GrantType(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func expandApiManagementAuthorizationServerAuthorizationMethods(input []any) *[]authorizationserver.AuthorizationMethod {
	outputs := make([]authorizationserver.AuthorizationMethod, 0)

	for _, v := range input {
		grantType := authorizationserver.AuthorizationMethod(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func expandApiManagementAuthorizationServerBearerTokenSendingMethods(input []any) *[]authorizationserver.BearerTokenSendingMethod {
	outputs := make([]authorizationserver.BearerTokenSendingMethod, 0)

	for _, v := range input {
		grantType := authorizationserver.BearerTokenSendingMethod(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func expandApiManagementAuthorizationServerClientAuthenticationMethods(input []any) *[]authorizationserver.ClientAuthenticationMethod {
	outputs := make([]authorizationserver.ClientAuthenticationMethod, 0)

	for _, v := range input {
		grantType := authorizationserver.ClientAuthenticationMethod(v.(string))
		outputs = append(outputs, grantType)
	}

	return &outputs
}

func expandApiManagementAuthorizationServerTokenBodyParameters(input []any) *[]authorizationserver.TokenBodyParameterContract {
	outputs := make([]authorizationserver.TokenBodyParameterContract, 0)

	for _, v := range input {
		vs := v.(map[string]any)

		output := authorizationserver.TokenBodyParameterContract{
			Name:  vs["name"].(string),
			Value: vs["value"].(string),
		}
		outputs = append(outputs, output)
	}

	return &outputs
}
