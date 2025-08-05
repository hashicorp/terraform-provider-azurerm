package authorizationserver

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationServerUpdateContractProperties struct {
	AuthorizationEndpoint      *string                       `json:"authorizationEndpoint,omitempty"`
	AuthorizationMethods       *[]AuthorizationMethod        `json:"authorizationMethods,omitempty"`
	BearerTokenSendingMethods  *[]BearerTokenSendingMethod   `json:"bearerTokenSendingMethods,omitempty"`
	ClientAuthenticationMethod *[]ClientAuthenticationMethod `json:"clientAuthenticationMethod,omitempty"`
	ClientId                   *string                       `json:"clientId,omitempty"`
	ClientRegistrationEndpoint *string                       `json:"clientRegistrationEndpoint,omitempty"`
	ClientSecret               *string                       `json:"clientSecret,omitempty"`
	DefaultScope               *string                       `json:"defaultScope,omitempty"`
	Description                *string                       `json:"description,omitempty"`
	DisplayName                *string                       `json:"displayName,omitempty"`
	GrantTypes                 *[]GrantType                  `json:"grantTypes,omitempty"`
	ResourceOwnerPassword      *string                       `json:"resourceOwnerPassword,omitempty"`
	ResourceOwnerUsername      *string                       `json:"resourceOwnerUsername,omitempty"`
	SupportState               *bool                         `json:"supportState,omitempty"`
	TokenBodyParameters        *[]TokenBodyParameterContract `json:"tokenBodyParameters,omitempty"`
	TokenEndpoint              *string                       `json:"tokenEndpoint,omitempty"`
	UseInApiDocumentation      *bool                         `json:"useInApiDocumentation,omitempty"`
	UseInTestConsole           *bool                         `json:"useInTestConsole,omitempty"`
}
