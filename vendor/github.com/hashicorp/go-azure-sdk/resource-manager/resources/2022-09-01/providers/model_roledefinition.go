package providers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleDefinition struct {
	Id            *string       `json:"id,omitempty"`
	IsServiceRole *bool         `json:"isServiceRole,omitempty"`
	Name          *string       `json:"name,omitempty"`
	Permissions   *[]Permission `json:"permissions,omitempty"`
	Scopes        *[]string     `json:"scopes,omitempty"`
}
