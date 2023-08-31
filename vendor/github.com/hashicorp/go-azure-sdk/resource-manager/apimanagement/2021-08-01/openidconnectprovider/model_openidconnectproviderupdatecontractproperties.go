package openidconnectprovider

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenidConnectProviderUpdateContractProperties struct {
	ClientId         *string `json:"clientId,omitempty"`
	ClientSecret     *string `json:"clientSecret,omitempty"`
	Description      *string `json:"description,omitempty"`
	DisplayName      *string `json:"displayName,omitempty"`
	MetadataEndpoint *string `json:"metadataEndpoint,omitempty"`
}
