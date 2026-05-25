package managementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupProperties struct {
	Children    *[]ManagementGroupChildInfo `json:"children,omitempty"`
	Details     *ManagementGroupDetails     `json:"details,omitempty"`
	DisplayName *string                     `json:"displayName,omitempty"`
	TenantId    *string                     `json:"tenantId,omitempty"`
}
