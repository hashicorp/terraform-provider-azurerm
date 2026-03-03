package managementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HierarchySettingsProperties struct {
	DefaultManagementGroup               *string `json:"defaultManagementGroup,omitempty"`
	RequireAuthorizationForGroupCreation *bool   `json:"requireAuthorizationForGroupCreation,omitempty"`
	TenantId                             *string `json:"tenantId,omitempty"`
}
