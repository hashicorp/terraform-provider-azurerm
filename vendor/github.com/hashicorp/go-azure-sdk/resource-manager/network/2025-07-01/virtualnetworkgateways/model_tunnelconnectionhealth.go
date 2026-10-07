package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TunnelConnectionHealth struct {
	ConnectionStatus                 *VirtualNetworkGatewayConnectionStatus `json:"connectionStatus,omitempty"`
	EgressBytesTransferred           *int64                                 `json:"egressBytesTransferred,omitempty"`
	IngressBytesTransferred          *int64                                 `json:"ingressBytesTransferred,omitempty"`
	LastConnectionEstablishedUtcTime *string                                `json:"lastConnectionEstablishedUtcTime,omitempty"`
	Tunnel                           *string                                `json:"tunnel,omitempty"`
}
