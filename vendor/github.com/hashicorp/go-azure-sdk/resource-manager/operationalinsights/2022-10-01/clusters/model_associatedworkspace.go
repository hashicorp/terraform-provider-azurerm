package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssociatedWorkspace struct {
	AssociateDate *string `json:"associateDate,omitempty"`
	ResourceId    *string `json:"resourceId,omitempty"`
	WorkspaceId   *string `json:"workspaceId,omitempty"`
	WorkspaceName *string `json:"workspaceName,omitempty"`
}
