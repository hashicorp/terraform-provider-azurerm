package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalytics struct {
	LogType             *LogAnalyticsLogType `json:"logType,omitempty"`
	Metadata            *map[string]string   `json:"metadata,omitempty"`
	WorkspaceId         string               `json:"workspaceId"`
	WorkspaceKey        string               `json:"workspaceKey"`
	WorkspaceResourceId *string              `json:"workspaceResourceId,omitempty"`
}
