package v2workspaceconnectionresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceConnectionManagedIdentity struct {
	ClientId   *string `json:"clientId,omitempty"`
	ResourceId *string `json:"resourceId,omitempty"`
}
