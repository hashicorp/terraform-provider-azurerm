package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAADProfile struct {
	AdminGroupObjectIDs *[]string `json:"adminGroupObjectIDs,omitempty"`
	ClientAppID         *string   `json:"clientAppID,omitempty"`
	EnableAzureRBAC     *bool     `json:"enableAzureRBAC,omitempty"`
	Managed             *bool     `json:"managed,omitempty"`
	ServerAppID         *string   `json:"serverAppID,omitempty"`
	ServerAppSecret     *string   `json:"serverAppSecret,omitempty"`
	TenantID            *string   `json:"tenantID,omitempty"`
}
