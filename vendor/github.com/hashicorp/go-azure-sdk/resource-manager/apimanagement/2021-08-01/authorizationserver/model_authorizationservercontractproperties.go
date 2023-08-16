package authorizationserver

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationServerContractProperties struct {
	AuthorizationEndpoint      string                        `json:"authorizationEndpoint"`
	AuthorizationMethods       *[]AuthorizationMethod        `json:"authorizationMethods,omitempty"`
	BearerTokenSendingMethods  *[]BearerTokenSendingMethod   `json:"bearerTokenSendingMethods,omitempty"`
	ClientAuthenticationMethod *[]ClientAuthenticationMethod `json:"clientAuthenticationMethod,omitempty"`
	ClientId                   string                        `json:"clientId"`
	ClientRegistrationEndpoint string                        `json:"clientRegistrationEndpoint"`
	ClientSecret               *string                       `json:"clientSecret,omitempty"`
	DefaultScope               *string                       `json:"defaultScope,omitempty"`
	Description                *string                       `json:"description,omitempty"`
	DisplayName                string                        `json:"displayName"`
	GrantTypes                 []GrantType                   `json:"grantTypes"`
	ResourceOwnerPassword      *string                       `json:"resourceOwnerPassword,omitempty"`
	ResourceOwnerUsername      *string                       `json:"resourceOwnerUsername,omitempty"`
	SupportState               *bool                         `json:"supportState,omitempty"`
	TokenBodyParameters        *[]TokenBodyParameterContract `json:"tokenBodyParameters,omitempty"`
	TokenEndpoint              *string                       `json:"tokenEndpoint,omitempty"`
}
