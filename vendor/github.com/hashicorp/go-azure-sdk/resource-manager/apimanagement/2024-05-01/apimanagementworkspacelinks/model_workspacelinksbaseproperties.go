package apimanagementworkspacelinks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceLinksBaseProperties struct {
	Gateways    *[]WorkspaceLinksGateway `json:"gateways,omitempty"`
	WorkspaceId *string                  `json:"workspaceId,omitempty"`
}
