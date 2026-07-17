package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnClientConnectionHealthDetail struct {
	EgressBytesTransferred    *int64  `json:"egressBytesTransferred,omitempty"`
	EgressPacketsTransferred  *int64  `json:"egressPacketsTransferred,omitempty"`
	IngressBytesTransferred   *int64  `json:"ingressBytesTransferred,omitempty"`
	IngressPacketsTransferred *int64  `json:"ingressPacketsTransferred,omitempty"`
	MaxBandwidth              *int64  `json:"maxBandwidth,omitempty"`
	MaxPacketsPerSecond       *int64  `json:"maxPacketsPerSecond,omitempty"`
	PrivateIPAddress          *string `json:"privateIpAddress,omitempty"`
	PublicIPAddress           *string `json:"publicIpAddress,omitempty"`
	VpnConnectionDuration     *int64  `json:"vpnConnectionDuration,omitempty"`
	VpnConnectionId           *string `json:"vpnConnectionId,omitempty"`
	VpnConnectionTime         *string `json:"vpnConnectionTime,omitempty"`
	VpnUserName               *string `json:"vpnUserName,omitempty"`
}
