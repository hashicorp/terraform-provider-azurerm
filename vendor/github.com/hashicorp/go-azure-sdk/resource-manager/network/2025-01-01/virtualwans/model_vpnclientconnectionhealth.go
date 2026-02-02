package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnClientConnectionHealth struct {
	AllocatedIPAddresses         *[]string `json:"allocatedIpAddresses,omitempty"`
	TotalEgressBytesTransferred  *int64    `json:"totalEgressBytesTransferred,omitempty"`
	TotalIngressBytesTransferred *int64    `json:"totalIngressBytesTransferred,omitempty"`
	VpnClientConnectionsCount    *int64    `json:"vpnClientConnectionsCount,omitempty"`
}
