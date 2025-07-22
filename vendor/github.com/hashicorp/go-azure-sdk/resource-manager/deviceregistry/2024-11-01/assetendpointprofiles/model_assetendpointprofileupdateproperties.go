package assetendpointprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetEndpointProfileUpdateProperties struct {
	AdditionalConfiguration *string               `json:"additionalConfiguration,omitempty"`
	Authentication          *AuthenticationUpdate `json:"authentication,omitempty"`
	EndpointProfileType     *string               `json:"endpointProfileType,omitempty"`
	TargetAddress           *string               `json:"targetAddress,omitempty"`
}
