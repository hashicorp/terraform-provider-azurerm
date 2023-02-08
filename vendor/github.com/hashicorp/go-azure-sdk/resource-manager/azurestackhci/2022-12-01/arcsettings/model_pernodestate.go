package arcsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerNodeState struct {
	ArcInstance *string       `json:"arcInstance,omitempty"`
	Name        *string       `json:"name,omitempty"`
	State       *NodeArcState `json:"state,omitempty"`
}
