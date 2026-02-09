package openidconnectprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenidConnectProviderContractProperties struct {
	ClientId              string  `json:"clientId"`
	ClientSecret          *string `json:"clientSecret,omitempty"`
	Description           *string `json:"description,omitempty"`
	DisplayName           string  `json:"displayName"`
	MetadataEndpoint      string  `json:"metadataEndpoint"`
	UseInApiDocumentation *bool   `json:"useInApiDocumentation,omitempty"`
	UseInTestConsole      *bool   `json:"useInTestConsole,omitempty"`
}
