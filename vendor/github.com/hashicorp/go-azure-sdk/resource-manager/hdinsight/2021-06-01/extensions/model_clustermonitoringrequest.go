package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterMonitoringRequest struct {
	PrimaryKey  *string `json:"primaryKey,omitempty"`
	WorkspaceId *string `json:"workspaceId,omitempty"`
}
