package networksecurityperimeterlinkreferences

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NspLinkReferenceProperties struct {
	Description               *string                   `json:"description,omitempty"`
	LocalInboundProfiles      *[]string                 `json:"localInboundProfiles,omitempty"`
	LocalOutboundProfiles     *[]string                 `json:"localOutboundProfiles,omitempty"`
	ProvisioningState         *NspLinkProvisioningState `json:"provisioningState,omitempty"`
	RemoteInboundProfiles     *[]string                 `json:"remoteInboundProfiles,omitempty"`
	RemoteOutboundProfiles    *[]string                 `json:"remoteOutboundProfiles,omitempty"`
	RemotePerimeterGuid       *string                   `json:"remotePerimeterGuid,omitempty"`
	RemotePerimeterLocation   *string                   `json:"remotePerimeterLocation,omitempty"`
	RemotePerimeterResourceId *string                   `json:"remotePerimeterResourceId,omitempty"`
	Status                    *NspLinkStatus            `json:"status,omitempty"`
}
