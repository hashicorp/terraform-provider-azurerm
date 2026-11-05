package playwrightworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PlaywrightWorkspaceUpdateProperties struct {
	LocalAuth        *EnablementStatus `json:"localAuth,omitempty"`
	RegionalAffinity *EnablementStatus `json:"regionalAffinity,omitempty"`
}
