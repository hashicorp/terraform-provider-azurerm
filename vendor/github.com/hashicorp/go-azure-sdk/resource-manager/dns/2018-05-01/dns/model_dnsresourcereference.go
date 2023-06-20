package dns

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsResourceReference struct {
	DnsResources   *[]SubResource `json:"dnsResources,omitempty"`
	TargetResource *SubResource   `json:"targetResource,omitempty"`
}
