package dicomservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DicomServiceProperties struct {
	AuthenticationConfiguration *DicomServiceAuthenticationConfiguration `json:"authenticationConfiguration,omitempty"`
	CorsConfiguration           *CorsConfiguration                       `json:"corsConfiguration,omitempty"`
	EnableDataPartitions        *bool                                    `json:"enableDataPartitions,omitempty"`
	Encryption                  *Encryption                              `json:"encryption,omitempty"`
	EventState                  *ServiceEventState                       `json:"eventState,omitempty"`
	PrivateEndpointConnections  *[]PrivateEndpointConnection             `json:"privateEndpointConnections,omitempty"`
	ProvisioningState           *ProvisioningState                       `json:"provisioningState,omitempty"`
	PublicNetworkAccess         *PublicNetworkAccess                     `json:"publicNetworkAccess,omitempty"`
	ServiceURL                  *string                                  `json:"serviceUrl,omitempty"`
	StorageConfiguration        *StorageConfiguration                    `json:"storageConfiguration,omitempty"`
}
