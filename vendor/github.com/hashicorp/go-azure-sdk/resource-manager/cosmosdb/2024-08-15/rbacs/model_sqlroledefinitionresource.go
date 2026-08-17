package rbacs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlRoleDefinitionResource struct {
	AssignableScopes *[]string           `json:"assignableScopes,omitempty"`
	Permissions      *[]Permission       `json:"permissions,omitempty"`
	RoleName         *string             `json:"roleName,omitempty"`
	Type             *RoleDefinitionType `json:"type,omitempty"`
}
