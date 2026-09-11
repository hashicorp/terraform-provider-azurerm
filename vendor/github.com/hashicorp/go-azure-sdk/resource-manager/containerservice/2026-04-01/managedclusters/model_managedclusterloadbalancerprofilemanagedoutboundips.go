package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterLoadBalancerProfileManagedOutboundIPs struct {
	Count     *int64 `json:"count,omitempty"`
	CountIPv6 *int64 `json:"countIPv6,omitempty"`
}
