package packetcaptures

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCaptureFilter struct {
	LocalIPAddress  *string     `json:"localIPAddress,omitempty"`
	LocalPort       *string     `json:"localPort,omitempty"`
	Protocol        *PcProtocol `json:"protocol,omitempty"`
	RemoteIPAddress *string     `json:"remoteIPAddress,omitempty"`
	RemotePort      *string     `json:"remotePort,omitempty"`
}
