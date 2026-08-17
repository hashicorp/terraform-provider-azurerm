package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodelessConnectorPollingAuthProperties struct {
	ApiKeyIdentifier                     *string      `json:"apiKeyIdentifier,omitempty"`
	ApiKeyName                           *string      `json:"apiKeyName,omitempty"`
	AuthType                             string       `json:"authType"`
	AuthorizationEndpoint                *string      `json:"authorizationEndpoint,omitempty"`
	AuthorizationEndpointQueryParameters *interface{} `json:"authorizationEndpointQueryParameters,omitempty"`
	FlowName                             *string      `json:"flowName,omitempty"`
	IsApiKeyInPostPayload                *string      `json:"isApiKeyInPostPayload,omitempty"`
	IsClientSecretInHeader               *bool        `json:"isClientSecretInHeader,omitempty"`
	RedirectionEndpoint                  *string      `json:"redirectionEndpoint,omitempty"`
	Scope                                *string      `json:"scope,omitempty"`
	TokenEndpoint                        *string      `json:"tokenEndpoint,omitempty"`
	TokenEndpointHeaders                 *interface{} `json:"tokenEndpointHeaders,omitempty"`
	TokenEndpointQueryParameters         *interface{} `json:"tokenEndpointQueryParameters,omitempty"`
}
