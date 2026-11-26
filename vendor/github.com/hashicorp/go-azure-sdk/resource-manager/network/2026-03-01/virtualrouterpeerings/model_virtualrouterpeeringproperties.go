package virtualrouterpeerings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualRouterPeeringProperties struct {
	PeerAsn           *int64             `json:"peerAsn,omitempty"`
	PeerIP            *string            `json:"peerIp,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
