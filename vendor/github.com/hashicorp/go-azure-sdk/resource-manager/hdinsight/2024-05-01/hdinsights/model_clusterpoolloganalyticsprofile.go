package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterPoolLogAnalyticsProfile struct {
	Enabled     bool    `json:"enabled"`
	WorkspaceId *string `json:"workspaceId,omitempty"`
}
