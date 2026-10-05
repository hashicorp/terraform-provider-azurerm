package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspacePatchInfo struct {
	Identity   *ManagedIdentity          `json:"identity,omitempty"`
	Properties *WorkspacePatchProperties `json:"properties,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
}
