package connectedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AadProfile struct {
	AdminGroupObjectIDs *[]string `json:"adminGroupObjectIDs,omitempty"`
	EnableAzureRBAC     *bool     `json:"enableAzureRBAC,omitempty"`
	TenantID            *string   `json:"tenantID,omitempty"`
}
