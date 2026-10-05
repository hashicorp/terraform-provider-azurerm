package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InfrastructureNetwork struct {
	DnsServers *[]string  `json:"dnsServers,omitempty"`
	Gateway    *string    `json:"gateway,omitempty"`
	IPPools    *[]IPPools `json:"ipPools,omitempty"`
	SubnetMask *string    `json:"subnetMask,omitempty"`
	UseDhcp    *bool      `json:"useDhcp,omitempty"`
}
