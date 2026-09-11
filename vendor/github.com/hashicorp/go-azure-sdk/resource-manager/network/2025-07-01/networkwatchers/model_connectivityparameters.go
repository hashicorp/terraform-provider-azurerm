package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityParameters struct {
	Destination           ConnectivityDestination `json:"destination"`
	PreferredIPVersion    *IPVersion              `json:"preferredIPVersion,omitempty"`
	Protocol              *Protocol               `json:"protocol,omitempty"`
	ProtocolConfiguration *ProtocolConfiguration  `json:"protocolConfiguration,omitempty"`
	Source                ConnectivitySource      `json:"source"`
}
