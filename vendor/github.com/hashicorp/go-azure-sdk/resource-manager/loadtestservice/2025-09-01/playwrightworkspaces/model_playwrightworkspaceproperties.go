package playwrightworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlaywrightWorkspaceProperties struct {
	DataplaneUri      *string            `json:"dataplaneUri,omitempty"`
	LocalAuth         *EnablementStatus  `json:"localAuth,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	RegionalAffinity  *EnablementStatus  `json:"regionalAffinity,omitempty"`
	WorkspaceId       *string            `json:"workspaceId,omitempty"`
}
