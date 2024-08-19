package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabricksProperties struct {
	DatabricksAccessToken *string `json:"databricksAccessToken,omitempty"`
	WorkspaceUrl          *string `json:"workspaceUrl,omitempty"`
}
