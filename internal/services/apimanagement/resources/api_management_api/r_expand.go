// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandApiManagementApiImport(importVs []any, apiType api.ApiType, soapApiType api.SoapApiType, path, serviceUrl, version, versionSetId string) *api.ApiCreateOrUpdateParameter {
	if len(importVs) == 0 || importVs[0] == nil {
		return nil
	}

	importV := importVs[0].(map[string]any)
	if len(importV) == 0 {
		return nil
	}

	contentFormat := importV["content_format"].(string)
	contentValue := importV["content_value"].(string)

	apiParams := api.ApiCreateOrUpdateParameter{
		Properties: &api.ApiCreateOrUpdateProperties{
			Type:    pointer.To(apiType),
			ApiType: pointer.To(soapApiType),
			Format:  pointer.ToEnum[api.ContentFormat](contentFormat),
			Value:   pointer.To(contentValue),
			Path:    path,
		},
	}

	wsdlSelectorVs := importV["wsdl_selector"].([]any)
	if len(wsdlSelectorVs) > 0 && wsdlSelectorVs[0] != nil {
		if wsdlSelectorV := wsdlSelectorVs[0].(map[string]any); len(wsdlSelectorV) > 0 {
			wSvcName := wsdlSelectorV["service_name"].(string)
			wEndpName := wsdlSelectorV["endpoint_name"].(string)

			apiParams.Properties.WsdlSelector = &api.ApiCreateOrUpdatePropertiesWsdlSelector{
				WsdlServiceName:  pointer.To(wSvcName),
				WsdlEndpointName: pointer.To(wEndpName),
			}
		}
	}
	if serviceUrl != "" {
		apiParams.Properties.ServiceURL = pointer.To(serviceUrl)
	}

	if version != "" {
		apiParams.Properties.ApiVersion = pointer.To(version)
	}

	if versionSetId != "" {
		apiParams.Properties.ApiVersionSetId = pointer.To(versionSetId)
	}

	return &apiParams
}

func expandApiManagementApiProtocols(input []any) *[]api.Protocol {
	if len(input) == 0 {
		return nil
	}
	results := make([]api.Protocol, 0)

	for _, v := range input {
		results = append(results, api.Protocol(v.(string)))
	}

	return &results
}

func expandApiManagementApiSubscriptionKeyParamNames(input []any) *api.SubscriptionKeyParameterNamesContract {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	query := v["query"].(string)
	header := v["header"].(string)
	contract := api.SubscriptionKeyParameterNamesContract{
		Query:  pointer.To(query),
		Header: pointer.To(header),
	}
	return &contract
}

func expandApiManagementOAuth2AuthenticationSettingsContract(input []any) *api.OAuth2AuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	oAuth2AuthorizationV := input[0].(map[string]any)
	return &api.OAuth2AuthenticationSettingsContract{
		AuthorizationServerId: pointer.To(oAuth2AuthorizationV["authorization_server_name"].(string)),
		Scope:                 pointer.To(oAuth2AuthorizationV["scope"].(string)),
	}
}

func expandApiManagementOpenIDAuthenticationSettingsContract(input []any) *api.OpenIdAuthenticationSettingsContract {
	if len(input) == 0 {
		return nil
	}

	openIDAuthorizationV := input[0].(map[string]any)
	return &api.OpenIdAuthenticationSettingsContract{
		OpenidProviderId:          pointer.To(openIDAuthorizationV["openid_provider_name"].(string)),
		BearerTokenSendingMethods: expandApiManagementOpenIDAuthenticationSettingsBearerTokenSendingMethods(openIDAuthorizationV["bearer_token_sending_methods"].(*pluginsdk.Set).List()),
	}
}

func expandApiManagementOpenIDAuthenticationSettingsBearerTokenSendingMethods(input []any) *[]api.BearerTokenSendingMethods {
	if input == nil {
		return nil
	}
	results := make([]api.BearerTokenSendingMethods, 0)

	for _, v := range input {
		results = append(results, api.BearerTokenSendingMethods(v.(string)))
	}

	return &results
}

func expandApiManagementApiContact(input []any) *api.ApiContactInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	return &api.ApiContactInformation{
		Email: pointer.To(v["email"].(string)),
		Name:  pointer.To(v["name"].(string)),
		Url:   pointer.To(v["url"].(string)),
	}
}

func expandApiManagementApiLicense(input []any) *api.ApiLicenseInformation {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)
	return &api.ApiLicenseInformation{
		Name: pointer.To(v["name"].(string)),
		Url:  pointer.To(v["url"].(string)),
	}
}
