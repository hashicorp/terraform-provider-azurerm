package servicelinker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VNetSolution struct {
	DeleteOrUpdateBehavior *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
	Type                   *VNetSolutionType       `json:"type,omitempty"`
}
