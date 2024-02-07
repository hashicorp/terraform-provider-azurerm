package roledefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleDefinitionProperties struct {
	AssignableScopes *[]string     `json:"assignableScopes,omitempty"`
	Description      *string       `json:"description,omitempty"`
	Permissions      *[]Permission `json:"permissions,omitempty"`
	RoleName         *string       `json:"roleName,omitempty"`
	Type             *string       `json:"type,omitempty"`
}
