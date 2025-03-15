package apimanagementworkspacelinks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementWorkspaceLinksResource struct {
	Etag       *string                      `json:"etag,omitempty"`
	Id         *string                      `json:"id,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties WorkspaceLinksBaseProperties `json:"properties"`
	Type       *string                      `json:"type,omitempty"`
}
