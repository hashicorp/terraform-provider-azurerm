package virtualnetworklinks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkLinkProperties struct {
	Metadata          *map[string]string `json:"metadata,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	VirtualNetwork    SubResource        `json:"virtualNetwork"`
}
