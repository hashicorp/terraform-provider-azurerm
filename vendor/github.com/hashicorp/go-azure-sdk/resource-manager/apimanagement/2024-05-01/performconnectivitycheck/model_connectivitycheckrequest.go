package performconnectivitycheck

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityCheckRequest struct {
	Destination           ConnectivityCheckRequestDestination            `json:"destination"`
	PreferredIPVersion    *PreferredIPVersion                            `json:"preferredIPVersion,omitempty"`
	Protocol              *ConnectivityCheckProtocol                     `json:"protocol,omitempty"`
	ProtocolConfiguration *ConnectivityCheckRequestProtocolConfiguration `json:"protocolConfiguration,omitempty"`
	Source                ConnectivityCheckRequestSource                 `json:"source"`
}
