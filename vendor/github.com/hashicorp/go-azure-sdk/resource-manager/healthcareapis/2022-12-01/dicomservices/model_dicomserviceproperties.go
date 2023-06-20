package dicomservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DicomServiceProperties struct {
	AuthenticationConfiguration *DicomServiceAuthenticationConfiguration `json:"authenticationConfiguration,omitempty"`
	CorsConfiguration           *CorsConfiguration                       `json:"corsConfiguration,omitempty"`
	EventState                  *ServiceEventState                       `json:"eventState,omitempty"`
	PrivateEndpointConnections  *[]PrivateEndpointConnection             `json:"privateEndpointConnections,omitempty"`
	ProvisioningState           *ProvisioningState                       `json:"provisioningState,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccess                     `json:"publicNetworkAccess,omitempty"`
	ServiceUrl                  *string                                  `json:"serviceUrl,omitempty"`
}
