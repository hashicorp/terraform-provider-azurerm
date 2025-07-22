package assetendpointprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetEndpointProfileProperties struct {
	AdditionalConfiguration           *string                     `json:"additionalConfiguration,omitempty"`
	Authentication                    *Authentication             `json:"authentication,omitempty"`
	DiscoveredAssetEndpointProfileRef *string                     `json:"discoveredAssetEndpointProfileRef,omitempty"`
	EndpointProfileType               string                      `json:"endpointProfileType"`
	ProvisioningState                 *ProvisioningState          `json:"provisioningState,omitempty"`
	Status                            *AssetEndpointProfileStatus `json:"status,omitempty"`
	TargetAddress                     string                      `json:"targetAddress"`
	Uuid                              *string                     `json:"uuid,omitempty"`
}
