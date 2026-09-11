package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalDNSProfile struct {
	KubeDNSOverrides *map[string]LocalDNSOverride `json:"kubeDNSOverrides,omitempty"`
	Mode             *LocalDNSMode                `json:"mode,omitempty"`
	State            *LocalDNSState               `json:"state,omitempty"`
	VnetDNSOverrides *map[string]LocalDNSOverride `json:"vnetDNSOverrides,omitempty"`
}
