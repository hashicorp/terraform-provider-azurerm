package dnsresolvers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsResolverProperties struct {
	DnsResolverState  *DnsResolverState  `json:"dnsResolverState,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceGuid      *string            `json:"resourceGuid,omitempty"`
	VirtualNetwork    SubResource        `json:"virtualNetwork"`
}
