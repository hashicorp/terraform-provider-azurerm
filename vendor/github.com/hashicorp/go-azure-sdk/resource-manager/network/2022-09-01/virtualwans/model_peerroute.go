package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeerRoute struct {
	AsPath       *string `json:"asPath,omitempty"`
	LocalAddress *string `json:"localAddress,omitempty"`
	Network      *string `json:"network,omitempty"`
	NextHop      *string `json:"nextHop,omitempty"`
	Origin       *string `json:"origin,omitempty"`
	SourcePeer   *string `json:"sourcePeer,omitempty"`
	Weight       *int64  `json:"weight,omitempty"`
}
