package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPodIdentityProfile struct {
	AllowNetworkPluginKubenet      *bool                                 `json:"allowNetworkPluginKubenet,omitempty"`
	Enabled                        *bool                                 `json:"enabled,omitempty"`
	UserAssignedIdentities         *[]ManagedClusterPodIdentity          `json:"userAssignedIdentities,omitempty"`
	UserAssignedIdentityExceptions *[]ManagedClusterPodIdentityException `json:"userAssignedIdentityExceptions,omitempty"`
}
