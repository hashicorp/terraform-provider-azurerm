package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAPIServerAccessProfile struct {
	AuthorizedIPRanges             *[]string `json:"authorizedIPRanges,omitempty"`
	DisableRunCommand              *bool     `json:"disableRunCommand,omitempty"`
	EnablePrivateCluster           *bool     `json:"enablePrivateCluster,omitempty"`
	EnablePrivateClusterPublicFQDN *bool     `json:"enablePrivateClusterPublicFQDN,omitempty"`
	PrivateDNSZone                 *string   `json:"privateDNSZone,omitempty"`
}
