package virtualnetworklinks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkLinkProperties struct {
	ProvisioningState       *ProvisioningState       `json:"provisioningState,omitempty"`
	RegistrationEnabled     *bool                    `json:"registrationEnabled,omitempty"`
	ResolutionPolicy        *ResolutionPolicy        `json:"resolutionPolicy,omitempty"`
	VirtualNetwork          *SubResource             `json:"virtualNetwork,omitempty"`
	VirtualNetworkLinkState *VirtualNetworkLinkState `json:"virtualNetworkLinkState,omitempty"`
}
