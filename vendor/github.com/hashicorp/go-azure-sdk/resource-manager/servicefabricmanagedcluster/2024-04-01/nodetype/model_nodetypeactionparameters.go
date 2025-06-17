package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeActionParameters struct {
	Force      *bool       `json:"force,omitempty"`
	Nodes      *[]string   `json:"nodes,omitempty"`
	UpdateType *UpdateType `json:"updateType,omitempty"`
}
