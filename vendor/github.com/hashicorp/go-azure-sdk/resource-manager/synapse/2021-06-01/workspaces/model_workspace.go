package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Workspace struct {
	Id         *string              `json:"id,omitempty"`
	Identity   *ManagedIdentity     `json:"identity,omitempty"`
	Location   string               `json:"location"`
	Name       *string              `json:"name,omitempty"`
	Properties *WorkspaceProperties `json:"properties,omitempty"`
	Tags       *map[string]string   `json:"tags,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
