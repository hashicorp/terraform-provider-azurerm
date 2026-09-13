package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceCustomBooleanParameter struct {
	Type  *CustomParameterType `json:"type,omitempty"`
	Value bool                 `json:"value"`
}
