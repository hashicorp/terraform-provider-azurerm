package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KekIdentityProperties struct {
	UseSystemAssignedIdentity *interface{} `json:"useSystemAssignedIdentity,omitempty"`
	UserAssignedIdentity      *string      `json:"userAssignedIdentity,omitempty"`
}
