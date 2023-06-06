package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkAdapter struct {
	AdapterId            *string                   `json:"adapterId,omitempty"`
	AdapterPosition      *NetworkAdapterPosition   `json:"adapterPosition,omitempty"`
	DhcpStatus           *NetworkAdapterDHCPStatus `json:"dhcpStatus,omitempty"`
	DnsServers           *[]string                 `json:"dnsServers,omitempty"`
	IPv4Configuration    *IPv4Config               `json:"ipv4Configuration,omitempty"`
	IPv6Configuration    *IPv6Config               `json:"ipv6Configuration,omitempty"`
	IPv6LinkLocalAddress *string                   `json:"ipv6LinkLocalAddress,omitempty"`
	Index                *int64                    `json:"index,omitempty"`
	Label                *string                   `json:"label,omitempty"`
	LinkSpeed            *int64                    `json:"linkSpeed,omitempty"`
	MacAddress           *string                   `json:"macAddress,omitempty"`
	NetworkAdapterName   *string                   `json:"networkAdapterName,omitempty"`
	NodeId               *string                   `json:"nodeId,omitempty"`
	RdmaStatus           *NetworkAdapterRDMAStatus `json:"rdmaStatus,omitempty"`
	Status               *NetworkAdapterStatus     `json:"status,omitempty"`
}
