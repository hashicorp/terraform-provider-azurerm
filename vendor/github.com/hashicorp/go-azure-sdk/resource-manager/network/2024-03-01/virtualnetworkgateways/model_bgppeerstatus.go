package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BgpPeerStatus struct {
	Asn               *int64        `json:"asn,omitempty"`
	ConnectedDuration *string       `json:"connectedDuration,omitempty"`
	LocalAddress      *string       `json:"localAddress,omitempty"`
	MessagesReceived  *int64        `json:"messagesReceived,omitempty"`
	MessagesSent      *int64        `json:"messagesSent,omitempty"`
	Neighbor          *string       `json:"neighbor,omitempty"`
	RoutesReceived    *int64        `json:"routesReceived,omitempty"`
	State             *BgpPeerState `json:"state,omitempty"`
}
