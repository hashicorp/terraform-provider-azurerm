package networkstatus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkStatusContract struct {
	ConnectivityStatus []ConnectivityStatusContract `json:"connectivityStatus"`
	DnsServers         []string                     `json:"dnsServers"`
}
