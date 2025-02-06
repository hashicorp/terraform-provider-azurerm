package managementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupInfoProperties struct {
	DisplayName *string `json:"displayName,omitempty"`
	TenantId    *string `json:"tenantId,omitempty"`
}
