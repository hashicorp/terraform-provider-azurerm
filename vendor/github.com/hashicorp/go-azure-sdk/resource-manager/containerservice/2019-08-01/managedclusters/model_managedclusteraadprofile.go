package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAADProfile struct {
	ClientAppID     string  `json:"clientAppID"`
	ServerAppID     string  `json:"serverAppID"`
	ServerAppSecret *string `json:"serverAppSecret,omitempty"`
	TenantID        *string `json:"tenantID,omitempty"`
}
