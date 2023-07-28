package v2workspaceconnectionresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceConnectionUsernamePassword struct {
	Password *string `json:"password,omitempty"`
	Username *string `json:"username,omitempty"`
}
